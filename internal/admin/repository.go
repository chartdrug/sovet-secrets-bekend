package admin

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

type Repository interface {
	GetAll(ctx context.Context) ([]entity.UsersAdmin, error)
	CreateFeedBack(ctx context.Context, antro entity.Feedback) error
	GetFeedBack(ctx context.Context, id string) (entity.Feedback, error)
	GetAllFeedback(ctx context.Context) ([]entity.Feedback, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetAll(ctx context.Context) ([]entity.UsersAdmin, error) {
	var users []entity.UsersAdmin

	err := r.db.With(ctx).NewQuery("select u.id, u.login, u.admin,u.blocked, u.blocked_at, u.email, u.date_registered, u.date_lastlogin, u.sex, u.birthday, u.type_sports, COALESCE(u.country,'-') as country, COALESCE(u.region,'-') as region, COALESCE(u.city,'-') as city, COALESCE(a.count_antro,0) as count_antro, COALESCE(i.count_injection,0) as count_injection " +
		"from (select distinct on (id) id, login, admin,blocked, blocked_at, email, date_registered, date_lastlogin, sex, birthday, type_sports, date_event, country, region, city from users u " +
		"left join history_login h on u.id = h.id_user " +
		"order by id, date_event desc) u " +
		"left join (select owner, count(1) as count_antro from antro group by owner) a on a.owner = u.id " +
		"left join (select owner, count(1) as count_injection from injection group by owner) i on i.owner = u.id " +
		"order by date_lastlogin desc").All(&users)

	return users, err
}

func (r repository) CreateFeedBack(ctx context.Context, antro entity.Feedback) error {
	return r.db.With(ctx).Model(&antro).Insert()
}

func (r repository) GetFeedBack(ctx context.Context, id string) (entity.Feedback, error) {
	var user entity.Feedback
	err := r.db.With(ctx).Select().Model(id, &user)
	return user, err
}

func (r repository) GetAllFeedback(ctx context.Context) ([]entity.Feedback, error) {
	var feedback []entity.Feedback

	err := r.db.With(ctx).NewQuery("SELECT id, owner, dt, email, name, feedback, location FROM feedback order by dt desc").All(&feedback)

	return feedback, err
}
