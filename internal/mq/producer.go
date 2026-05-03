package mq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"

	"online-judge/internal/config"
)

type Producer struct {
	p     rocketmq.Producer
	topic string
}

func NewProducer(cfg *config.RocketMQConfig) (*Producer, error) {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.NameServers),
		producer.WithRetry(cfg.Retry),
	)
	if err != nil {
		return nil, fmt.Errorf("create rocketmq producer: %w", err)
	}
	if err := p.Start(); err != nil {
		return nil, fmt.Errorf("start rocketmq producer: %w", err)
	}
	return &Producer{p: p, topic: cfg.Topic}, nil
}

func (p *Producer) SendJudgeMessage(ctx context.Context, msg *JudgeMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	m := &primitive.Message{
		Topic: p.topic,
		Body:  body,
	}
	_, err = p.p.SendSync(ctx, m)
	return err
}

func (p *Producer) Close() error {
	return p.p.Shutdown()
}
