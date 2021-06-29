package rabsub

import (
	"sync"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/handler"
	"github.com/PASPARTUUU/go_for_example/service/rabbitmq/rabbit"
	"github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

type Subscriber struct {
	Rabbit  *rabbit.Rabbit
	Encoder rabbit.Encoder
	Handler *handler.Handler

	wg     sync.WaitGroup
	closed chan struct{}
}

// Listen - create subscriber instance and start listening rabbit queues
func Listen(rabConn *rabbit.Rabbit, hndl *handler.Handler) (*Subscriber, error) {
	sub := Subscriber{
		Rabbit:  rabConn,
		Encoder: &rabbit.JsonEncoder{},
		Handler: hndl,
	}
	if err := sub.init(); err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *Subscriber) init() error {
	s.closed = make(chan struct{})

	// call all the initializers here, multierr package might be useful
	return multierr.Combine(
		s.initUserListener(),
	)
}

func (s *Subscriber) Wait(shutdownTimeout time.Duration) {
	// try to shutdown the listener gracefully
	stoppedGracefully := make(chan struct{}, 1)
	go func() {
		// Notify subscribers about exit, wait for their work to be finished
		close(s.closed)
		s.wg.Wait()
		stoppedGracefully <- struct{}{}
	}()

	// wait for a graceful shutdown and then stop forcibly
	select {
	case <-stoppedGracefully:
		logrus.Infoln("subscriber stopped gracefully")
	case <-time.After(shutdownTimeout):
		logrus.Infoln("subscriber stopped forcibly")
	}
}
