package consumerProducers

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
)

type RedisConsumerProducerOptions struct {
	redis.Options
	Channel string
}

type RedisConsumerProducer struct {
	client         *redis.Client
	consumer       *redis.PubSub
	options        RedisConsumerProducerOptions
	ctx            context.Context
	messageChannel chan string
	errorChannel   chan error
	stopChannel    chan bool
}

func (r *RedisConsumerProducer) Subscribe(parameters string) error {
	if err := json.Unmarshal([]byte(parameters), &r.options); err != nil {
		return err
	}

	r.client = redis.NewClient(&r.options.Options)
	r.consumer = r.client.Subscribe(r.ctx, r.options.Channel)

	go func() {
		defer r.client.Close()
		defer r.consumer.Close()
		channel := r.consumer.Channel()
		for {
			select {
			case <-r.stopChannel:
				break
			case msg := <-channel:
				r.messageChannel <- msg.Payload
			}
		}
	}()

	return nil
}

func (r *RedisConsumerProducer) Send(message []byte) error {
	return r.client.Publish(r.ctx, r.options.Channel, message).Err()
}

func (r *RedisConsumerProducer) Close() error {
	r.stopChannel <- true
	return nil
}

func (r *RedisConsumerProducer) Context() context.Context {
	return r.ctx
}

func (r *RedisConsumerProducer) MessageChannel() <-chan string {
	return r.messageChannel
}

func (r *RedisConsumerProducer) ErrorChannel() <-chan error {
	return r.errorChannel
}

func NewRedisConsumerProducer() ConsumerProducer {
	return &RedisConsumerProducer{
		messageChannel: make(chan string),
		errorChannel:   make(chan error),
		stopChannel:    make(chan bool),
		ctx:            context.Background(),
	}
}
