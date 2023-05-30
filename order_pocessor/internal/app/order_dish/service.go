package order_dish

import (
	"context"
)

type repository interface {
	add(ctx context.Context, order *OrderDishModel) (int64, error)
	getByID(ctx context.Context, id int64) (*OrderDishModel, error)
}

type OrderDishRepository struct {
	repository repository
}

func NewOrderDishRepo(repository repository) *OrderDishRepository {
	return &OrderDishRepository{
		repository: repository,
	}
}

func (s *OrderDishRepository) Add(ctx context.Context, order *OrderDishModel) (int64, error) {
	return s.repository.add(ctx, order)
}

func (s *OrderDishRepository) GetById(ctx context.Context, id int64) (*OrderDishModel, error) {
	return s.repository.getByID(ctx, id)
}
