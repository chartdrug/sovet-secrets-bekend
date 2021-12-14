package profile

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Users, error)
	GetByLogin(ctx context.Context, login string) (entity.Users, error)
	GetByEmail(ctx context.Context, email string) (entity.Users, error)
	Create(ctx context.Context, album entity.Users) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

func (r repository) Create(ctx context.Context, antro entity.Users) error {
	return r.db.With(ctx).Model(&antro).Insert()
}

func (r repository) GetByLogin(ctx context.Context, login string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"login": login}).One(&user)
	return user, err
}

func (r repository) GetByEmail(ctx context.Context, email string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"email": email}).One(&user)
	return user, err
}
