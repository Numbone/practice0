package cache

import (
	"github.com/Numbone/practice0/internal/model"
	"sync"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]model.Order),
	}
}

func (c *Cache) Get(id string) (model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orders[id]
	return order, ok
}

func (c *Cache) Set(order model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) LoadFromDB(orders []model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, o := range orders {
		c.orders[o.OrderUID] = o
	}
}

func (c *Cache) GetAll() []model.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	orders := make([]model.Order, 0, len(c.orders))
	for _, o := range c.orders {
		orders = append(orders, o)
	}
	return orders
}
