package injections

import (
	"context"
	"database/sql"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"strconv"
	"strings"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id string) ([]entity.Injection, error)
	GetByDate(ctx context.Context, id string, sd time.Time, ed time.Time) ([]entity.Injection, error)
	GetAll(ctx context.Context, id string) ([]entity.Injection, error)
	GetInjectionByCourse(ctx context.Context, owner string, id string) ([]entity.Injection, error)
	GetDose(ctx context.Context, id string) ([]entity.Injection_Dose, error)
	GetAllDose(ctx context.Context, owner string) ([]entity.Injection_Dose, error)
	GetAllDoseCourse(ctx context.Context, owner string, id string) ([]entity.Injection_Dose, error)
	GetOne(ctx context.Context, id string) (entity.Injection, error)
	Delete(ctx context.Context, id string) error
	GetOneDose(ctx context.Context, id string) (entity.Injection_Dose, error)
	DeleteDose(ctx context.Context, idDose string) error
	SaveConcentration(ctx context.Context, concentration []entity.Concentration) error
	SaveConcentration2(ctx context.Context, concentration []entity.Concentration)
	DeleteConcentration(ctx context.Context, owner string, id string) error
	CreateInjection(ctx context.Context, injection entity.Injection) error
	CreateInjectionDose(ctx context.Context, injectionDose entity.Injection_Dose) error
	GetConcentrationDrugs(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Concentration, error)
	GetConcentrationDrugsCourse(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Concentration, error)
	GetCountCalcProcessInjection(ctx context.Context, owner string) (int, error)
	GetConcentration(ctx context.Context, owner string, drug string, sd time.Time, ed time.Time) ([]entity.Concentration, error)
	GetConcentration2(ctx context.Context, owner string, id string) ([]entity.Concentration, error)
	UpdateInjection(ctx context.Context, injection entity.Injection) error
	GetForBloodVolume(ctx context.Context, owner string) ([]entity.Antro, error)
	GetInjectionLimit(ctx context.Context) ([]entity.Injection, error)
	UpdateInjectionCalc(ctx context.Context, id string) error
	GetCourse(ctx context.Context, id string) (entity.Course, error)
}

type repository struct {
	db     *dbcontext.DB
	db2    *sql.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, db2 *sql.DB, logger log.Logger) Repository {
	return repository{db, db2, logger}
}

func (r repository) Get(ctx context.Context, owner string) ([]entity.Injection, error) {
	var injection []entity.Injection
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("dt desc").All(&injection)
	return injection, err
}

func (r repository) GetByDate(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Injection, error) {
	var injection []entity.Injection
	//err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("dt desc").All(&injection)
	err := r.db.With(ctx).NewQuery("select * from injection where owner = {:owner} and dt >= {:sd} and dt <= {:ed} order by dt desc").Bind(dbx.Params{"owner": owner, "sd": sd, "ed": ed}).All(&injection)
	return injection, err
}

func (r repository) GetAll(ctx context.Context, owner string) ([]entity.Injection, error) {
	var injection []entity.Injection
	//err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).OrderBy("dt desc").All(&injection)
	err := r.db.With(ctx).NewQuery("select * from injection where owner = {:owner} order by dt desc").Bind(dbx.Params{"owner": owner}).All(&injection)
	return injection, err
}

func (r repository) GetInjectionByCourse(ctx context.Context, owner string, id string) ([]entity.Injection, error) {
	var injection []entity.Injection
	err := r.db.With(ctx).NewQuery("select * from injection where owner = {:owner} and (id = {:id} or course = {:id}) order by dt desc").Bind(dbx.Params{"owner": owner, "id": id}).All(&injection)
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

func (r repository) GetAllDose(ctx context.Context, owner string) ([]entity.Injection_Dose, error) {
	var injectionDose []entity.Injection_Dose
	err := r.db.With(ctx).NewQuery("select *  from injection_dose where id_injection in (select id from injection where owner = {:owner})").Bind(dbx.Params{"owner": owner}).All(&injectionDose)
	return injectionDose, err
}

func (r repository) GetAllDoseCourse(ctx context.Context, owner string, id string) ([]entity.Injection_Dose, error) {
	var injectionDose []entity.Injection_Dose
	err := r.db.With(ctx).NewQuery("select *  from injection_dose where id_injection in (select id from injection where owner = {:owner} and (id = {:id} or course = {:id}))").Bind(dbx.Params{"owner": owner, "id": id}).All(&injectionDose)
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

func (r repository) SaveConcentration(ctx context.Context, concentration []entity.Concentration) error {

	for j := 0; j < len(concentration); j = j + 5000 {
		vals := []interface{}{}
		count := 0
		//сохраняем каждую 15ю минуту
		//fmt.Println(concentration[0].Drug)
		//fmt.Println(len(vals))
		for i := j; i < len(concentration) && i < (j+5000); i++ {
			count++
			vals = append(vals, concentration[i].Owner, concentration[i].Id_injection, concentration[i].Drug,
				concentration[i].Dt, concentration[i].C, concentration[i].CC, concentration[i].CCT, concentration[i].CT)

		}
		//fmt.Println(count)
		//fmt.Println(len(vals))
		sqlStr := `INSERT INTO concentration (owner, id_injection, drug, dt, c, cc, cct, ct) VALUES %s`
		sqlStr = ReplaceSQL(sqlStr, "(?, ?, ?, ?, ?, ?, ?, ?)", count)
		//fmt.Println(sqlStr)
		//Prepare and execute the statement
		stmt, _ := r.db2.Prepare(sqlStr)
		_, err := stmt.Exec(vals...)

		if err != nil {
			return err
		}

	}
	return nil
}

func (r repository) SaveConcentration2(ctx context.Context, concentration []entity.Concentration) {
	logger := r.logger.With(ctx)
	for j := 0; j < len(concentration); j = j + 5000 {
		vals := []interface{}{}
		count := 0
		for i := j; i < len(concentration) && i < (j+5000); i++ {
			count++
			vals = append(vals, concentration[i].Owner, concentration[i].Id_injection, concentration[i].Drug,
				concentration[i].Dt, concentration[i].C, concentration[i].CC, concentration[i].CCT, concentration[i].CT)

		}
		sqlStr := `INSERT INTO concentration (owner, id_injection, drug, dt, c, cc, cct, ct) VALUES %s`
		sqlStr = ReplaceSQL(sqlStr, "(?, ?, ?, ?, ?, ?, ?, ?)", count)
		stmt, _ := r.db2.Prepare(sqlStr)
		_, err := stmt.Exec(vals...)

		if err != nil {
			logger.Error("SaveConcentration2 error to save = " + err.Error())
		}

	}
}

func (r repository) GetConcentrationDrugs(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Concentration, error) {
	var concentrationDrugs []entity.Concentration
	//err := r.db.With(ctx).Select("drug").Where(dbx.HashExp{"owner": owner}).Distinct(true).All(&concentrationDrugs)
	err := r.db.With(ctx).NewQuery("select distinct drug " +
		"from concentration where owner = {:owner} and dt >= {:sd} and dt <= {:ed} and id_injection in (select id from injection where owner = {:owner})").Bind(dbx.Params{"owner": owner, "sd": sd.Unix() * 1000, "ed": ed.Unix() * 1000}).All(&concentrationDrugs)
	return concentrationDrugs, err
}

func (r repository) GetConcentrationDrugsCourse(ctx context.Context, owner string, sd time.Time, ed time.Time) ([]entity.Concentration, error) {
	var concentrationDrugs []entity.Concentration
	err := r.db.With(ctx).NewQuery("select distinct drug " +
		"from concentration where owner = {:owner} and id_injection in (select id from injection where owner = {:owner} and dt >= {:sd} and dt <= {:ed})").Bind(dbx.Params{"owner": owner, "sd": sd.Unix() * 1000, "ed": ed.Unix() * 1000}).All(&concentrationDrugs)
	return concentrationDrugs, err
}

func (r repository) GetCountCalcProcessInjection(ctx context.Context, owner string) (int, error) {
	var count int
	//err := r.db.With(ctx).Select("drug").Where(dbx.HashExp{"owner": owner}).Distinct(true).All(&concentrationDrugs)
	err := r.db.With(ctx).NewQuery("select count(1) as count from injection where owner = {:owner} " +
		"and calc=false and delete = false").Bind(dbx.Params{"owner": owner}).Row(&count)
	return count, err
}

func (r repository) GetConcentration(ctx context.Context, owner string, drug string, sd time.Time, ed time.Time) ([]entity.Concentration, error) {
	var concentration []entity.Concentration
	/*q := r.db.With(ctx).NewQuery("select drug, CAST (round(dt/(1000*60*15))*(1000*60*15) AS BIGINT) as dt, max(CCT) as CCT " +
	"from concentration where owner = {:owner} and dt >= {:sd} and dt <= {:ed} and id_injection in (select id from injection where owner = {:owner})" +
	"group by 1,2 order by 2").Bind(dbx.Params{"owner": owner, "sd": sd.Unix() * 1000, "ed": ed.Unix() * 1000})*/
	q := r.db.With(ctx).NewQuery("select drug, dt, sum(CCT) as CCT from(select drug, a.id_injection, CAST (round(dt/(1000*60*60))*(1000*60*60) AS BIGINT) as dt, max(CCT) as CCT " +
		"from concentration a where owner = {:owner} and dt >= {:sd} and dt <= {:ed} " +
		"and id_injection in (select id from injection where owner = {:owner}) " +
		"group by 1,2,3) a " +
		"group by 1,2 order by 2").Bind(dbx.Params{"owner": owner, "sd": sd.Unix() * 1000, "ed": ed.Unix() * 1000})
	err := q.All(&concentration)
	return concentration, err
}

func (r repository) GetConcentration2(ctx context.Context, owner string, id string) ([]entity.Concentration, error) {
	var concentration []entity.Concentration
	q := r.db.With(ctx).NewQuery("select drug, dt, sum(CCT) as CCT from(select drug, a.id_injection, CAST (round(dt/(1000*60*60))*(1000*60*60) AS BIGINT) as dt, max(CCT) as CCT " +
		"from concentration a where owner = {:owner} " +
		"and id_injection in (select id from injection where owner = {:owner} and (id = {:id} or course = {:id})) " +
		"group by 1,2,3) a " +
		"group by 1,2 order by 2").Bind(dbx.Params{"owner": owner, "id": id})
	err := q.All(&concentration)
	return concentration, err
}

func (r repository) DeleteConcentration(ctx context.Context, owner string, id string) error {
	var concentration []entity.Concentration
	return r.db.With(ctx).Delete("concentration", dbx.HashExp{"id_injection": id, "owner": owner}).All(&concentration)

}

func (r repository) CreateInjectionDose(ctx context.Context, injectionDose entity.Injection_Dose) error {
	return r.db.With(ctx).Model(&injectionDose).Insert()
}

func ReplaceSQL(stmt, pattern string, len int) string {
	pattern += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(stmt, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		stmt = strings.Replace(stmt, "?", param, 1)
	}
	return strings.TrimSuffix(stmt, ",")
}

func (r repository) UpdateInjection(ctx context.Context, injection entity.Injection) error {
	fmt.Println(injection.Calc)
	_, err := r.db.With(ctx).Update("injection", dbx.Params{"calc": injection.Calc}, dbx.NewExp("id={:id}", dbx.Params{"id": injection.ID})).Execute()
	return err
}

func (r repository) GetForBloodVolume(ctx context.Context, owner string) ([]entity.Antro, error) {
	var antro []entity.Antro
	err := r.db.With(ctx).NewQuery("(select dt,general_weight,result_fat from antro where owner={:owner} /*and dt<$2::date*/ limit 1) " +
		"union " +
		"(select dt,general_weight, result_fat from antro where owner={:owner}  /*and dt>=$2::date and dt<=$3::date*/) " +
		"order by dt desc").Bind(dbx.Params{"owner": owner}).All(&antro)
	return antro, err
}

func (r repository) GetInjectionLimit(ctx context.Context) ([]entity.Injection, error) {
	var injection []entity.Injection
	err := r.db.With(ctx).NewQuery("select * from injection where calc = false and delete = false limit 1").All(&injection)
	return injection, err
}

func (r repository) UpdateInjectionCalc(ctx context.Context, id string) error {
	_, err := r.db.With(ctx).Update("injection", dbx.Params{"calc": "true"}, dbx.HashExp{"id": id}).Execute()
	return err
}

func (r repository) GetCourse(ctx context.Context, id string) (entity.Course, error) {
	var course entity.Course
	err := r.db.With(ctx).Select().Model(id, &course)
	return course, err
}
