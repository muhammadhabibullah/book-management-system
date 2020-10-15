package middlewares

import (
	"sync"

	"book-management-system/configs"
)

var (
	m    *Middleware
	once sync.Once
)

// Middleware type
type Middleware struct {
	jwtKey string
}

// InitMiddleware returns new Middleware
func InitMiddleware() *Middleware {
	once.Do(func() {
		jwtCfg := configs.GetConfig().JWT
		m.jwtKey = jwtCfg.Key
	})

	return m
}
