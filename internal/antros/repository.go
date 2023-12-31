package antros

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Antro, error)
	GetOne(ctx context.Context, id string) (entity.Antro, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, album entity.Antro) error
	Update(ctx context.Context, album entity.Antro) error
	GetProfile(ctx context.Context, id string) (entity.Users, error)
	GetByLogin(ctx context.Context, login string) (entity.Users, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, owner string) ([]entity.Antro, error) {
	var antro []entity.Antro
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("dt desc").All(&antro)
	return antro, err
}

func (r repository) GetOne(ctx context.Context, id string) (entity.Antro, error) {
	var antro entity.Antro
	err := r.db.With(ctx).Select().Model(id, &antro)
	return antro, err
}

func (r repository) Delete(ctx context.Context, id string) error {
	antro, err := r.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&antro).Delete()
}

func (r repository) Create(ctx context.Context, antro entity.Antro) error {
	return r.db.With(ctx).Model(&antro).Insert()
}

func (r repository) Update(ctx context.Context, antro entity.Antro) error {
	return r.db.With(ctx).Model(&antro).Update()
}

func (r repository) GetProfile(ctx context.Context, id string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

func (r repository) GetByLogin(ctx context.Context, login string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"login": login}).One(&user)
	return user, err
}
