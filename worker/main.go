package main

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	rabbitConnectionString = "amqp://guest:guest@rmq:5672/"
)

type requestMsg struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type responseMsg struct {
	Sentence string `json:"sentence"`
}

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

			reqMsg := requestMsg{}
			_ = json.Unmarshal([]byte(msgBody), &reqMsg)

			resMsg := responseMsg{
				Sentence: fmt.Sprintf("%s %s is %d year's old", reqMsg.FirstName, reqMsg.LastName, reqMsg.Age),
			}

			resMsgJson, _ := json.Marshal(resMsg)

			err = channel.Publish(
				"",
				d.ReplyTo,
				false,
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: d.CorrelationId,
					Body:          resMsgJson,
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
