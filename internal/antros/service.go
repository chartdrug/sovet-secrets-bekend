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
	Id                string    `json:"id"`
	Dt                time.Time `json:"dt"`
	General_age       int       `json:"general_age"`
	General_hip       float32   `json:"general_hip"`
	General_height    float32   `json:"general_height"`
	General_leglen    float32   `json:"general_leglen"`
	General_weight    float32   `json:"general_weight"`
	General_handlen   float32   `json:"general_handlen"`
	General_shoulders float32   `json:"general_shoulders"`

	Fold_anterrior_iliac float32 `json:"fold_anterrior_iliac"`
	Fold_back            float32 `json:"fold_back"`
	Fold_belly           float32 `json:"fold_belly"`
	Fold_chest           float32 `json:"fold_chest"`
	Fold_forearm         float32 `json:"fold_forearm"`
	Fold_hip_front       float32 `json:"fold_hip_front"`
	Fold_hip_inside      float32 `json:"fold_hip_inside"`
	Fold_hip_rear        float32 `json:"fold_hip_rear"`
	Fold_hip_side        float32 `json:"fold_hip_side"`
	Fold_scapula         float32 `json:"fold_scapula"`
	Fold_shin            float32 `json:"fold_shin"`
	Fold_shoulder_front  float32 `json:"fold_shoulder_front"`
	Fold_shoulder_rear   float32 `json:"fold_shoulder_rear"`
	Fold_waist_side      float32 `json:"fold_waist_side"`
	Fold_wrist           float32 `json:"fold_wrist"`
	Fold_xiphoid         float32 `json:"fold_xiphoid"`

	Notes string `json:"notes"`
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

	e := entity.Antro{
		ID:                id,
		Owner:             owner,
		Dt:                req.Dt,
		General_age:       40, //TODO сделать расчёт возраста
		General_hip:       req.General_hip,
		General_height:    req.General_height,
		General_leglen:    req.General_leglen,
		General_weight:    req.General_weight,
		General_handlen:   req.General_handlen,
		General_shoulders: req.General_shoulders,

		Fold_anterrior_iliac: req.Fold_anterrior_iliac,
		Fold_back:            req.Fold_back,
		Fold_belly:           req.Fold_belly,
		Fold_chest:           req.Fold_chest,
		Fold_forearm:         req.Fold_forearm,
		Fold_hip_front:       req.Fold_hip_front,
		Fold_hip_inside:      req.Fold_hip_inside,
		Fold_hip_rear:        req.Fold_hip_rear,
		Fold_hip_side:        req.Fold_hip_side,
		Fold_scapula:         req.Fold_scapula,
		Fold_shin:            req.Fold_shin,
		Fold_shoulder_front:  req.Fold_shoulder_front,
		Fold_shoulder_rear:   req.Fold_shoulder_rear,
		Fold_waist_side:      req.Fold_waist_side,
		Fold_wrist:           req.Fold_wrist,
		Fold_xiphoid:         req.Fold_xiphoid,

		Notes:         req.Notes,
		Result_fat:    result_fat,
		Result_nofat:  result_nofat,
		Result_energy: result_energy,
	}

	if req.Id != "" {
		e.ID = req.Id
		err := s.repo.Update(ctx, e)
		if err != nil {
			return Antro{}, err
		}
	} else {
		err := s.repo.Create(ctx, e)
		if err != nil {
			return Antro{}, err
		}
	}

	return s.GetOne(ctx, e.ID)
}
