package profile

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Users, error)
	GetHistoryLogin(ctx context.Context, id string) ([]entity.HistoryLogin, error)
	GetByLogin(ctx context.Context, login string) (entity.Users, error)
	GetByEmail(ctx context.Context, email string) (entity.Users, error)
	Create(ctx context.Context, album entity.Users) error
	Update(ctx context.Context, album entity.Users) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, id string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Model(id, &user)
	//err := r.db.With(ctx).NewQuery("select id, login, passwd, email, date_registered, date_lastlogin, sex, birthday, array_to_json(type_sports::int[]) as type_sports from users where id = {:id}").Bind(dbx.Params{"id": id}).All(&user)
	return user, err
}

func (r repository) GetHistoryLogin(ctx context.Context, id string) ([]entity.HistoryLogin, error) {
	var HistoryLogin []entity.HistoryLogin
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"id_user": id}).
		OrderBy("date_event desc").
		//Offset(int64(offset)).
		//Limit(int64(limit)).
		All(&HistoryLogin)
	return HistoryLogin, err
}

func (r repository) Create(ctx context.Context, antro entity.Users) error {
	return r.db.With(ctx).Model(&antro).Insert()
}
func (r repository) Update(ctx context.Context, user entity.Users) error {
	_, err := r.db.With(ctx).Update("users", dbx.Params{"login": user.Login, "email": user.Email, "sex": user.Sex, "birthday": user.Birthday, "type_sports": user.TypeSports}, dbx.HashExp{"id": user.ID}).Execute()
	return err
}

func (r repository) GetByLogin(ctx context.Context, login string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"login": login}).One(&user)
	return user, err
}

func (r repository) GetByEmail(ctx context.Context, email string) (entity.Users, error) {
	var user entity.Users
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"email": email}).One(&user)
	return user, err
}
