package app

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/uuid-disted/producer/internal/services/broker"
	"github.com/uuid-disted/producer/internal/services/generator"
)

type Application struct {
	mu      sync.Mutex
	brokers []*broker.RabbitMQBroker
	gen     generator.UUIDGenerator
}

var (
	uuidGenerationDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "uuid_generation_duration_seconds",
			Help:    "Histogram of durations for UUID generation",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"worker"},
	)
	messagePublishDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "message_publish_duration_seconds",
			Help:    "Histogram of durations for message publishing",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"broker"},
	)
)

func init() {
	prometheus.MustRegister(uuidGenerationDuration)
	prometheus.MustRegister(messagePublishDuration)
}

func NewApplication(brokersIP []string) *Application {
	epoch := time.Now()
	app := &Application{
		brokers: make([]*broker.RabbitMQBroker, len(brokersIP)),
		gen:     generator.New(1, epoch),
	}

	for i, ip := range brokersIP {
		newBroker := broker.NewRabbitMQBroker(ip)
		err := newBroker.Connect()
		if err != nil {
			fmt.Printf("Error connecting to broker %s: %v\n", ip, err)
			continue
		}
		app.brokers[i] = newBroker
	}

	return app
}

func (app *Application) ListBrokers() []string {
	app.mu.Lock()
	defer app.mu.Unlock()

	ips := make([]string, 0, len(app.brokers))
	for _, b := range app.brokers {
		if b != nil {
			ips = append(ips, b.Host())
		}
	}
	return ips
}

func (app *Application) PublishMessage(broker *broker.RabbitMQBroker, queueName string, message []byte) error {
	start := time.Now()
	err := broker.Publish(queueName, message)
	duration := time.Since(start).Seconds()
	messagePublishDuration.WithLabelValues(broker.Host()).Observe(duration)

	if err != nil {
		return err
	}
	return nil
}

func (app *Application) Run(queueName string, numUUIDs int, workerCount int) error {
	results := make(chan error, numUUIDs)
	var wg sync.WaitGroup

	worker := func(id int, uuids <-chan string) {
		defer wg.Done()

		for uuid := range uuids {
			start := time.Now()
			app.mu.Lock()
			for len(app.brokers) == 0 {
				app.mu.Unlock()
				time.Sleep(100 * time.Millisecond) // Wait and retry if no brokers are available
				app.mu.Lock()
			}

			for _, b := range app.brokers {
				if b != nil {
					err := app.PublishMessage(b, queueName, []byte(uuid))
					if err != nil {
						results <- fmt.Errorf("worker %d: error publishing to broker %s: %v", id, b.Host(), err)
					}
				}
			}
			app.mu.Unlock()
			duration := time.Since(start).Seconds()
			uuidGenerationDuration.WithLabelValues(fmt.Sprintf("worker-%d", id)).Observe(duration)
		}
	}

	uuidChan := make(chan string, numUUIDs)
	go func() {
		for i := 0; i < numUUIDs; i++ {
			uuidChan <- app.gen.Generate(time.Now())
		}
		close(uuidChan)
	}()

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(i, uuidChan)
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
