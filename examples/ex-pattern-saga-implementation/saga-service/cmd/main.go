package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/config"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/handler"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/logger"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/saga"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/server"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/saga-service/store"
	"github.com/PASPARTUUU/go_for_example/pkg/dye"
	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
)

const (
	defaultConfigPath     = "configs/home_pc_config.toml"
	defaultConfigPathLnx  = "configs/linux_notebook_config.toml"
	serverShutdownTimeout = 30 * time.Second
)

type ff func(ctx context.Context,arg interface{}) error

func main() {
	fmt.Println("i am alive")

	ctx := context.Background()

	// Parse flags
	configPath := flag.String("config", defaultConfigPathLnx, "configuration file path")
	flag.Parse()

	// Create logger
	logger := logger.New()

	// Parse configs
	cfg, err := config.Parse(*configPath)
	if err != nil {
		logger.Fatal(errpath.Err(err))
	}

	// ---

	store, err := store.New(ctx, cfg, logger)
	if err != nil {
		logger.Fatal(errpath.Err(err))
	}
	defer store.Pg.DB.Close()
	logger.Infoln(errpath.Infof("%+v", store.Pg.DB))

	hndl := handler.New(store, logger)

	router := server.NewRouter(logger)

	server.RestInit(router, hndl)

	// ---

	go server.Start(router, cfg.ServerPort)
	defer server.Stop(router, serverShutdownTimeout)

	// Wait for program exit
	<-NotifyAboutExit()
}

// NotifyAboutExit -
func NotifyAboutExit() <-chan os.Signal {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	return exit
}

func localsaga() {

	s := saga.NewSaga("saga name")

	x := 0 // saga will change x by adding 10 than adding 100
	err := s.AddStep(&saga.Step{
		Name:           "1",
		Func:           func(context.Context) error { x += 10; return nil },
		CompensateFunc: func(context.Context) error { x -= 10; return nil },
	})
	if err != nil {
		dye.Next(err)
		return
	}
	// err = s.AddStep(&saga.Step{
	// 	Name: "2",
	// 	// suppose function in second step returns error
	// 	Func:           func(context.Context) error { x += 100; return errors.New("step 2 err") },
	// 	CompensateFunc: func(context.Context) error { x -= 100; return nil },
	// })
	// if err != nil {
	// 	dye.Next(err)
	// }
	memory := saga.New()
	c := saga.NewCoordinator(context.Background(), context.Background(), s, memory, "id123")
	playRes := c.Play()
	err = playRes.ExecutionError
	if err != nil {
		dye.Next(err)
		// return
	}
	// x is still 0, because saga rolled back all applied steps
	dye.Next(x)

	lg, _ := memory.GetAllLogsByExecutionID("id123")

	for _, l := range lg {
		dye.Next(*l)
	}

}


/*

подключаемся к сага-оунеру

указываем исполняемую и комписационную функции

включаем ожидание сигнала




*/