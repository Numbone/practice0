package kafka

import (
	"context"
	"encoding/json"
	"github.com/Numbone/practice0/internal/model"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll, // ждём подтверждения от всех реплик
		},
	}
}

func (p *Producer) SendOrder(order model.Order) error {
	// сериализация в JSON
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := kafka.Message{
		Key:   []byte(order.OrderUID), // ключ для партиционирования
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return err
	}

	log.Println("Order sent to Kafka:", order.OrderUID)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
