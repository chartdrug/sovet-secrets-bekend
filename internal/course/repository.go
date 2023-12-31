package course

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Course, error)
	GetOne(ctx context.Context, id string) (entity.Course, error)
	GetByDate(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Course, error)
	Delete(ctx context.Context, id string) error
	Create(ctx context.Context, album entity.Course) error
	Update(ctx context.Context, album entity.Course) error
	UntieCourse(ctx context.Context, courseId string) error
	СontactCourse(ctx context.Context, courseId string, owner string, sd time.Time, ed time.Time) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Get(ctx context.Context, owner string) ([]entity.Course, error) {
	var course []entity.Course
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("course_start desc").All(&course)
	return course, err
}

func (r repository) GetByDate(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Course, error) {
	var course []entity.Course
	err := r.db.With(ctx).NewQuery("select * from course " +
		"where owner = {:owner} and ((course_start <= {:sd} and course_end >= {:sd}) or (course_start <= {:ed} and course_end >= {:ed}))").Bind(dbx.Params{"owner": owner, "sd": sd, "ed": ed}).All(&course)
	return course, err
}

func (r repository) GetOne(ctx context.Context, id string) (entity.Course, error) {
	var course entity.Course
	err := r.db.With(ctx).Select().Model(id, &course)
	return course, err
}

func (r repository) Delete(ctx context.Context, id string) error {
	course, err := r.GetOne(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&course).Delete()
}

func (r repository) UntieCourse(ctx context.Context, courseId string) error {
	_, err := r.db.With(ctx).Update("injection", dbx.Params{"course": "369cf9b0-f652-11ec-b939-0242ac120002"}, dbx.HashExp{"course": courseId}).Execute()
	return err
}

func (r repository) СontactCourse(ctx context.Context, courseId string, owner string, sd time.Time, ed time.Time) error {
	_, err := r.db.With(ctx).NewQuery("update injection set course = {:course} where owner = {:owner} and dt >= {:sd} and DATE_TRUNC('day',dt) <= {:ed}").
		Bind(dbx.Params{"course": courseId, "owner": owner, "sd": sd, "ed": ed}).Execute()
	return err

}

func (r repository) Create(ctx context.Context, course entity.Course) error {
	return r.db.With(ctx).Model(&course).Insert()
}

func (r repository) Update(ctx context.Context, course entity.Course) error {
	return r.db.With(ctx).Model(&course).Update()
}
