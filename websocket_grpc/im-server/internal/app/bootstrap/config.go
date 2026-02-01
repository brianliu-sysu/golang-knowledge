package bootstrap

import (
	"time"
)

type Config struct {
	NodeID         string        `yaml:"node_id"`
	HTTPAddr       string        `yaml:"http_addr"`        // for metrics/health
	WSAddr         string        `yaml:"ws_addr"`          // for websocket
	WSPath         string        `yaml:"ws_path"`          // for websocket
	ReadLimitBytes int64         `yaml:"read_limit_bytes"` // for websocket
	SendQueueSize  int           `yaml:"send_queue_size"`  // for websocket
	WriteTimeout   time.Duration `yaml:"write_timeout"`    // for websocket
	ReadTimeout    time.Duration `yaml:"read_timeout"`     // for websocket
	PingInterval   time.Duration `yaml:"ping_interval"`    // for websocket
	PongWait       time.Duration `yaml:"pong_wait"`        // for websocket

	PostgresDSN string `yaml:"postgres_dsn"`
}

func DefaultConfig() *Config {
	return &Config{
		NodeID:         "default_node_id",
		HTTPAddr:       ":8080",
		WSAddr:         ":8081",
		WSPath:         "/ws",
		ReadLimitBytes: 1024 * 1024,
		SendQueueSize:  100,
		WriteTimeout:   time.Second * 10,
		ReadTimeout:    time.Second * 10,
		PingInterval:   time.Second * 30,
		PongWait:       time.Second * 10,

		PostgresDSN: "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable",
	}
}
