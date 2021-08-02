package antros

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"math/rand"
	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, owner string) ([]Antro, error)
	Delete(ctx context.Context, id string) (Antro, error)
	Create(ctx context.Context, input CreateAntroRequest, owner string) (Antro, error)
}

// Album represents the data about an album.
type Antro struct {
	entity.Antro
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
func (s service) Get(ctx context.Context, owner string) ([]Antro, error) {
	items, err := s.repo.Get(ctx, owner)
	if err != nil {
		return nil, err
	}
	result := []Antro{}
	for _, item := range items {
		result = append(result, Antro{item})
	}
	return result, nil
}

func (s service) GetOne(ctx context.Context, id string) (Antro, error) {
	antro, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return Antro{}, err
	}
	return Antro{antro}, nil
}

func (s service) Delete(ctx context.Context, id string) (Antro, error) {
	antro, err := s.GetOne(ctx, id)
	if err != nil {
		return Antro{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Antro{}, err
	}
	return antro, nil
}

func (s service) Create(ctx context.Context, req CreateAntroRequest, owner string) (Antro, error) {

	id := entity.GenerateID()

	//TODO сделать
	result_fat := rand.Float32()
	result_nofat := rand.Float32()
	result_energy := rand.Float32()

	err := s.repo.Create(ctx, entity.Antro{
		ID:                id,
		Owner:             owner,
		Dt:                req.Dt,
		General_age:       40, //TODO сделать расчёт возраста
		General_hip:       req.General_hip,
		General_height:    req.General_height,
		General_leglen:    req.General_height,
		General_weight:    req.General_weight,
		General_handlen:   req.General_handlen,
		General_shoulders: req.General_shoulders,
		Notes:             req.Notes,
		Result_fat:        result_fat,
		Result_nofat:      result_nofat,
		Result_energy:     result_energy,
	})
	if err != nil {
		return Antro{}, err
	}
	return s.GetOne(ctx, id)
}
