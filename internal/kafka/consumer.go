package kafka

import (
	"context"
	"encoding/json"
	"github.com/Numbone/practice0/internal/model"
	"log"

	"github.com/Numbone/practice0/internal/cache"
	"github.com/Numbone/practice0/internal/db"
	"github.com/segmentio/kafka-go"
)

func StartConsumer(brokers []string, topic, groupID string, dbConn *db.DB, cache *cache.Cache) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	log.Println("Kafka consumer started...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		if err := order.Validate(); err != nil {
			log.Println("Validation failed:", err)
			continue
		}

		if err := dbConn.SaveOrder(order); err != nil {
			log.Println("DB save error:", err)
			continue
		}

		cache.Set(order)
		log.Println("Order saved & cached:", order.OrderUID)
	}
}
