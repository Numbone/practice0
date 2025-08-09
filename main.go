package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SMID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func main() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		log.Fatal("failed to connect kafka ", err)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("failed receive controller ", err)
	}
	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))

	if err != nil {
		log.Fatal("failed to connect controller ", err)
	}
	defer func(controllerConn *kafka.Conn) {
		err := controllerConn.Close()
		if err != nil {
			log.Fatal("failed to close controller ", err)
		}
	}(controllerConn)

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             "orders",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatal("Ошибка создания топика:", err)
	}

	fmt.Println("Топик orders готов")

	ctx := context.Background()
	writer := kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "orders",
		Balancer: &kafka.LeastBytes{},
	}

	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			panic(err)
		}
	}(&writer)
	err = writer.WriteMessages(ctx,
		kafka.Message{

			Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		kafka.Message{

			Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		kafka.Message{
			Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "my-groupID",
	})
	defer reader.Close()

	msg, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Fatal("Ошибка при получении:", err)
	}

	fmt.Println(string(msg.Value))
}
