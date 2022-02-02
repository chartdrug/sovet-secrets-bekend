package admin

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context) ([]entity.UsersAdmin, error)
}

// Album represents the data about an album.
type Users struct {
	entity.Users
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context) ([]entity.UsersAdmin, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return []entity.UsersAdmin{}, err
	}
	return users, nil
}
