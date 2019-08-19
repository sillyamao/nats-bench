package config

import "time"

type BaseConfig struct {
	URL string

	Duration       time.Duration
	ReportInterval time.Duration
	Rate           float64

	MsgNumber int // MsgNumber Must be greater than 0
	MsgSize   int

	SubjectPrefix     string
	SubjectNumber     int
	PubNumber         int
	SubNumber         int
	Async             bool
	PubAck            time.Duration
	PubAckMaxInflight int

	PingInterval int
	PingMax      int
}
