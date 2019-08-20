package config

import "time"

type BaseConfig struct {
	URL                string
	NatsConnectTimeout time.Duration
	PingInterval       int
	PingMax            int

	Duration       time.Duration
	ReportInterval time.Duration
	Rate           float64

	MsgNumber int // total msg number which will be published. if Duration is set, it will be ignored.
	MsgSize   int

	SubjectPrefix string
	SubjectNum    int

	PubNum     int // publisher number per subject
	SubNum     int // subscriber number per subject
	QueueGroup string

	PubAsync bool
}
