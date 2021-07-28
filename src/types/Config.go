package types

import "time"

type Config struct {
	// The interval which the client is required to respond by
	Heartbeat time.Duration

	// The duration to wait before a graceful client connection close
	CloseWait time.Duration
}
