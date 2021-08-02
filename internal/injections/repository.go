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
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).All(&antro)
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
