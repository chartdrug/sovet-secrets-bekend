package injections

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Injection, error)
	//GetOne(ctx context.Context, id string) (entity.Injection, error)
	//Delete(ctx context.Context, id string) error
	//Create(ctx context.Context, album entity.Injection) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, owner string) ([]entity.Injection, error) {
	var injection []entity.Injection
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).All(&injection)
	return injection, err
}

func (r repository) GetOne(ctx context.Context, id string) (entity.Injection, error) {
	var injection entity.Injection
	err := r.db.With(ctx).Select().Model(id, &injection)
	return injection, err
}

func (r repository) Delete(ctx context.Context, id string) error {
	injection, err := r.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&injection).Delete()
}

func (r repository) Create(ctx context.Context, injection entity.Injection) error {
	return r.db.With(ctx).Model(&injection).Insert()
}
