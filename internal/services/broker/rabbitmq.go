package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQBroker struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	host       string
}

func NewRabbitMQBroker(host string) *RabbitMQBroker {
	return &RabbitMQBroker{
		host: host,
	}
}

func (r *RabbitMQBroker) Connect() error {
	var err error
	r.connection, err = amqp.Dial(r.host)
	if err != nil {
		return err
	}

	r.channel, err = r.connection.Channel()
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQBroker) Publish(queueName string, message []byte) error {
	queue, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQBroker) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	if err := r.connection.Close(); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQBroker) Host() string {
	return r.host
}
