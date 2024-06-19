package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/uuid-disted/producer/app"
)

func main() {
	brokersFile := flag.String("f", "brokers.txt", "The file containing host address of RabbitMQ brokers IP")
	uuidNumber := flag.Int("n", 1000, "The number of UUIDs that must be generated overall")

	flag.Parse()

	brokersHost, err := readBrokersFile(*brokersFile)
	if err != nil {
		fmt.Printf("Error reading brokers file: %v\n", err)
		return
	}

	app := app.NewApplication(brokersHost)

	// Run the application
	err = app.Run("uuids", *uuidNumber)
	if err != nil {
		fmt.Printf("Error running generation process: %v\n", err)
	} else {
		fmt.Println("Generation process completed successfully")
	}
}

func readBrokersFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	brokers := strings.Split(strings.TrimSpace(string(data)), "\n")
	return brokers, nil
}
