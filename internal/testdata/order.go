package testdata

import (
	"time"

	"github.com/Numbone/practice0/internal/model"
)

func NewTestOrder() model.Order {
	return model.Order{
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
		Locale:          "en",
		CustomerID:      "cust-123",
		DeliveryService: "DHL",
		ShardKey:        "1",
		SMID:            1,
		DateCreated:     time.Now(),
		OofShard:        "1",
	}
}
