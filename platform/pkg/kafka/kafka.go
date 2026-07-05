package kafka

import "context"

type Message struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int32
	Offset    int64
}

type Handler func(context.Context, Message) error

type Producer interface {
	Send(context.Context, []byte, []byte) error
	Close() error
}

type Consumer interface {
	Consume(context.Context, Handler) error
	Close() error
}
