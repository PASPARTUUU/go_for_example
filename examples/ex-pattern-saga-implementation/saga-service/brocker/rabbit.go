package brocker

import (
	"fmt"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/config"
)

// Rabbit структура для работы с rabbit
type Rabbit struct {
	Connection      *Connection
	Encoder         Encoder
	Credits         ConnCredits
	headerPublisher string // имя отправителя для заголовока
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
		headerPublisher: "SagaOwner",
	}, nil
}
