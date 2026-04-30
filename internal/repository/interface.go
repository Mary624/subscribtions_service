package repository

import (
	"subscriptions_rest/internal/domain"
)

type Repository interface {
	GetSubscribe(clientId, seviceName string) (domain.Subscribe, error)
	GetSubscribesPrice(requestStruct domain.RepositoryRequest) (int, error)
	DeleteSubscribe(clientId, seviceName string) error
	UpdateSubscribe(subscribe domain.Subscribe) error
	AddSubscribe(subscribe domain.Subscribe) error
}
