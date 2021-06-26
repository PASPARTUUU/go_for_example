package brocker

import (
	"github.com/PASPARTUUU/go_for_example/pkg/dye"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	maxNewUsersAllowed = 10
)

func (r *Rabbit) HandleNewEvent() <-chan []byte {
	return r.msgs
}

// -------------------------------------------------

func (r *Rabbit) Listen() error {

	chann, err := r.Connection.Channel()
	if err != nil {
		return errors.Wrap(err, "get rabbit channel")
	}

	err = chann.ExchangeDeclare(
		ExchangeEvent, // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments

	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a exchange")
	}

	queue, err := chann.QueueDeclare(
		r.Queue, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
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
		queue.Name, // queue
		r.Consumer, // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	go func(messages <-chan amqp.Delivery) {
		for {
			dye.Next()
			msg := <-messages
			dye.Next(string(msg.Body))
			r.msgs <- msg.Body
		}
	}(msgs)

	return nil
}
