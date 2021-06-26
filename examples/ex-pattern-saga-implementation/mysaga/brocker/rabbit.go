package brocker

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// Exchanges
const ExchangeEvent = "saga_event"

// Rabbit структура для работы с rabbit
type Rabbit struct {
	Connection *Connection
	Encoder    Encoder
	Credits    ConnCredits
	Queue      string
	Consumer   string
	msgs       chan []byte
}

// ConnCredits data for connection
type ConnCredits struct {
	User string
	URL  string
}

// NewConnection создает Rabbit инстанс
func NewConnection(c ConnCredits) (*Rabbit, error) {

	connString := fmt.Sprintf("amqp://%s@%s/", c.User, c.URL)
	conn, err := Dial(connString)
	if err != nil {
		err := fmt.Errorf("Error dialing Rabbit 0_o. %s", err)
		return nil, err
	}

	return &Rabbit{
		Connection: conn,
		Credits: ConnCredits{
			URL:  c.URL,
			User: c.User,
		},
		Encoder:  &JsonEncoder{},
		Queue:    fmt.Sprint("saga_queue_", uuid.Must(uuid.NewV4())),
		Consumer: fmt.Sprint("saga_consumer_", uuid.Must(uuid.NewV4())),
		msgs:     make(chan []byte),
	}, nil
}
