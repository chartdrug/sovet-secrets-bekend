package tbot

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	GetChatsQuestion(ctx context.Context) ([]entity.ChatsQuestion, error)
	CreateChatsQuestion(ctx context.Context, in entity.ChatsQuestion) error
	UpdateChatsQuestion(ctx context.Context, res entity.ChatsQuestion) error
	DeleteChatsQuestion(ctx context.Context, id string) error

	GetChatsAnswer(ctx context.Context) ([]entity.ChatsAnswer, error)
	CreateChatsAnswer(ctx context.Context, in entity.ChatsAnswer) error
	UpdateChatsAnswer(ctx context.Context, res entity.ChatsAnswer) error
	DeleteChatsAnswer(ctx context.Context, id string) error

	GetUsersQuestion(ctx context.Context) ([]entity.UsersQuestion, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// блок вопросов
func (r repository) GetChatsQuestion(ctx context.Context) ([]entity.ChatsQuestion, error) {
	var chatsQuestion []entity.ChatsQuestion
	err := r.db.With(ctx).Select().OrderBy("id desc").All(&chatsQuestion)
	return chatsQuestion, err
}

func (r repository) GetOneChatsQuestion(ctx context.Context, id string) (entity.ChatsQuestion, error) {
	var res entity.ChatsQuestion
	err := r.db.With(ctx).Select().Model(id, &res)
	return res, err
}

func (r repository) CreateChatsQuestion(ctx context.Context, in entity.ChatsQuestion) error {
	return r.db.With(ctx).Model(&in).Insert()
}

func (r repository) UpdateChatsQuestion(ctx context.Context, res entity.ChatsQuestion) error {
	return r.db.With(ctx).Model(&res).Update()
}

func (r repository) DeleteChatsQuestion(ctx context.Context, id string) error {
	res, err := r.GetOneChatsQuestion(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&res).Delete()
}

// блок для ответов
func (r repository) GetChatsAnswer(ctx context.Context) ([]entity.ChatsAnswer, error) {
	var chatsAnswer []entity.ChatsAnswer
	err := r.db.With(ctx).Select().OrderBy("id desc").All(&chatsAnswer)
	return chatsAnswer, err
}

func (r repository) GetOneChatsAnswer(ctx context.Context, id string) (entity.ChatsAnswer, error) {
	var res entity.ChatsAnswer
	err := r.db.With(ctx).Select().Model(id, &res)
	return res, err
}

func (r repository) CreateChatsAnswer(ctx context.Context, in entity.ChatsAnswer) error {
	return r.db.With(ctx).Model(&in).Insert()
}

func (r repository) UpdateChatsAnswer(ctx context.Context, res entity.ChatsAnswer) error {
	return r.db.With(ctx).Model(&res).Update()
}

func (r repository) DeleteChatsAnswer(ctx context.Context, id string) error {
	res, err := r.GetOneChatsAnswer(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&res).Delete()
}

// вопросы пользователей
func (r repository) GetUsersQuestion(ctx context.Context) ([]entity.UsersQuestion, error) {
	var rq []entity.UsersQuestion
	err := r.db.With(ctx).Select().OrderBy("id desc").All(&rq)
	return rq, err
}
