package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Numbone/practice0/internal/cache"
	"github.com/Numbone/practice0/internal/db"
	"github.com/Numbone/practice0/internal/kafka"
	"github.com/Numbone/practice0/internal/model"
)

func main() {
	_ = godotenv.Load()

	brokers := []string{"localhost:9092"}
	topic := "orders"
	groupID := "order-consumers"

	// Подключаемся к Postgres
	connStr := os.Getenv("DB_URL")
	dbConn, err := db.Connect(connStr)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer dbConn.Pool.Close()

	// Инициализация кэша
	cacheLayer := cache.NewCache()

	// Загружаем все заказы в кэш при старте
	orders, _ := dbConn.LoadAllOrders()
	for _, order := range orders {
		cacheLayer.Set(order)
	}

	// Стартуем консьюмер Kafka в отдельной горутине
	go kafka.StartConsumer(brokers, topic, groupID, dbConn, cacheLayer)

	// Создаём продюсер Kafka
	producer := kafka.NewProducer(brokers, topic)
	defer producer.Close()

	// Ждём подключения консьюмера
	time.Sleep(2 * time.Second)

	// Генерим тестовый заказ
	testOrder := model.Order{
		OrderUID:    "test-order-1",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:    "John Doe",
			Phone:   "+123456789",
			Zip:     "123456",
			City:    "Almaty",
			Address: "Some street 123",
			Region:  "KZ",
			Email:   "john@example.com",
		},
		Payment: model.Payment{
			Transaction:  "txn-123",
			Currency:     "KZT",
			Provider:     "Payme",
			Amount:       10000,
			PaymentDt:    time.Now().Unix(),
			Bank:         "Kaspi",
			DeliveryCost: 500,
			GoodsTotal:   9500,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      1,
				TrackNumber: "WBILMTESTTRACK",
				Price:       9500,
				RID:         "rid-123",
				Name:        "Sneakers",
				Sale:        0,
				Size:        "42",
				TotalPrice:  9500,
				NMID:        12345,
				Brand:       "Nike",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "cust-123",
		DeliveryService:   "DHL",
		ShardKey:          "1",
		SMID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Отправляем тестовый заказ в Kafka
	if err := producer.SendOrder(testOrder); err != nil {
		log.Fatal("failed to send order:", err)
	}

	log.Println("Test order sent. Waiting for consumer to process...")
	time.Sleep(5 * time.Second)

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		orders, err := dbConn.LoadAllOrders()
		if err != nil {
			log.Printf("failed to load orders: %v\n", err)
			http.Error(w, "failed to load orders", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			log.Printf("failed to encode orders response: %v\n", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	})

	// HTTP обработчик для получения заказа по ID
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		orderID := strings.TrimPrefix(r.URL.Path, "/order/")
		if orderID == "" {
			http.Error(w, "order id is required", http.StatusBadRequest)
			return
		}

		// Пытаемся достать из кэша
		if order, ok := cacheLayer.Get(orderID); ok {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(order)
			return
		}

		// Если нет в кэше — достаём из БД
		order, err := dbConn.GetOrder(orderID)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		// Кладём в кэш и возвращаем
		cacheLayer.Set(order)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(order)
	})

	// Отдаём index.html и статику
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		http.FileServer(http.Dir("web")).ServeHTTP(w, r)
	})

	log.Println("HTTP server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
