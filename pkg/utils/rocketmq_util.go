package utils

import (
	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
)

func NewProducer(topic string) (golang.Producer, error) {
	producer, err := golang.NewProducer(&golang.Config{
		Endpoint: Config.RocketMQ.Endpoint,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    Config.RocketMQ.AccessKey,
			AccessSecret: Config.RocketMQ.SecretKey,
		},
	},
		golang.WithTopics(topic),
	)
	if err != nil {
		return nil, err
	}
	// start producer
	err = producer.Start()
	if err != nil {
		return nil, err
	}

	return producer, nil
}
