package dish

import "context"

type repository interface {
	add(ctx context.Context, dish *DishModel) (int64, error)
	getById(ctx context.Context, id int64) (*DishModel, error)
	list(ctx context.Context) ([]*DishModel, error)
	update(ctx context.Context, dish *DishModel) (bool, error)
	delete(ctx context.Context, id int64) (bool, error)
}

type DishRepository struct {
	repository repository
}

func NewDishRepo(repository repository) *DishRepository {
	return &DishRepository{
		repository: repository,
	}
}

func (s *DishRepository) Add(ctx context.Context, dish *DishModel) (int64, error) {
	return s.repository.add(ctx, dish)
}

func (s *DishRepository) GetById(ctx context.Context, id int64) (*DishModel, error) {
	return s.repository.getById(ctx, id)
}

func (s *DishRepository) List(ctx context.Context) ([]*DishModel, error) {
	return s.repository.list(ctx)
}

func (s *DishRepository) Update(ctx context.Context, dish *DishModel) (bool, error) {
	return s.repository.update(ctx, dish)
}

func (s *DishRepository) Delete(ctx context.Context, id int64) (bool, error) {
	return s.repository.delete(ctx, id)
}
