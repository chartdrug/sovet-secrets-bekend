package spr

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Spr, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, sprName string) ([]entity.Spr, error) {
	var spr []entity.Spr
	err := r.db.With(ctx).NewQuery("select *  from spr where name_spr = {:name_spr}").Bind(dbx.Params{"name_spr": sprName}).All(&spr)
	return spr, err
}
