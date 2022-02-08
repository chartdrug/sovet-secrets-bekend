package admin

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context) ([]entity.UsersAdmin, error)
	CreateFeedBack(ctx context.Context, req CreateFeedBack, owner string) (entity.Feedback, error)
	getFeedback(ctx context.Context) ([]entity.Feedback, error)
}

// Album represents the data about an album.
type Users struct {
	entity.Users
}

type service struct {
	repo   Repository
	logger log.Logger
}

type CreateFeedBack struct {
	Email    string `json:"email"`
	Feedback string `json:"feedback"`
	Name     string `json:"name"`
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

func (s service) CreateFeedBack(ctx context.Context, req CreateFeedBack, owner string) (entity.Feedback, error) {

	id := entity.GenerateID()

	err := s.repo.CreateFeedBack(ctx, entity.Feedback{
		Id:       id,
		Owner:    owner,
		Dt:       time.Now(),
		Email:    req.Email,
		Name:     req.Name,
		Feedback: req.Feedback,
	})
	if err != nil {
		return entity.Feedback{}, err
	}
	return s.repo.GetFeedBack(ctx, id)
}

func (s service) getFeedback(ctx context.Context) ([]entity.Feedback, error) {
	feedback, err := s.repo.GetAllFeedback(ctx)
	if err != nil {
		return []entity.Feedback{}, err
	}
	return feedback, nil
}
