package main

import (
	"math/rand"

	"github.com/streadway/amqp"
)

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func doRPC(jsonMessage string) (res string, err error) {

	conn, err := amqp.Dial(rabbitConnectionString)
	fatalOnError(err, "failed to connect to the rabbit")
	defer conn.Close()

	channel, err := conn.Channel()
	fatalOnError(err, "failed to open channel")
	defer channel.Close()

	// create an empty exclusive queue for the reply
	queue, err := channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	fatalOnError(err, "failed to declare reply queue")

	// consume the temporary queue that we created
	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	fatalOnError(err, "failed to register a consumer for reply_to queue")

	correlationId := randomString(32)

	err = channel.Publish(
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationId,
			ReplyTo:       queue.Name,
			Body:          []byte(jsonMessage),
		},
	)
	fatalOnError(err, "unable to publish message")

	for d := range msgs {
		if correlationId == d.CorrelationId {

			return string(d.Body), nil
		}
	}

	return "", nil
}
