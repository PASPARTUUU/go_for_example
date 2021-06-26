package brocker1

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const headerPublisher = "SagaService1" // имя отправителя для заголовока

func (r *Rabbit) Publish(exchange, routingKey string, payload interface{}) error {

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

	err = channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  r.Encoder.ContentType(),
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
