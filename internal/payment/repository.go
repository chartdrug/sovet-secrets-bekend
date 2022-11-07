package payment

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	GetAll(ctx context.Context, owner string) ([]entity.Cryptocloud, error)
	CreateCryptocloud(ctx context.Context, album entity.Cryptocloud) error
	CreateCryptocloudPostback(ctx context.Context, album entity.CryptocloudPostback) error
	UpdateCryproInvoice(ctx context.Context, id string, statusinvoice string, resthttpstatus string) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetAll(ctx context.Context, owner string) ([]entity.Cryptocloud, error) {
	var antro []entity.Cryptocloud
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("dt desc").All(&antro)
	return antro, err
}

func (r repository) CreateCryptocloud(ctx context.Context, album entity.Cryptocloud) error {
	return r.db.With(ctx).Model(&album).Insert()
}

func (r repository) CreateCryptocloudPostback(ctx context.Context, album entity.CryptocloudPostback) error {
	return r.db.With(ctx).Model(&album).Insert()
}

func (r repository) UpdateCryproInvoice(ctx context.Context, id string, statusinvoice string, resthttpstatus string) error {
	_, err := r.db.With(ctx).Update("cryptocloud", dbx.Params{"dtpaym": "now()", "statusinvoice": statusinvoice, "resthttpstatus": resthttpstatus},
		dbx.NewExp("id={:id}", dbx.Params{"id": id})).Execute()
	return err
}
