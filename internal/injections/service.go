package injections

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, owner string) ([]Injection, error)
	//Delete(ctx context.Context, id string) (Injection, error)
	//Create(ctx context.Context, input CreateAntroRequest, owner string) (Injection, error)
}

// Album represents the data about an album.
type Injection struct {
	entity.Injection
}

type CreateAntroRequest struct {
	Dt                time.Time `json:"dt"`
	General_age       int       `json:"general_age"`
	General_hip       float32   `json:"general_hip"`
	General_height    float32   `json:"general_height"`
	General_leglen    float32   `json:"general_leglen"`
	General_weight    float32   `json:"general_weight"`
	General_handlen   float32   `json:"general_handlen"`
	General_shoulders float32   `json:"general_shoulders"`
	Notes             string    `json:"notes"`
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context, owner string) ([]Injection, error) {
	items, err := s.repo.Get(ctx, owner)
	if err != nil {
		return nil, err
	}
	result := []Injection{}
	for _, item := range items {
		result = append(result, Injection{item})
	}
	return result, nil
}
