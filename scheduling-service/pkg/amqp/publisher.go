package amqp

import (
	"encoding/json"
	gorabbitmq "github.com/hadihammurabi/go-rabbitmq"
	"github.com/hadihammurabi/go-rabbitmq/exchange"
	"github.com/ronilsonalves/GoLang-in-a-spring-cloud-architecture/scheduling-service/internal/domain"
	amqpi "github.com/streadway/amqp"
	"log"
	"os"
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

// PublishMessage - send a msg to RabbitMQ queue when an appointment is made or updated
func PublishMessage(a domain.AppointmentDTO) {
	mq, err := ConnectRabbitMQ(os.Getenv("RABBIT_MQ_URL_CONN"), "appointment-service")
	log.Println(err)
	body, _ := json.Marshal(a)
	err = mq.Publish(&gorabbitmq.MQConfigPublish{
		RoutingKey: mq.Queue().Name,
		Message: amqpi.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	})
	defer mq.Close()
}
