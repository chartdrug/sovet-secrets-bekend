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

type Injection_Dose struct {
	entity.Injection_Dose
}

type CreateInjectionsRequest struct {
	Dt      time.Time       `json:"dt"`
	Course  string          `json:"course"`
	What    string          `json:"what"`
	Dose    map[float32]int `json:"dose"`
	Drug    map[string]int  `json:"drug"`
	Volume  map[float32]int `json:"volume"`
	Solvent map[string]int  `json:"solvent"`
	Points  map[int]int     `json:"points"`
	/*
	   	CREATE TABLE public.injection
	   (
	   	id uuid NOT NULL,
	   	owner uuid NOT NULL,
	   	dt timestamp without time zone NOT NULL DEFAULT ('now'::text)::date,
	   	course uuid,
	   	what character(1) COLLATE pg_catalog."default" NOT NULL DEFAULT '?'::bpchar,
	   	dose double precision[],
	   	drug uuid[],
	   	volume double precision[],
	   	solvent character(1)[] COLLATE pg_catalog."default",
	   	points integer[],
	   	zerodt timestamp without time zone,
	   	hashid uuid,
	   	cutoff integer[],
	   	CONSTRAINT injection_pkey PRIMARY KEY (id)
	   )*/
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
		itemsDose, errDose := s.repo.GetDose(ctx, item.ID)
		if errDose != nil {
			return nil, errDose
		}
		//resultDose := []Injection_Dose{}
		item.Injection_Dose = itemsDose
		result = append(result, Injection{item})
	}
	return result, nil
}
