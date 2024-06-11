package main

import (
	"encoding/json"
	"log"
	"wbtech_l0/cache"
	db "wbtech_l0/database"
	"wbtech_l0/models"
	server "wbtech_l0/server"

	stan "github.com/nats-io/stan.go"
)

var storage *db.Storage

func main() {

	var err error // Подключение к БД, html серверу и nats-steaming, обработка ошибки в случае некорректного подключения к бд
	storage = &db.Storage{}
	storage.Db, err = db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	go server.StartServer(storage)
	SubscribeToNATS()

}

// SubscribeToNATS подключается к серверу NATS и подписывается на канал 'orders'
func SubscribeToNATS() {
	log.Println("Connecting to NATS...")
	sc, err := stan.Connect("test-cluster", "subscriber", stan.NatsURL("nats://localhost:4222")) //Подключение к кластеру
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer sc.Close()
	log.Println("Subscriber connected to NATS")

	_, err = sc.Subscribe("orders", func(m *stan.Msg) { //Подписка на orders, и Message handler(обработчик сообщений поступающих из этого канала)
		handleMessage(m) //Вызов обработчика во время получения сообщений в orders
	}, stan.DeliverAllAvailable(), stan.DurableName("my_durable"))
	if err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}
	log.Println("Subscribed to 'orders' topic")

	// Keep the application running to listen to messages
	select {}
}

// handleMessage обрабатывает полученные сообщения из топика 'orders'
func handleMessage(m *stan.Msg) {
	log.Println("Received a message")
	var order models.Order
	if err := json.Unmarshal(m.Data, &order); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}
	log.Printf("Message received: %+v", order)

	cache.SaveToCache(order)
	log.Println("Order saved to cache")

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshaling order to JSON: %v", err)
		return
	}

	if err := storage.SaveOrder(orderJSON); err != nil {
		log.Printf("Error saving order to DB: %v", err)
	} else {
		log.Printf("Order saved to database: %+v", order)
	}
}
