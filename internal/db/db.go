package db

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Numbone/practice0/internal/model"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect(connStr string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) SaveOrder(order model.Order) error {
	ctx := context.Background()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SMID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return fmt.Errorf("insert orders: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO deliveries (
            order_uid, name, phone, zip, city, address, region, email
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return fmt.Errorf("insert deliveries: %w", err)
	}

	_, err = tx.Exec(ctx, `
        INSERT INTO payments (
            order_uid, transaction, request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("insert payments: %w", err)
	}

	for _, item := range order.Items {
		_, err = tx.Exec(ctx, `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid, name, sale, size,
                total_price, nm_id, brand, status
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            ON CONFLICT DO NOTHING
        `,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status,
		)
		if err != nil {
			return fmt.Errorf("insert items: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (db *DB) GetOrder(id string) (model.Order, error) {
	var data []byte
	err := db.Pool.QueryRow(context.Background(),
		`SELECT data FROM orders_json WHERE order_uid=$1`, id).Scan(&data)
	if err != nil {
		return model.Order{}, nil
	}

	var order model.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return model.Order{}, nil
	}
	return order, err
}

func (db *DB) LoadAllOrders() ([]model.Order, error) {
	ctx := context.Background()

	rows, err := db.Pool.Query(ctx, `SELECT order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders`)
	if err != nil {
		return nil, fmt.Errorf("query orders: %w", err)
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SMID, &o.DateCreated, &o.OofShard)
		if err != nil {
			log.Println("scan order:", err)
			continue
		}

		err = db.Pool.QueryRow(ctx, `SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE order_uid=$1`, o.OrderUID).
			Scan(&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email)
		if err != nil {
			log.Println("load delivery:", err)
		}

		err = db.Pool.QueryRow(ctx, `SELECT transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE order_uid=$1`, o.OrderUID).
			Scan(&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount,
				&o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee)
		if err != nil {
			log.Println("load payment:", err)
		}

		itemRows, err := db.Pool.Query(ctx, `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid=$1`, o.OrderUID)
		if err != nil {
			log.Println("load items:", err)
		} else {
			var items []model.Item
			for itemRows.Next() {
				var i model.Item
				if err := itemRows.Scan(&i.ChrtID, &i.TrackNumber, &i.Price, &i.RID, &i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NMID, &i.Brand, &i.Status); err != nil {
					log.Println("scan item:", err)
					continue
				}
				items = append(items, i)
			}
			itemRows.Close()
			o.Items = items
		}

		orders = append(orders, o)
	}

	return orders, nil
}
