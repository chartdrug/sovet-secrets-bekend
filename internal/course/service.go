package course

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, owner string) ([]entity.Course, error)
	Delete(ctx context.Context, id string) (entity.Course, error)
	Create(ctx context.Context, input entity.Course, owner string) (entity.Course, error)
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
func (s service) Get(ctx context.Context, owner string) ([]entity.Course, error) {
	items, err := s.repo.Get(ctx, owner)
	if err != nil {
		return nil, err
	}
	return items, nil
	/*
		result := []entity.Course{}
		for _, item := range items {
			result = append(result, entity.Course{item})
		}
		return result, nil*/
}

func (s service) GetOne(ctx context.Context, id string) (entity.Course, error) {
	course, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return entity.Course{}, err
	}
	return course, nil
}

func (s service) Delete(ctx context.Context, id string) (entity.Course, error) {
	course, err := s.GetOne(ctx, id)
	if err != nil {
		return entity.Course{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return entity.Course{}, err
	}
	return course, nil
}

func (s service) Create(ctx context.Context, req entity.Course, owner string) (entity.Course, error) {

	logger := s.logger.With(ctx, "owner", owner)

	logger.Infof("Create start ")

	id := entity.GenerateID()

	logger.Infof("Start get profile")

	println(id)

	return entity.Course{}, nil
	/*
		profile, err := s.repo.GetProfile(ctx, owner)
		if err != nil {
			return Antro{}, err
		}
		logger.Infof("End get profile")

		var General_age int

		year, month, day, hour, min, sec := utils.Diff(profile.Birthday, req.Dt)
		_, _, _, _, _ = month, day, hour, min, sec

		logger.Infof("General_age=" + strconv.Itoa(year))

		General_age = year

		var s3 float64
		var fat float64
		var energy float64
		var nofat float64

		if (profile.Sex) == "M" {
			s3 = req.Fold_chest + req.Fold_belly + req.Fold_hip_front
			//todo нужно проверить формулу
			fat = 495/(1.109380-0.0008267*s3+0.0000016*s3*s3-0.0002574*float64(General_age)) - 450
			nofat = req.General_weight - (req.General_weight/100)*fat
			//energy=66+(13.7*req.General_weight)+(5*req.General_height)-(6.8*General_age);
			energy = 370 + (21.6 * nofat)
		} else {
			s3 = req.Fold_shoulder_rear + req.Fold_anterrior_iliac + req.Fold_hip_front
			fat = 495/(1.099421-0.0009929*s3+0.0000023*s3*s3-0.0001392*float64(General_age)) - 450
			nofat = req.General_weight - (req.General_weight/100)*fat
			//energy=655+(9.6*General_age)+(1.8*req.General_height)-(4.7*General_age);
			energy = 370 + (21.6 * nofat)
		}
		if fat > 100 {
			return Antro{}, errors.BadRequest("МДЖ>100%")
		}

		e := entity.Antro{
			ID:                id,
			Owner:             owner,
			Dt:                req.Dt,
			General_age:       General_age,
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
			Result_fat:    fat,
			Result_nofat:  nofat,
			Result_energy: energy,
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

	*/
}
