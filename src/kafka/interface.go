package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kController interface {
	Write(ctx context.Context, key, message []byte, clientId string, topic string, writer kafka.Writer) error
	Read(ctx context.Context, topic string, clientId string, reader kafka.Reader) (kafka.Message, error)
}

// build kafka reader and writer instances to use for the producer/consumer methods
