package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PASPARTUUU/go_for_example/service/config"
	"github.com/PASPARTUUU/go_for_example/service/handler"
	"github.com/PASPARTUUU/go_for_example/service/logger"
	"github.com/PASPARTUUU/go_for_example/service/server"
	"github.com/PASPARTUUU/go_for_example/service/store"
	"github.com/PASPARTUUU/go_for_example/tools/errpath"
)

const (
	defaultConfigPath     = "configs/home_pc_config.toml"
	serverShutdownTimeout = 30 * time.Second
)

func main() {
	fmt.Println("i am alive")

	ctx := context.Background()

	// Parse flags
	configPath := flag.String("config", defaultConfigPath, "configuration file path")
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
