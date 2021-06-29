package rabpub

import (
	"sync"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/rabbitmq/rabbit"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.uber.org/multierr"
)

const headerPublisher = "mySuperService" // имя отправителя для заголовока
const channelUnreserved = "unreserved"

type Publisher struct {
	Rabbit  *rabbit.Rabbit
	Encoder rabbit.Encoder

	wg sync.WaitGroup
}

func New(rabConn *rabbit.Rabbit) (*Publisher, error) {
	pub := Publisher{
		Rabbit:  rabConn,
		Encoder: &rabbit.JsonEncoder{},
	}
	if err := pub.init(); err != nil {
		return nil, err
	}

	return &pub, nil
}

func (p *Publisher) init() error {
	// call all the initializers here, multierr package might be useful
	return multierr.Combine(
		p.initUsersTransferExchange(),
	)
}

func (p *Publisher) Wait(shutdownTimeout time.Duration) {
	// try to shutdown the listener gracefully
	stoppedGracefully := make(chan struct{}, 1)
	go func() {
		p.wg.Wait()
		stoppedGracefully <- struct{}{}
	}()

	// wait for a graceful shutdown and then stop forcibly
	select {
	case <-stoppedGracefully:
		logrus.Infoln("publisher stopped gracefully")
	case <-time.After(shutdownTimeout):
		logrus.Infoln("publisher stopped forcibly")
	}
}

func (p *Publisher) Publish(exchange, routingKey string, payload interface{}) error {
	senderChannel, err := p.Rabbit.GetSender(channelUnreserved)
	if err != nil {
		return errors.Wrapf(err, "failed to get a sender channel")
	}
	return p.publish(senderChannel, exchange, routingKey, payload)
}

func (p *Publisher) publish(channel *rabbit.Channel, exchange, routingKey string, payload interface{}) error {
	p.wg.Add(1)
	defer p.wg.Done()

	headers := make(amqp.Table)
	headers["publisher"] = headerPublisher

	body, err := p.Encoder.Encode(payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode the message")
	}

	err = channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  p.Encoder.ContentType(),
			Body:         body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return errors.Wrapf(err, "failed to send a message, exchange = %s, routing key = %s", exchange, routingKey)
	}

	logrus.WithFields(logrus.Fields{
		"event":   "publish to RabbitMQ",
		"payload": payload,
	}).Debugf("exchange = %s; key = %s", exchange, routingKey)

	return nil
}
