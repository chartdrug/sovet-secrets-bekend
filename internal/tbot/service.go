package tbot

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	GetChatsQuestion(ctx context.Context) ([]entity.ChatsQuestion, error)
	CreateChatsQuestion(ctx context.Context, req entity.ChatsQuestion) error
	UpdateChatsQuestion(ctx context.Context, req entity.ChatsQuestion) error
	DeleteChatsQuestion(ctx context.Context, id string) error

	GetChatsAnswer(ctx context.Context) ([]entity.ChatsAnswer, error)
	CreateChatsAnswer(ctx context.Context, req entity.ChatsAnswer) error
	UpdateChatsAnswer(ctx context.Context, req entity.ChatsAnswer) error
	DeleteChatsAnswer(ctx context.Context, id string) error

	GetUsersQuestion(ctx context.Context) ([]entity.UsersQuestion, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

//блок вопросов

func (s service) GetChatsQuestion(ctx context.Context) ([]entity.ChatsQuestion, error) {
	items, err := s.repo.GetChatsQuestion(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s service) CreateChatsQuestion(ctx context.Context, req entity.ChatsQuestion) error {
	err := s.repo.CreateChatsQuestion(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdateChatsQuestion(ctx context.Context, req entity.ChatsQuestion) error {
	err := s.repo.UpdateChatsQuestion(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteChatsQuestion(ctx context.Context, id string) error {
	return s.repo.DeleteChatsQuestion(ctx, id)
}

//блок ответов

func (s service) GetChatsAnswer(ctx context.Context) ([]entity.ChatsAnswer, error) {
	items, err := s.repo.GetChatsAnswer(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s service) CreateChatsAnswer(ctx context.Context, req entity.ChatsAnswer) error {
	err := s.repo.CreateChatsAnswer(ctx, entity.ChatsAnswer{
		Answer:   req.Answer,
		Category: req.Category,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdateChatsAnswer(ctx context.Context, req entity.ChatsAnswer) error {
	err := s.repo.UpdateChatsAnswer(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteChatsAnswer(ctx context.Context, id string) error {
	return s.repo.DeleteChatsAnswer(ctx, id)
}

// вопросы пользователей
func (s service) GetUsersQuestion(ctx context.Context) ([]entity.UsersQuestion, error) {
	items, err := s.repo.GetUsersQuestion(ctx)
	if err != nil {
		return nil, err
	}
	return items, nil
}
