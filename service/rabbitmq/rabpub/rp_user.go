package rabpub

import (
	"github.com/PASPARTUUU/go_for_example/service/models"
	"github.com/PASPARTUUU/go_for_example/service/rabbitmq/rabbit"
	"github.com/pkg/errors"
)

const (
	channelNameUsersTransfer = "usersTransfer"
)

func (p *Publisher) initUsersTransferExchange() error {
	senderChannel, err := p.Rabbit.GetSender(channelNameUsersTransfer)
	if err != nil {
		return errors.Wrapf(err, "failed to get a sender channel")
	}

	err = senderChannel.ExchangeDeclare(
		rabbit.ExchangeUser, // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	return errors.Wrap(err, "failed to create an exchange")
}

func (p *Publisher) TransferNewUser(user models.User) error {
	senderChannel, err := p.Rabbit.GetSender(channelNameUsersTransfer)
	if err != nil {
		return errors.Wrapf(err, "failed to get a sender channel")
	}
	return p.publish(senderChannel, rabbit.ExchangeUser, rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyNew), user)
}

func (p *Publisher) TransferUpdatedUser(user models.User) error {
	senderChannel, err := p.Rabbit.GetSender(channelNameUsersTransfer)
	if err != nil {
		return errors.Wrapf(err, "failed to get a sender channel")
	}
	return p.publish(senderChannel, rabbit.ExchangeUser, rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyUpdate), user)
}

func (p *Publisher) TransferDeletedUser(user models.User) error {
	senderChannel, err := p.Rabbit.GetSender(channelNameUsersTransfer)
	if err != nil {
		return errors.Wrapf(err, "failed to get a sender channel")
	}
	return p.publish(senderChannel, rabbit.ExchangeUser, rabbit.KeysCombine(rabbit.KeyUser, rabbit.KeyDelete), user)
}
