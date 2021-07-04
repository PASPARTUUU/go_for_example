package rabsub

import (
	"context"

	"github.com/PASPARTUUU/go_for_example/pkg/lang"
	"github.com/PASPARTUUU/go_for_example/service/models"
	"github.com/PASPARTUUU/go_for_example/service/rabbitmq/rabbit"
	"github.com/korovkin/limiter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"go.uber.org/multierr"
)

const (
	maxNewUsersAllowed = 10
)

func (s *Subscriber) HandleNewUser(ctx context.Context, msg amqp.Delivery) error {
	// Decode incoming message
	var user models.User
	if err := s.Encoder.Decode(msg.Body, &user); err != nil {
		return errors.Wrap(err, "failed to decode a new user")
	}

	_, err := s.Handler.Storage.Pg.User.CreateUser(ctx, &user)
	if err != nil {
		return errors.Wrap(err, "failed to create new usert")
	}

	logrus.WithField("userID", user.ID).Debug("new user created successfully")

	return nil
}

func (s *Subscriber) HandleUpdatedUser(ctx context.Context, msg amqp.Delivery) error {
	// Decode incoming message
	var user models.User
	if err := s.Encoder.Decode(msg.Body, &user); err != nil {
		return errors.Wrap(err, "failed to decode user")
	}

	_, err := s.Handler.Storage.Pg.User.UpdateUser(ctx, &user)
	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	logrus.WithField("userID", user.ID).Debug("user updated successfully")

	return nil
}

func (s *Subscriber) HandleDeleteUser(ctx context.Context, msg amqp.Delivery) error {
	// Decode incoming message
	var user models.User
	if err := s.Encoder.Decode(msg.Body, &user); err != nil {
		return errors.Wrap(err, "failed to decode user")
	}

	err := s.Handler.Storage.Pg.User.DeleteUser(ctx, user.ID)
	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	logrus.WithField("userID", user.ID).Debug("user deleted successfully")

	return nil
}

func (s *Subscriber) initUserListener() error {
	return multierr.Combine(
		s.initNewUserListener(),
		s.initUpdateUserListener(),
		s.initDeleteUserListener(),
	)
}

// -------------------------------------------------

func (s *Subscriber) initNewUserListener() error {
	receiverChannel, err := s.Rabbit.GetReceiver(rabbit.QueueNewUser)
	if err != nil {
		return errors.Wrapf(err, "failed to get a receiver channel")
	}

	queue, err := receiverChannel.QueueDeclare(
		rabbit.QueueNewUser, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}

	err = receiverChannel.QueueBind(
		queue.Name, // queue name
		rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyNew), // routing key
		rabbit.ExchangeUser, // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}

	msgs, err := receiverChannel.Consume(
		queue.Name,             // queue
		rabbit.ConsumerNewUser, // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	s.wg.Add(1)
	go s.handleNewUser(msgs) // handle incoming messages
	return nil
}

func (s *Subscriber) handleNewUser(messages <-chan amqp.Delivery) {
	defer s.wg.Done()

	limit := limiter.NewConcurrencyLimiter(maxNewUsersAllowed)
	defer limit.Wait()

	for {
		select {
		case <-s.closed:
			return
		case msg := <-messages:
			// Start a new goroutine to handle multiple requests at the same time
			limit.Execute(lang.Recover(
				func() {
					if err := s.HandleNewUser(context.Background(), msg); err != nil {
						logrus.Errorf("failed to handle a new user: %v", err)
					}
				},
			))
		}
	}
}

// ---

func (s *Subscriber) initUpdateUserListener() error {
	receiverChannel, err := s.Rabbit.GetReceiver(rabbit.QueueNewUser)
	if err != nil {
		return errors.Wrapf(err, "failed to get a receiver channel")
	}

	queue, err := receiverChannel.QueueDeclare(
		rabbit.QueueUpdatedUser, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}

	err = receiverChannel.QueueBind(
		queue.Name, // queue name
		rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyUpdate), // routing key
		rabbit.ExchangeUser, // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}

	msgs, err := receiverChannel.Consume(
		queue.Name,                 // queue
		rabbit.ConsumerUpdatedUser, // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	s.wg.Add(1)
	go s.handleUpdateUser(msgs) // handle incoming messages
	return nil
}

func (s *Subscriber) handleUpdateUser(messages <-chan amqp.Delivery) {
	defer s.wg.Done()

	limit := limiter.NewConcurrencyLimiter(maxNewUsersAllowed)
	defer limit.Wait()

	for {
		select {
		case <-s.closed:
			return
		case msg := <-messages:
			// Start a new goroutine to handle multiple requests at the same time
			limit.Execute(lang.Recover(
				func() {
					if err := s.HandleUpdatedUser(context.Background(), msg); err != nil {
						logrus.Errorf("failed to handle update user: %v", err)
					}
				},
			))
		}
	}
}

// ---

func (s *Subscriber) initDeleteUserListener() error {
	receiverChannel, err := s.Rabbit.GetReceiver(rabbit.QueueNewUser)
	if err != nil {
		return errors.Wrapf(err, "failed to get a receiver channel")
	}

	queue, err := receiverChannel.QueueDeclare(
		rabbit.QueueDeletedUser, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return errors.Wrap(err, "failed to declare a queue")
	}

	err = receiverChannel.QueueBind(
		queue.Name, // queue name
		rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyDelete), // routing key
		rabbit.ExchangeUser, // exchange
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "failed to bind a queue")
	}

	msgs, err := receiverChannel.Consume(
		queue.Name,                 // queue
		rabbit.ConsumerDeletedUser, // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return errors.Wrap(err, "failed to consume from a channel")
	}

	s.wg.Add(1)
	go s.handleDeleteUser(msgs) // handle incoming messages
	return nil
}

func (s *Subscriber) handleDeleteUser(messages <-chan amqp.Delivery) {
	defer s.wg.Done()

	limit := limiter.NewConcurrencyLimiter(maxNewUsersAllowed)
	defer limit.Wait()

	for {
		select {
		case <-s.closed:
			return
		case msg := <-messages:
			// Start a new goroutine to handle multiple requests at the same time
			limit.Execute(lang.Recover(
				func() {
					if err := s.HandleDeleteUser(context.Background(), msg); err != nil {
						logrus.Errorf("failed to handle delete user: %v", err)
					}
				},
			))
		}
	}
}
