package user

import "context"

type repository interface {
	add(ctx context.Context, u *UserModel) (int64, error)
	get(ctx context.Context, email string) (*UserModel, error)
	getByID(ctx context.Context, id int64) (*UserDTO, error)
}

type UsersRepository struct {
	repository repository
}

func NewUsersRepo(repository repository) *UsersRepository {
	return &UsersRepository{
		repository: repository,
	}
}

func (s *UsersRepository) Add(ctx context.Context, u *UserModel) (int64, error) {
	return s.repository.add(ctx, u)
}

func (s *UsersRepository) Get(ctx context.Context, email string) (*UserModel, error) {
	return s.repository.get(ctx, email)
}

func (s *UsersRepository) GetByID(ctx context.Context, id int64) (*UserDTO, error) {
	return s.repository.getByID(ctx, id)
}
