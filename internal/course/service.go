package course

import (
	"context"
	"fmt"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
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
	//отвязываем всю фарму от этого укола
	if err = s.repo.UntieCourse(ctx, id); err != nil {
		return entity.Course{}, err
	}

	return course, nil
}

func (s service) Create(ctx context.Context, req entity.Course, owner string) (entity.Course, error) {

	logger := s.logger.With(ctx, "owner", owner)

	logger.Infof("Create start ")

	id := entity.GenerateID()

	e := entity.Course{
		Id:           id,
		Owner:        owner,
		Course_start: req.Course_start,
		Course_end:   req.Course_end,
		Type:         req.Type,
		Target:       "{" + req.Target + "}",
		Notes:        req.Notes,
	}

	//ищем вхождения в другой крус
	CourseBydate, errGetByDate := s.repo.GetByDate(ctx, owner, req.Course_start, req.Course_end)
	if errGetByDate != nil {
		utils.SendMailError("Course GetByDate ", e.Id+" - "+errGetByDate.Error())
		return entity.Course{}, errGetByDate
	}

	logger.Info("CourseBydate length" + fmt.Sprintf(":%v", len(CourseBydate)))

	//если нашли один
	if (len(CourseBydate) == 1 && CourseBydate[0].Id != req.Id) || len(CourseBydate) > 1 {
		return entity.Course{}, errors.BadRequest("Найдено пересечение курсов. Пересечений быть не должно!")
	}

	if req.Id != "" {
		e.Id = req.Id
		err := s.repo.Update(ctx, e)
		if err != nil {
			utils.SendMailError("Course Create1 ", e.Id+" - "+err.Error())
			return entity.Course{}, err
		}
		//отвязываем всю фарму от этого укола
		logger.Info("UntieCourse:" + e.Id)
		if err = s.repo.UntieCourse(ctx, e.Id); err != nil {
			return entity.Course{}, err
		}
	} else {
		err := s.repo.Create(ctx, e)
		if err != nil {
			utils.SendMailError("Course Create2", e.Id+" - "+err.Error())
			return entity.Course{}, err
		}
	}

	//всю фарму нужно привязать к новому курсу
	err := s.repo.СontactCourse(ctx, e.Id, e.Owner, e.Course_start, e.Course_end)
	if err != nil {
		utils.SendMailError("СontactCourse", e.Id+" - "+err.Error())
		return entity.Course{}, err
	}

	return s.GetOne(ctx, e.Id)
}
