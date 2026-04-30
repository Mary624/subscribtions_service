package repository

import (
	"context"
	"subscriptions_rest/internal/domain"
)

type Repository interface {
	GetSubscribtions(ctx context.Context) ([]domain.Subscribtion, error)
	GetSubscribtion(ctx context.Context, userId, seviceName string) (domain.Subscribtion, error)
	GetSubscribtionsPrice(ctx context.Context, requestStruct domain.Subscribtion) (int, error)
	DeleteSubscribtion(ctx context.Context, userId, seviceName string) error
	UpdateSubscribtion(ctx context.Context, subscribtion domain.Subscribtion) error
	AddSubscribtion(ctx context.Context, subscribtion domain.Subscribtion) error
}
