package rabbit

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Rabbit структура для работы с rabbit
type Rabbit struct {
	Connection *Connection
	Credits    ConnCredits

	channels struct {
		m   map[string]*Channel
		mtx sync.Mutex
	}
}

// ConnCredits data for connection
type ConnCredits struct {
	User string
	URL  string
}

// NewConnection создает Rabbit инстанс
func NewConnection(cfg config.Rabbit) (*Rabbit, error) {

	connString := fmt.Sprintf("amqp://%s@%s/", cfg.RabbitUser, cfg.RabbitURL)
	conn, err := Dial(connString)
	if err != nil {
		err := fmt.Errorf("Error dialing Rabbit 0_o. %s", err)
		return nil, err
	}

	return &Rabbit{
		Connection: conn,
		Credits: ConnCredits{
			URL:  cfg.RabbitURL,
			User: cfg.RabbitUser,
		},
		channels: struct {
			m   map[string]*Channel
			mtx sync.Mutex
		}{
			m:   make(map[string]*Channel),
			mtx: sync.Mutex{},
		},
	}, nil
}

// CloseRabbit закрывает rabbitMQ соединение
func (rb *Rabbit) CloseRabbit() error {
	// Close all the channels first
	rb.channels.mtx.Lock()
	defer rb.channels.mtx.Unlock()

	var err error
	for _, ch := range rb.channels.m {
		err = multierr.Append(err, ch.Close())
	}

	// Then close the connection itself
	return multierr.Append(err, rb.Connection.Close())
}

func (rb *Rabbit) GetReceiver(name string) (*Channel, error) {
	return rb.getChannel("receiver." + name)
}

func (rb *Rabbit) GetSender(name string) (*Channel, error) {
	return rb.getChannel("sender." + name)
}

func (rb *Rabbit) getChannel(name string) (*Channel, error) {
	rb.channels.mtx.Lock()
	defer rb.channels.mtx.Unlock()

	if ch, found := rb.channels.m[name]; found {
		return ch, nil
	}

	ch, err := rb.Connection.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a channel")
	}
	rb.channels.m[name] = ch
	return ch, nil
}

func KeysCombine(keys ...string) string {
	if len(keys) == 0 {
		return ""
	}
	if len(keys) == 1 {
		return keys[0]
	}

	return strings.Join(keys, ".")
}
