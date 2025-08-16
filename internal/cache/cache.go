package cache

import (
	"github.com/Numbone/practice0/internal/model"
	"sync"
	"time"
)

type cacheItem struct {
	order     model.Order
	expiresAt time.Time
}

type Cache struct {
	mu     sync.RWMutex
	orders map[string]cacheItem
	ttl    time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		orders: make(map[string]cacheItem),
		ttl:    ttl,
	}
}

func (c *Cache) Get(id string) (model.Order, bool) {
	c.mu.RLock()
	item, ok := c.orders[id]
	c.mu.RUnlock()
	if !ok || time.Now().After(item.expiresAt) {
		return model.Order{}, false
	}
	return item.order, true
}

func (c *Cache) Set(order model.Order) {
	c.mu.Lock()
	c.orders[order.OrderUID] = cacheItem{
		order:     order,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

func (c *Cache) GetAll() []model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	orders := make([]model.Order, 0, len(c.orders))
	for _, item := range c.orders {
		if now.Before(item.expiresAt) {
			orders = append(orders, item.order)
		}
	}
	return orders
}

func (c *Cache) DeleteUnCached(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.mu.Lock()
			now := time.Now()
			for k, v := range c.orders {
				if now.After(v.expiresAt) {
					delete(c.orders, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
