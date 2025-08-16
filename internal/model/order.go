package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var validate = validator.New()

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SMID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
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
	ChrtID      int    `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NMID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func (o *Order) Validate() error {
	return validate.Struct(o)
}
