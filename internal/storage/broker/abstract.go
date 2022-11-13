package broker

import (
	"context"
)

// Handler provides message processing capabilities.
type Handler func(ctx context.Context, msg []byte) error

type Consumer interface {
	Run(ctx context.Context, handler Handler) error
	Shutdown() error
}

type Producer interface {
	WriteDeadMessages(ctx context.Context, srcTopic string, err error, msg []byte) error
	Shutdown() error
}
