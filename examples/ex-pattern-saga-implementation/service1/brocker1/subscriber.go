package brocker1

import (
	"context"

	"github.com/korovkin/limiter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	maxNewUsersAllowed = 10
)

func (r *Rabbit) HandleNewEvent(ctx context.Context, msg amqp.Delivery) error {

	// Decode incoming message
	var event string
	if err := r.Encoder.Decode(msg.Body, &event); err != nil {
		return errors.Wrap(err, "failed to decode a new user")
	}

	// _, err := s.Handler.Storage.Pg.User.CreateUser(ctx, &user)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to create new usert")
	// }

	logrus.Info("get event: ", event)

	return nil
}

// -------------------------------------------------

func (r *Rabbit) Listen() error {

	chann, err := r.Connection.Channel()
	if err != nil {
		return errors.Wrap(err, "get rabbit channel")
	}

	queue, err := chann.QueueDeclare(
		QueueNewEvent, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}

	err = chann.QueueBind(
		queue.Name, // queue name
		"",
		ExchangeEvent, // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}

	msgs, err := chann.Consume(
		queue.Name,          // queue
		ConsumerNewConsumer, // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	go r.handleNewUser(msgs) // handle incoming messages
	return nil
}

func (r *Rabbit) handleNewUser(messages <-chan amqp.Delivery) {

	limit := limiter.NewConcurrencyLimiter(maxNewUsersAllowed)
	defer limit.Wait()

	for {
		select {
		case msg := <-messages:

			if err := r.HandleNewEvent(context.Background(), msg); err != nil {
				logrus.Errorf("failed to handle a new event: %v", err)
			}

		}
	}
}
