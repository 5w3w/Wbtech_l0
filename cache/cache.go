package cache

import (
	"sync"
	"wbtech_l0/models"
)

var (
	orderCache = make(map[string]models.Order)
	mu         sync.Mutex
)

func SaveToCache(order models.Order) {
	mu.Lock()
	defer mu.Unlock()
	orderCache[order.OrderUID] = order
}

func GetOrderFromCache(id string) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()
	order, found := orderCache[id]
	return order, found
}
