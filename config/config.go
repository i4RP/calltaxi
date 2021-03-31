package config

import "os"

type Config struct {
	DatabaseURL            string
	PassengerChannelSecret string
	PassengerChannelToken  string
	OperatorChannelSecret  string
	OperatorChannelToken   string
}

func New() *Config {
	return &Config{
		DatabaseURL:            os.Getenv("DATABASE_URL"),
		PassengerChannelSecret: os.Getenv("PASSENGER_CHANNEL_SECRET"),
		PassengerChannelToken:  os.Getenv("PASSENGER_CHANNEL_TOKEN"),
		OperatorChannelSecret:  os.Getenv("OPERATOR_CHANNEL_SECRET"),
		OperatorChannelToken:   os.Getenv("OPERATOR_CHANNEL_TOKEN"),
	}
}
