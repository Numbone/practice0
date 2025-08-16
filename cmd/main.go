package main

import (
	"github.com/Numbone/practice0/internal/cache"
	"github.com/Numbone/practice0/internal/db"
	"github.com/Numbone/practice0/internal/handler"
	"github.com/Numbone/practice0/internal/kafka"
	"github.com/Numbone/practice0/internal/testdata"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	_ = godotenv.Load()

	brokers := []string{os.Getenv("KAFKA_URL")}
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP_ID")

	connStr := os.Getenv("DB_URL")
	dbConn, err := db.Connect(connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer dbConn.Pool.Close()

	cacheLayer := cache.NewCache(1 * time.Minute)
	cacheLayer.DeleteUnCached(1 * time.Minute)

	orders, _ := dbConn.LoadAllOrders()
	for _, order := range orders {
		cacheLayer.Set(order)
	}

	go kafka.StartConsumer(brokers, topic, groupID, dbConn, cacheLayer)

	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()

	time.Sleep(2 * time.Second)

	testOrder := testdata.NewTestOrder()
	if err := testOrder.Validate(); err != nil {
		log.Fatal("invalid test order:", err)
	}
	if err := producer.SendOrder(testOrder); err != nil {
		log.Fatal("failed to send order:", err)
	}

	log.Println("Test order sent. Waiting for consumer to process...")
	time.Sleep(5 * time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", handler.OrdersHandler(cacheLayer))
	mux.HandleFunc("/order/", handler.OrderHandler(cacheLayer, dbConn))
	mux.HandleFunc("/", handler.RootHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP server started at %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")
	_ = srv.Shutdown(nil)
	log.Println("Server stopped")
}
