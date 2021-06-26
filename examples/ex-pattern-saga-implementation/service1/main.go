package main

import (
	"context"
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
	fmt.Println("i am service1")

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
		"saga-queue-service1", "saga-consumer-service1",
	)
	if err != nil {
		logrus.Fatal(errpath.Err(err))
	}

	step, err := sagaClient.AddStep(
		"regular",
		"service1-step",
		nil,
		mysaga.Begining,
		func(num int) (int, error) {
			num = num + 100
			transcheck = num
			fmt.Println("!!!!!", transcheck)
			return num, nil
		},
		// func(num int) (int, error) {
		// 	num = num + 3
		// 	transcheck = num
		// 	return num, errors.New("am errrrror")
		// },
		func(num int) (int, error) {
			num = num - 100
			transcheck = num
			fmt.Println("!!!!!", transcheck)
			return num, nil
		},
		0,
		nil,
	)
	if err != nil {
		logrus.Fatal(errpath.Err(err))
	}

	if err = step.Play(13); err != nil {
		// logrus.Fatal(errpath.Err(err))
		logrus.Error(errpath.Err(err))
	}

	step.Listen()

	// ---

	http.HandleFunc("/", Hello)
	http.HandleFunc("/trans/", Transcheck)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Service1")
}

func Transcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service1  Transcheck: %v", transcheck)
}
