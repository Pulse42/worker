package consumerProducers

import "context"

type Type string

const (
	RedisConsumerProducerType Type = "redis"
)

type Generator func() ConsumerProducer

type ConsumerProducer interface {
	Subscribe(parameters string) error
	Send(message []byte) error
	Close() error
	Context() context.Context
	MessageChannel() <-chan string
	ErrorChannel() <-chan error
}

var consumerProducerMap = map[Type]Generator{
	RedisConsumerProducerType: NewRedisConsumerProducer,
}

func Factory(input Type) ConsumerProducer {
	if generator, ok := consumerProducerMap[input]; ok {
		return generator()
	}
	panic("unknown consumer type:" + input)
}
