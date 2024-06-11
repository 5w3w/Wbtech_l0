package main

import (
	"encoding/json"
	"log"
	"os"
	s "wbtech_l0/models"

	stan "github.com/nats-io/stan.go"
)

func main() {
	publisher()

}

func publisher() {
	// Подключение к кластеру в NATS
	log.Println("Connecting to NATS...")
	sc, err := stan.Connect("test-cluster", "publisher", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer sc.Close()
	log.Println("Publisher connected to NATS")

	var order s.Order
	data, err := os.ReadFile("model.json") //Чтение битовой последовательности данных
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	err = json.Unmarshal(data, &order) //Создание и наполнение структуры данных
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	log.Printf("Order UID: %s", order.OrderUID)
	log.Printf("Track Number: %s", order.TrackNumber)

	jsonData, err := json.Marshal(order) //Сериализация обратно в JSON
	if err != nil {
		log.Fatalf("Error serializing JSON: %v", err)
	}

	log.Println("Publishing message to 'orders' topic...")

	err = sc.Publish("orders", jsonData) //Публикация данных
	if err != nil {
		log.Printf("Error publishing message: %v", err)
	} else {
		log.Println("Message published to NATS Streaming")
	}

}
