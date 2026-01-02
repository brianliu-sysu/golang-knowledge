package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	MaxConnectionKey = "max_connection"
	EnvKey           = "env"
)

type ConfigManager struct {
	config sync.Map
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (c *ConfigManager) Set(key string, valude any) {
	c.config.Store(key, valude)
}

func (c *ConfigManager) Load(key string) (any, bool) {
	return c.config.Load(key)
}

func main() {
	cm := NewConfigManager()
	cm.Set(MaxConnectionKey, 123)
	cm.Set(EnvKey, "prod")

	go func() {
		for {
			time.Sleep(time.Second)
			cm.Set(MaxConnectionKey, rand.Int31()%1000)
		}
	}()

	for {
		time.Sleep(time.Second)
		maxConn, _ := cm.Load(MaxConnectionKey)
		env, _ := cm.Load(EnvKey)
		fmt.Println("env:", env, ", maxConn:", maxConn)
	}
}
