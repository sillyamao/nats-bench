package config

import "time"

// NatsConfig - nats benchmark client configurations
type NatsConfig struct {
	BaseConfig

	ConnTimeout time.Duration
}
