package tickers

import (
	"sync"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/handler"
	"github.com/sirupsen/logrus"
)

var (
	wg     sync.WaitGroup
	closed = make(chan struct{}) // проверяется в самих тикерах, дабы не запускать новую итерацию
)

// Ticker -
type Ticker struct {
	handler *handler.Handler
}

// New -
func New(hndl *handler.Handler) *Ticker {
	return &Ticker{
		handler: hndl,
	}
}

// Start - run tickers
func (t *Ticker) Start() {
	t.startCheckUserRegTime()
}

// Wait -
func (*Ticker) Wait(shutdownTimeout time.Duration) {
	// try to shutdown the listener gracefully
	stoppedGracefully := make(chan struct{}, 1)
	go func() {
		// Notify subscribers about exit, wait for their work to be finished
		close(closed)
		wg.Wait()
		stoppedGracefully <- struct{}{}
	}()

	// wait for a graceful shutdown and then stop forcibly
	select {
	case <-stoppedGracefully:
		logrus.Infoln("tickers stopped gracefully")
	case <-time.After(shutdownTimeout):
		logrus.Infoln("tickers stopped forcibly")
	}
}
