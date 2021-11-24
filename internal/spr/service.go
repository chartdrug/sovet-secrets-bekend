package spr

import (
	"context"
	//"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	GetSports(ctx context.Context, language string) ([]CreateSpr, error)
}

// Album represents the data about an album.
//type Spr struct {
//	entity.Spr
//}
type ChildSpr struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CreateSpr struct {
	Id    int        `json:"id"`
	Name  string     `json:"name"`
	Child []ChildSpr `json:"child"`
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
func (s service) GetSports(ctx context.Context, language string) ([]CreateSpr, error) {
	//logger := s.logger.With(ctx, "language", language)

	ArrayCreateSpr := []CreateSpr{}
	sprs, err := s.repo.Get(ctx, "sport")
	if err != nil {
		return ArrayCreateSpr, err
	}

	// собираем родительский объект
	for _, spr := range sprs {
		if spr.Id_parent == 0 {
			CSpr := CreateSpr{}
			CSpr.Child = []ChildSpr{}
			CSpr.Id = spr.Id
			if language == "eng" {
				CSpr.Name = spr.Name_eng
			} else {
				CSpr.Name = spr.Name_ru
			}
			//собираем дочернии
			for _, sprChild := range sprs {
				if sprChild.Id_parent == CSpr.Id {
					CSprChild := ChildSpr{}
					CSprChild.Id = sprChild.Id
					if language == "eng" {
						CSprChild.Name = sprChild.Name_eng
					} else {
						CSprChild.Name = sprChild.Name_ru
					}
					CSpr.Child = append(CSpr.Child, CSprChild)
				}
			}
			ArrayCreateSpr = append(ArrayCreateSpr, CSpr)
		}
	}

	return ArrayCreateSpr, nil
}
