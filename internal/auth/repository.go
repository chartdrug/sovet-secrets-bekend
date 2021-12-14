package auth

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// Repository encapsulates the logic to access albums from the data source.
type Repository interface {
	// Get returns the album with the specified album ID.
	Get(ctx context.Context, login string, password string) (entity.Users, error)
	UpdateTimeLastLogin(ctx context.Context, id string) error
}

// repository persists albums in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new album repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the album with the specified ID from the database.
func (r repository) Get(ctx context.Context, login string, password string) (entity.Users, error) {
	var user entity.Users

	//err := r.db.With(ctx).Select().Bind(dbx.Params{"id": 100}).One(&user)
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"login": login, "passwd": password}).One(&user)
	//.dbx.NewExp("id={:id}", dbx.Params{"id":100})
	//.Bind(dbx.Params{"id": 100}).One(&user)
	//Model(login, &user)
	return user, err
}

func (r repository) UpdateTimeLastLogin(ctx context.Context, id string) error {
	_, err := r.db.With(ctx).Update("users", dbx.Params{"date_lastlogin": "now()"}, dbx.NewExp("id={:id}", dbx.Params{"id": id})).Execute()
	return err
}
