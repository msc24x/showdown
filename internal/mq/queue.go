// Communication middleware between the application and the message broker.
package mq

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/msc24x/showdown/internal/config"
	"github.com/msc24x/showdown/internal/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

var connection *amqp.Connection
var channel *amqp.Channel
var queues = make(map[string]*amqp.Queue)

// Declare and create multiple rabbit mq queues if required.
func declareQueues(names ...string) {
	for _, name := range names {
		q_exe, err := channel.QueueDeclare(
			name,
			false,
			false,
			false,
			false,
			nil,
		)
		utils.PanicIf(err)
		queues[name] = &q_exe

		log.Printf("Declared RabbitMQ queue '%s'", name)
	}
}

// Prepares Showdown for queue operations and consumption.
func InitMessageQueue() func() {
	rmq_url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		config.RABBIT_MQ_USER,
		config.RABBIT_MQ_PASSWORD,
		config.RABBIT_MQ_HOST,
		config.RABBIT_MQ_PORT,
	)
	var err error

	connection, err = amqp.Dial(rmq_url)
	utils.PanicIf(err)
	channel, err = connection.Channel()
	utils.PanicIf(err)

	declareQueues("executables")

	return func() {
		connection.Close()
	}
}

// Queues given message into given queue.
func Queue(q_name string, retries int, body []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	queue := queues[q_name]

	if queue == nil {
		log.Fatalf("Trying to queue to undeclared queue '%s'", q_name)
	}

	err := channel.PublishWithContext(ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		if retries == 0 {
			utils.PanicIf(err)
		}

		Queue(q_name, retries-1, body)
	}
}

// Returns consumable stream of array like object of rabbit mq messages
// for given queue.
func Consume(q_name string) <-chan amqp.Delivery {
	deliveries, err := channel.Consume(
		q_name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.PanicIf(err)

	return deliveries
}
