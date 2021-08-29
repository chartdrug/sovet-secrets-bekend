package injections

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Injection, error)
	GetDose(ctx context.Context, id string) ([]entity.Injection_Dose, error)
	GetOne(ctx context.Context, id string) (entity.Injection, error)
	Delete(ctx context.Context, id string) error
	GetOneDose(ctx context.Context, id string) (entity.Injection_Dose, error)
	DeleteDose(ctx context.Context, idDose string) error
	CreateInjection(ctx context.Context, injection entity.Injection) error
	CreateInjectionDose(ctx context.Context, injectionDose entity.Injection_Dose) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, owner string) ([]entity.Injection, error) {
	var injection []entity.Injection
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).All(&injection)
	return injection, err
}

func (r repository) GetDose(ctx context.Context, id_injection string) ([]entity.Injection_Dose, error) {
	var injection_dose []entity.Injection_Dose
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"id_injection": id_injection}).All(&injection_dose)
	return injection_dose, err
}

func (r repository) GetOne(ctx context.Context, id string) (entity.Injection, error) {
	var injection entity.Injection
	err := r.db.With(ctx).Select().Model(id, &injection)
	return injection, err
}

func (r repository) GetOneDose(ctx context.Context, id string) (entity.Injection_Dose, error) {
	var injectionDose entity.Injection_Dose
	err := r.db.With(ctx).Select().Model(id, &injectionDose)
	return injectionDose, err
}

func (r repository) Delete(ctx context.Context, id string) error {
	injection, err := r.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&injection).Delete()
}

func (r repository) DeleteDose(ctx context.Context, idDose string) error {
	injection, err := r.GetOneDose(ctx, idDose)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&injection).Delete()
}

func (r repository) CreateInjection(ctx context.Context, injection entity.Injection) error {
	return r.db.With(ctx).Model(&injection).Insert()
}

func (r repository) CreateInjectionDose(ctx context.Context, injectionDose entity.Injection_Dose) error {
	return r.db.With(ctx).Model(&injectionDose).Insert()
}
