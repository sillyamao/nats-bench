package config

import "time"

// StanConfig - nats-streaming benchmark client configurations
type StanConfig struct {
	BaseConfig

	Cluster string

	ConnTimeout time.Duration
	StanTimeout time.Duration
}
