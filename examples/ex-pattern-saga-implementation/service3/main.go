package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/mysaga"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/mysaga/brocker"
	"github.com/PASPARTUUU/go_for_example/examples/ex-pattern-saga-implementation/service1/config"
	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/sirupsen/logrus"
)

// const defaultConfigPath = "examples/ex-pattern-saga-implementation/service1/config/cfg.toml"
const defaultConfigPath = "config/cfg.toml"

var transcheck int = -1

func main() {
	fmt.Println("i am service3")

	ctx := context.Background()
	_ = ctx

	// Parse flags
	configPath := flag.String("config", defaultConfigPath, "configuration file path")
	flag.Parse()

	// Parse configs
	cfg, err := config.Parse(*configPath)
	if err != nil {
		logrus.Fatal(errpath.Err(err))
	}

	// ---

	sagaClient, err := mysaga.NewClient(
		"",
		brocker.ConnCredits{
			User: cfg.Rabbit.RabbitUser,
			URL:  cfg.Rabbit.RabbitURL,
		},
		"saga-queue-service3", "saga-consumer-service3",
	)
	if err != nil {
		logrus.Fatal(errpath.Err(err))
	}

	step, err := sagaClient.AddStep(
		"regular",
		"service3-step",
		[]string{"service2-step"},
		mysaga.Begining,
		func(num int) (int, error) {
			num = num + 3
			transcheck = num
			return num, errors.New("am errrrror")
		},
		func(num int) (int, error) { num = num - 3; transcheck = num; return num, nil },
		0,
		nil,
	)
	if err != nil {
		logrus.Fatal(errpath.Err(err))
	}

	step.Listen()

	// ---

	http.HandleFunc("/", Hello)
	http.HandleFunc("/trans/", Transcheck)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Service3")
}

func Transcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service3  Transcheck: %v", transcheck)
}
