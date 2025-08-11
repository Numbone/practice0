package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func KafkaTest() {
	// Sample order data based on the provided model
	order := map[string]interface{}{
		"order_uid":    "b563feb7b2b84b6test",
		"track_number": "WBILMTESTTRACK",
		"entry":        "WBIL",
		"delivery": map[string]interface{}{
			"name":    "Test Testov",
			"phone":   "+9720000000",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "test@gmail.com",
		},
		"payment": map[string]interface{}{
			"transaction":   "b563feb7b2b84b6test",
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"goods_total":   317,
			"custom_fee":    0,
		},
		"items": []map[string]interface{}{
			{
				"chrt_id":      9934930,
				"track_number": "WBILMTESTTRACK",
				"price":        453,
				"rid":          "ab4219087a764ae0btest",
				"name":         "Mascaras",
				"sale":         30,
				"size":         "0",
				"total_price":  317,
				"nm_id":        2389212,
				"brand":        "Vivienne Sabo",
				"status":       202,
			},
		},
		"locale":             "en",
		"internal_signature": "",
		"customer_id":        "test",
		"delivery_service":   "meest",
		"shardkey":           "9",
		"sm_id":              99,
		"date_created":       time.Now().Format(time.RFC3339),
		"oof_shard":          "1",
	}

	// Create Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
	})
	defer writer.Close()

	// Convert to JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	// Send message
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(order["order_uid"].(string)),
			Value: orderJSON,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Message sent successfully: %s", order["order_uid"])

	// Send a few more test orders with different UIDs
	testOrders := []string{
		"order_001_test",
		"order_002_test",
		"order_003_test",
	}

	for i, uid := range testOrders {
		testOrder := order
		testOrder["order_uid"] = uid
		testOrder["customer_id"] = "customer_" + string(rune('1'+i))

		orderJSON, _ := json.Marshal(testOrder)
		err = writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(uid),
				Value: orderJSON,
			},
		)
		if err != nil {
			log.Printf("Error sending test order %s: %v", uid, err)
		} else {
			log.Printf("Test message sent: %s", uid)
		}

		time.Sleep(time.Second)
	}
}
