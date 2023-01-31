package amqp

import (
	gorabbitmq "github.com/hadihammurabi/go-rabbitmq"
	"github.com/hadihammurabi/go-rabbitmq/exchange"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// ConnectRabbitMQ connect and setup RabbitMQ channel and queue
func ConnectRabbitMQ(urlConn, name string) (gorabbitmq.MQ, error) {
	mq, err := gorabbitmq.New(urlConn)
	failOnError(err, "Failed to create a MQ")

	err = mq.Exchange().
		WithName(name).
		WithType(exchange.TypeDirect).
		Declare()
	failOnError(err, "Failed to create a channel")

	q, err := mq.Queue().
		WithName(name).
		Declare()
	failOnError(err, "Failed to create a queue")

	err = q.Binding().
		WithExchange(name).
		Bind()
	failOnError(err, "Failed to bind queue")

	return *mq, nil
}
