package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"

	"online-judge/internal/config"
)

type MessageHandler func(ctx context.Context, msg *JudgeMessage) error

type Consumer struct {
	c rocketmq.PushConsumer
}

func NewConsumer(cfg *config.RocketMQConfig, handler MessageHandler) (*Consumer, error) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(cfg.NameServers),
		consumer.WithGroupName(cfg.ConsumerGroup),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset),
	)
	if err != nil {
		return nil, fmt.Errorf("create rocketmq consumer: %w", err)
	}

	err = c.Subscribe(cfg.Topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, m := range msgs {
			var judgeMsg JudgeMessage
			if err := json.Unmarshal(m.Body, &judgeMsg); err != nil {
				return consumer.ConsumeSuccess, nil
			}
			if err := handler(ctx, &judgeMsg); err != nil {
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe topic: %w", err)
	}

	return &Consumer{c: c}, nil
}

func (c *Consumer) Start() error {
	return c.c.Start()
}

func (c *Consumer) Close() error {
	return c.c.Shutdown()
}
