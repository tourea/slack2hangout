package main

import (
    "context"
    "github.com/go-kit/kit/log"
)

type Service interface {
    Forward(ctx context.Context, space string, key string, token string, slackNotification SlackNotification) error
}

type service struct {
    logger log.Logger
}

func NewService(logger log.Logger) Service {
    return &service{logger: logger}
}

func (svc *service) Forward(ctx context.Context, space string, key string, token string, slackNotification SlackNotification) error {
    // TODO: convert to hangout message and send it
    return nil
}
