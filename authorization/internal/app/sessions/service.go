package sessions

import (
	"context"
)

type repository interface {
	add(ctx context.Context, s *SessionModel) (int64, error)
	get(ctx context.Context, session string) (*SessionModel, error)
	update(ctx context.Context, s *SessionModel)
}

type SessionRepository struct {
	repository repository
}

func NewSessionRepo(repository repository) *SessionRepository {
	return &SessionRepository{
		repository: repository,
	}
}

func (s *SessionRepository) Add(ctx context.Context, session *SessionModel) (int64, error) {
	return s.repository.add(ctx, session)
}

func (s *SessionRepository) Get(ctx context.Context, session string) (*SessionModel, error) {
	return s.repository.get(ctx, session)
}

func (s *SessionRepository) Update(ctx context.Context, session *SessionModel) {
	s.repository.update(ctx, session)
}
