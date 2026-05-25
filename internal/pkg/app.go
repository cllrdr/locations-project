package pkg

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"locations-project/internal/app/config"
	"locations-project/internal/app/handler"
)

type Application struct {
	Config    *config.Config
	Router    *gin.Engine
	Handler   *handler.Handler
	RedisPool *redis.Pool
}

// NewApp инициализирует приложение. Принимает готовый пул Redis.
func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler, pool *redis.Pool) *Application {
	return &Application{
		Config:    c,
		Router:    r,
		Handler:   h,
		RedisPool: pool,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")
	a.Handler.RegisterAPI(a.Router)

	addr := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(addr); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}

// Close gracefully closes Redis connections
func (a *Application) Close() {
	if a.RedisPool != nil {
		a.RedisPool.Close()
		logrus.Info("Redis pool closed")
	}
}

// NewRedisPool создаёт пул соединений для переиспользования
func NewRedisPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}