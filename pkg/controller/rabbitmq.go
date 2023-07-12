package controller

import (
	"fmt"
	"os"

	"vngitSub/pkg/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

//ReceiveMsg - Waiting for message from RabbitMQ
func ReceiveMsg() {
	logger := utils.ConfigZap()

	address := os.Getenv("ADDRESSRB")
	username := os.Getenv("USERRB")
	password := os.Getenv("PASSRB")
	port := os.Getenv("PORTRB")
	queue := os.Getenv("QUEUE")

	var _amqp string = fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, address, port)

	//Connect RabbitMQ Instance
	conn, err := amqp.Dial(_amqp)
	if err != nil {
		logger.Errorf("Connecting to RabbitMQ [%s:%s]...FAILED: %s", address, port, err)
	} else {
		logger.Debug("Connecting to RabbitMQ [%s:%s]...OK", address, port)
	}
	defer conn.Close()

	//Get RabbitMQ Channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Errorf("Openning a channel...FAILED: %s", err)
	} else {
		logger.Debug("Openning a channel...OK")
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queue, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		logger.Errorf("Registering a consumer...FAILED: %s", err)
	} else {
		logger.Debug("Registering a consumer...OK")
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logger.Infof("Incoming message: %s", d.Body)
			ChangeImage(d.Body)
		}
	}()

	logger.Info("[*] Waiting for messages. To exit press CTRL+C")
	<-forever
}