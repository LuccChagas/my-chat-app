package config

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"os"
)

func ConnRabbit() (*amqp091.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("AMQP_USER"),
		os.Getenv("AMQP_PASS"),
		os.Getenv("AMQP_HOST"),
		os.Getenv("AMQP_PORT"))

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
