package historyupdate

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context) ([]entity.HistoryUpdate, error)
	Create(ctx context.Context, historyUpdate entity.HistoryUpdate) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context) ([]entity.HistoryUpdate, error) {
	var historyupdate []entity.HistoryUpdate
	err := r.db.With(ctx).Select().OrderBy("created_at desc").All(&historyupdate)
	return historyupdate, err
}

func (r repository) Create(ctx context.Context, historyUpdate entity.HistoryUpdate) error {
	return r.db.With(ctx).Model(&historyUpdate).Insert()
}
