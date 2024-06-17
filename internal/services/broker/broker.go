package broker

type Broker interface {
	Connect() error
	Publish(queueName string, message []byte) error
	Close() error
}
