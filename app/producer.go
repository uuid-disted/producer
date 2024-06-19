package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/uuid-disted/producer/internal/services/broker"
	"github.com/uuid-disted/producer/internal/services/generator"
)

type Application struct {
	brokers []*broker.RabbitMQBroker
	gen     generator.UUIDGenerator
}

func NewApplication(brokersHost []string) *Application {
	epoch := time.Now()
	app := &Application{
		brokers: make([]*broker.RabbitMQBroker, len(brokersHost)),
		gen:     generator.New(1, epoch),
	}

	for i, host := range brokersHost {
		newBroker := broker.NewRabbitMQBroker(host)
		err := newBroker.Connect()
		if err != nil {
			fmt.Printf("Error connecting to broker %s: %v\n", host, err)
			continue
		}
		app.brokers[i] = newBroker
	}

	return app
}

func (app *Application) PublishMessage(broker *broker.RabbitMQBroker, queueName string, message []byte) error {
	return broker.Publish(queueName, message)
}

func (app *Application) Run(queueName string, numUUIDs int) error {
	workerCount := len(app.brokers)
	results := make(chan error, numUUIDs)
	var wg sync.WaitGroup

	worker := func(id int, broker *broker.RabbitMQBroker, uuids chan string) {
		defer wg.Done()

		go func() {
			for i := 0; i < numUUIDs/workerCount; i++ {
				uuids <- app.gen.Generate(time.Now())
			}
			close(uuids)
		}()

		for uuid := range uuids {
			if broker != nil {
				err := app.PublishMessage(broker, queueName, []byte(uuid))
				if err != nil {
					results <- fmt.Errorf("worker %d: error publishing to broker %s: %v", id, broker.Host(), err)
				}
			}
		}
	}

	uuidsChannels := make([]chan string, workerCount)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		uuidsChannels[i] = make(chan string, numUUIDs/workerCount)
		go worker(i, app.brokers[i], uuidsChannels[i])
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for err := range results {
		if err != nil {
			fmt.Printf("Error occurred: %v\n", err)
			return err
		}
	}

	fmt.Println("Generation process completed successfully")
	return nil
}
