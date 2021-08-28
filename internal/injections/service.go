package injections

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, owner string) ([]InjectionModel, error)
	//Delete(ctx context.Context, id string) (Injection, error)
	Create(ctx context.Context, input CreateInjectionsRequest, owner string) (InjectionModel, error)
}

// Album represents the data about an album.
type InjectionModel struct {
	entity.InjectionModel
}
type Injection struct {
	entity.Injection
}

type Injection_Dose struct {
	entity.Injection_Dose
}

type CreateInjectionsRequest struct {
	entity.InjectionModel
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
func (s service) Get(ctx context.Context, owner string) ([]InjectionModel, error) {
	items, err := s.repo.Get(ctx, owner)
	if err != nil {
		return nil, err
	}
	result := []InjectionModel{}
	for _, item := range items {
		itemsDose, errDose := s.repo.GetDose(ctx, item.ID)
		if errDose != nil {
			return nil, errDose
		}
		//resultDose := []Injection_Dose{}
		//item.Injection_Dose = itemsDose
		resultModel := InjectionModel{}
		resultModel.Injection = item
		resultModel.Injection_Dose = itemsDose
		result = append(result, resultModel)
	}
	return result, nil
}

func (s service) GetOne(ctx context.Context, id string) (InjectionModel, error) {
	injection, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return InjectionModel{}, err
	}

	itemsDose, errDose := s.repo.GetDose(ctx, injection.ID)
	if errDose != nil {
		return InjectionModel{}, errDose
	}

	resultModel := InjectionModel{}
	resultModel.Injection = injection
	resultModel.Injection_Dose = itemsDose

	return resultModel, nil
}

func (s service) Create(ctx context.Context, req CreateInjectionsRequest, owner string) (InjectionModel, error) {

	id := entity.GenerateID()
	// записфваем сначало применение доз
	for _, item := range req.Injection_Dose {
		err := s.repo.CreateInjectionDose(ctx, entity.Injection_Dose{
			ID:           entity.GenerateID(),
			Id_injection: id,
			Dose:         item.Dose,
			Drug:         item.Drug,
			Volume:       item.Volume,
			Solvent:      item.Solvent,
		})
		if err != nil {
			return InjectionModel{}, err
		}
	}
	// создаём сущность самого применения
	err := s.repo.CreateInjection(ctx, entity.Injection{
		ID:    id,
		Owner: owner,
		Dt:    req.Injection.Dt,
		What:  req.Injection.What,
	})
	if err != nil {
		return InjectionModel{}, err
	}
	return s.GetOne(ctx, id)
}
