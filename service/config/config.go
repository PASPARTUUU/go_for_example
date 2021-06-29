package config

import (
	"github.com/PASPARTUUU/go_for_example/pkg/errpath"
	"github.com/spf13/viper"
)

const (
	// ServerPort -
	ServerPort string = "1324"
)

// Config -
type Config struct {
	ServerPort int

	Postgres Postgres
	Rabbit   Rabbit
}

// Postgres -
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Rabbit -
type Rabbit struct {
	RabbitURL  string
	RabbitUser string
}

// Parse will parse the configuration from the environment variables and a file with the specified path.
// Environment variables have more priority than ones specified in the file.
func Parse(filepath string) (*Config, error) {
	var cfg Config

	setDefaults()

	// Parse the file
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errpath.Err(err, "failed to read the config file")
	}

	// Unmarshal the config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errpath.Err(err, "failed to unmarshal the configuration")
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("ServerPort", 1324)
}
