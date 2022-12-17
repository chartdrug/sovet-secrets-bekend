package historyupdate

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"time"
)

type Service interface {
	Get(ctx context.Context) ([]entity.HistoryUpdate, error)
	Create(ctx context.Context, input entity.HistoryUpdate) ([]entity.HistoryUpdate, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) Create(ctx context.Context, req entity.HistoryUpdate) ([]entity.HistoryUpdate, error) {

	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.HistoryUpdate{
		ID:            id,
		CreatedAt:     now,
		DescriptionRu: req.DescriptionRu,
		DescriptionEn: req.DescriptionEn,
	})

	if err != nil {
		return nil, err
	}
	return s.Get(ctx)
}

func (s service) Get(ctx context.Context) ([]entity.HistoryUpdate, error) {
	historyupdate, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return historyupdate, nil
}
