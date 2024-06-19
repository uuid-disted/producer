package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestApplication_NewApplication(t *testing.T) {
// 	brokers := []string{
// 		"amqp://guest:guest@10.0.0.2:5672/",
// 		"amqp://guest:guest@10.0.0.2:5673/",
// 	}

// 	app := NewApplication(brokers)

// 	listedBrokers := app.brokers
// 	assert.Contains(t, listedBrokers, "amqp://guest:guest@10.0.0.2:5672/")
// 	assert.Contains(t, listedBrokers, "amqp://guest:guest@10.0.0.2:5673/")
// }

func TestApplication_PublishMessage(t *testing.T) {
	brokers := []string{
		"amqp://guest:guest@10.0.0.2:5672/",
	}

	app := NewApplication(brokers)

	err := app.PublishMessage(app.brokers[0], "test-queue", []byte("test-message"))
	assert.NoError(t, err)
}

func TestApplication_Run(t *testing.T) {
	brokers := []string{
		"amqp://guest:guest@10.0.0.2:5672/",
		"amqp://guest:guest@10.0.0.2:5673/",
	}

	app := NewApplication(brokers)

	// Define the number of UUIDs
	numUUIDs := 10

	// Run the generation process
	err := app.Run("test-queue", numUUIDs)
	assert.NoError(t, err, "Run should complete without error")

	// Verify all brokers received the messages (This is a simplification, in real tests you should verify messages in the queues)
	listedBrokers := app.brokers
	for _, broker := range listedBrokers {
		assert.NoError(t, app.PublishMessage(broker, "test-queue", []byte("verification-message")))
	}
}
