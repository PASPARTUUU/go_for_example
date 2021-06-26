package brocker

import (
	"github.com/PASPARTUUU/go_for_example/pkg/dye"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const headerPublisher = "SagaService" // имя отправителя для заголовока

func (r *Rabbit) Publish(payload interface{}) error {

	headers := make(amqp.Table)
	headers["publisher"] = headerPublisher
	channel, err := r.Connection.Channel()
	if err != nil {
		return errors.Wrap(err, "failed to encode the message")
	}

	body, err := r.Encoder.Encode(payload)
	if err != nil {
		return errors.Wrap(err, "failed to encode the message")
	}

	dye.Next("publish", payload)

	err = channel.Publish(
		ExchangeEvent, // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:  r.Encoder.ContentType(),
			Body:         body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return errors.Wrapf(err, "failed to send a message, exchange = %s", ExchangeEvent)
	}

	logrus.WithFields(logrus.Fields{
		"event":   "publish to RabbitMQ",
		"payload": payload,
	}).Debugf("exchange = %s; ", ExchangeEvent)

	return nil
}
