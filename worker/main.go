package main

import (
	"github.com/streadway/amqp"
	"github.com/Sirupsen/logrus"
)

const (
	rabbitConnectionString = "amqp://guest:guest@rmq:5672/"
)

func main() {

	conn, err := amqp.Dial(rabbitConnectionString)
	fatalOnError(err, "failed to connect to the rabbit")
	defer conn.Close()

	channel, err := conn.Channel()
	fatalOnError(err, "failed to open channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	fatalOnError(err, "failed to declare queue")

	err = channel.Qos(1, 0, false)
	fatalOnError(err, "failed to set Quality of Service")

	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	fatalOnError(err, "failed to register a consumer")

	go func() {
		for d := range msgs {

			msgBody := string(d.Body)

			logrus.Infof("msg: %s", string(msgBody))

			err = channel.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          []byte(msgBody + " reply"),
				},
			)
			fatalOnError(err, "failed to publish reply")

			d.Ack(false)
		}
	}()

	forever := make(chan bool)

	logrus.Info("awaiting RPC requests")
	<-forever
}

func fatalOnError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}
