package pkg

import "github.com/rabbitmq/amqp091-go"

type RabbitMQConnection interface {
	Channel() (*amqp091.Channel, error)
}
