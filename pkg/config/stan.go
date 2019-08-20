package config

import "time"

// StanConfig - nats-streaming benchmark client configurations
type StanConfig struct {
	BaseConfig

	Cluster            string
	StanConnectTimeout time.Duration
	PubAckWait         time.Duration
	PubAckMaxInflight  int
}
