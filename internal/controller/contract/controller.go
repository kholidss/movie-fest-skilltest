package contract

import (
	"context"
	"github.com/kholidss/movie-fest-skilltest/internal/appctx"
	"github.com/kholidss/movie-fest-skilltest/pkg/pubsubx"
	"github.com/rabbitmq/amqp091-go"
)

type PubSubMessageController interface {
	Serve(ctx context.Context, message *pubsubx.Message)
}

type MessageController interface {
	Serve(data amqp091.Delivery) error
}

type Controller interface {
	Serve(xCtx appctx.Data) appctx.Response
}
