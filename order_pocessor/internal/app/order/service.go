package order

import (
	"context"
)

type repository interface {
	add(ctx context.Context, order *OrderModel) (int64, error)
	getByID(ctx context.Context, id int64) (*OrderModel, error)
	changeStatus(ctx context.Context, id int64, status string) (bool, error)
	getPendingOrder(ctx context.Context) (*OrderModel, error)
}

type OrderRepository struct {
	repository repository
}

func NewOrderRepo(repository repository) *OrderRepository {
	return &OrderRepository{
		repository: repository,
	}
}

func (s *OrderRepository) Add(ctx context.Context, order *OrderModel) (int64, error) {
	return s.repository.add(ctx, order)
}

func (s *OrderRepository) GetById(ctx context.Context, id int64) (*OrderModel, error) {
	return s.repository.getByID(ctx, id)
}

func (s *OrderRepository) ChangeStatus(ctx context.Context, id int64, status string) (bool, error) {
	return s.repository.changeStatus(ctx, id, status)
}

func (s *OrderRepository) GetPendingOrder(ctx context.Context) (*OrderModel, error) {
	return s.repository.getPendingOrder(ctx)
}
