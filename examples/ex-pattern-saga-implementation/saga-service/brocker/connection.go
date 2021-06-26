package brocker

import (
	"time"

	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const reconnectTimeout = 1 * time.Second

// Connection is an amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
}

// Dial wraps amqp.Dial, dials and returns an auto-reconnected connection
func Dial(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			errChan := make(chan *amqp.Error)
			reason, ok := <-connection.Connection.NotifyClose(errChan)
			// Exit the goroutine if the connection intentionally closed by developer
			if !ok {
				break
			}
			logrus.Infoln("connection is accidentally closed: %v", reason)

			// Reconnect if not closed by developer
			for {
				time.Sleep(reconnectTimeout)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					logrus.Infoln("successfully reconnected")
					break
				}

				logrus.Warnf("reconnection failed: %v", err)
			}
		}
	}()

	return connection, nil
}

// Channel wraps amqp.Connection.Channel to return an auto-reconnected channel
func (c *Connection) Channel() (*Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			errChan := make(chan *amqp.Error)
			reason, ok := <-channel.Channel.NotifyClose(errChan)
			// Exit the goroutine if the connection intentionally closed by developer
			if !ok || channel.IsClosed() {
				_ = channel.Close() // close the channel, ensure closed flag is set
				break
			}
			logrus.Infof("channel is accidentally closed: %v", reason)

			// Reconnect if not closed by developer
			for {
				time.Sleep(reconnectTimeout)

				ch, err := c.Connection.Channel()
				if err == nil {
					logrus.Infoln("channel is successfully recreated")
					channel.Channel = ch
					break
				}

				logrus.Infof("channel recreation failed: %v", err)
			}
		}
	}()

	return channel, nil
}

// Channel wraps amqp.Channel
type Channel struct {
	*amqp.Channel
	closed int32
}

// IsClosed indicates the channel is closed by developer
func (c *Channel) IsClosed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

// Close ensures closed flag is set
func (c *Channel) Close() error {
	if c.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&c.closed, 1)
	return c.Channel.Close()
}

// Consume wraps amqp.Channel.Consume, the returned delivery will only end only when channel is closed by developer
func (c *Channel) Consume(
	queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table,
) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := c.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {

				time.Sleep(reconnectTimeout)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			time.Sleep(reconnectTimeout) // make sure closed flag is set

			if c.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
