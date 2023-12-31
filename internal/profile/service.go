package profile

import (
	"context"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"strings"
	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, id string) (Users, error)
	GetHistoryLogin(ctx context.Context, id string) ([]entity.HistoryLogin, error)
	Create(ctx context.Context, req CreateUser) (Users, error)
	Update(ctx context.Context, id string, req CreateUser) (Users, error)
	GetByLogin(ctx context.Context, login string) (Users, error)
	GetByEmail(ctx context.Context, email string) (Users, error)
}

// Album represents the data about an album.
type Users struct {
	entity.Users
}

type CreateUser struct {
	Login      string    `json:"login"`
	Pass       string    `json:"pass"`
	Email      string    `json:"email"`
	Sex        string    `json:"sex"`
	Birthday   time.Time `json:"birthday"`
	TypeSports []int     `json:"typesports"`
}

type service struct {
	repo   Repository
	logger log.Logger
}

func (m CreateUser) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Login, validation.Required, validation.Length(1, 30)),
		validation.Field(&m.Pass, validation.Required, validation.Length(1, 30)),
		//validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Email, validation.Required, validation.Length(1, 40)),
		validation.Field(&m.Sex, validation.Required, validation.In("M", "F")),
		validation.Field(&m.Birthday, validation.Required, validation.NotNil),
		validation.Field(&m.TypeSports, validation.Required, validation.NotNil),
	)
}

func (m CreateUser) ValidateUpdate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Login, validation.Required, validation.Length(1, 30)),
		//validation.Field(&m.Pass, validation.Required, validation.Length(1, 30)),
		//validation.Field(&m.Email, validation.Required, is.Email),
		validation.Field(&m.Email, validation.Required, validation.Length(1, 40)),
		validation.Field(&m.Sex, validation.Required, validation.In("M", "F")),
		validation.Field(&m.Birthday, validation.Required, validation.NotNil),
		validation.Field(&m.TypeSports, validation.Required, validation.NotNil),
	)
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the album with the specified the album ID.
func (s service) Get(ctx context.Context, id string) (Users, error) {
	users, err := s.repo.Get(ctx, id)
	if err != nil {
		return Users{}, err
	}
	return Users{users}, nil
}

func (s service) GetHistoryLogin(ctx context.Context, id string) ([]entity.HistoryLogin, error) {
	HistoryLogin, err := s.repo.GetHistoryLogin(ctx, id)
	if err != nil {
		return []entity.HistoryLogin{}, err
	}
	return HistoryLogin, nil
}

func (s service) GetByLogin(ctx context.Context, login string) (Users, error) {
	users, err := s.repo.GetByLogin(ctx, login)
	if err != nil {
		return Users{}, err
	}
	return Users{users}, nil
}

func (s service) GetByEmail(ctx context.Context, email string) (Users, error) {
	users, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return Users{}, err
	}
	return Users{users}, nil
}

func (s service) Create(ctx context.Context, req CreateUser) (Users, error) {

	if err := req.ValidateUpdate(); err != nil {
		return Users{}, err
	}

	id := entity.GenerateID()
	fmt.Println(req.TypeSports)
	err := s.repo.Create(ctx, entity.Users{
		ID:             id,
		Login:          req.Login,
		Passwd:         req.Pass,
		Email:          req.Email,
		DateRegistered: time.Now(),
		Sex:            req.Sex,
		Birthday:       req.Birthday,
		//TypeSports:     "{" + strings.Join(req.TypeSports, ",") + "}",
		TypeSports: "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.TypeSports)), ","), "[]") + "}",
	})
	if err != nil {
		return Users{}, err
	}
	return s.Get(ctx, id)
}

func (s service) Update(ctx context.Context, id string, req CreateUser) (Users, error) {

	if err := req.ValidateUpdate(); err != nil {
		return Users{}, err
	}

	fmt.Println(req.TypeSports)
	err := s.repo.Update(ctx, entity.Users{
		ID:    id,
		Login: req.Login,
		//Passwd:         req.Pass,
		Email: req.Email,
		//DateRegistered: time.Now(),
		Sex:        req.Sex,
		Birthday:   req.Birthday,
		TypeSports: "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(req.TypeSports)), ","), "[]") + "}",
	})
	if err != nil {
		return Users{}, err
	}
	return s.Get(ctx, id)
}
