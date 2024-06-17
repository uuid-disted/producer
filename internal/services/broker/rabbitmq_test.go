package broker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRabbitMQBroker(t *testing.T) {
	url := "amqp://guest:guest@10.0.0.2:5672/"
	broker := NewRabbitMQBroker(url)

	// Test Connect
	err := broker.Connect()
	assert.NoError(t, err, "should connect to RabbitMQ without error")

	defer func() {
		err := broker.Close()
		assert.NoError(t, err, "should close RabbitMQ connection without error")
	}()

	// Test Host method
	assert.Equal(t, url, broker.Host(), "should return the correct host")

	// Test Publish
	queueName := "test-queue"
	message := []byte("test-message")
	err = broker.Publish(queueName, message)
	assert.NoError(t, err, "should publish message to RabbitMQ without error")

	// Test if message is published
	ch, err := broker.connection.Channel()
	assert.NoError(t, err, "should create a channel without error")
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	assert.NoError(t, err, "should declare queue without error")

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	assert.NoError(t, err, "should consume messages without error")

	go func() {
		for d := range msgs {
			assert.Equal(t, message, d.Body, "should receive the published message")
			break
		}
	}()

	// Wait a bit for the message to be consumed
	time.Sleep(1 * time.Second)
}
