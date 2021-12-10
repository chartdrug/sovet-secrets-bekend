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
	GetDose(ctx context.Context, id string) ([]entity.Injection_Dose, error)
	GetAllDose(ctx context.Context, owner string) ([]entity.Injection_Dose, error)
	GetOne(ctx context.Context, id string) (entity.Injection, error)
	Delete(ctx context.Context, id string) error
	GetOneDose(ctx context.Context, id string) (entity.Injection_Dose, error)
	DeleteDose(ctx context.Context, idDose string) error
	SaveConcentration(ctx context.Context, concentration []entity.Concentration) error
	DeleteConcentration(ctx context.Context, owner string, id string) error
	CreateInjection(ctx context.Context, injection entity.Injection) error
	CreateInjectionDose(ctx context.Context, injectionDose entity.Injection_Dose) error
	GetConcentrationDrugs(ctx context.Context, owner string) ([]entity.Concentration, error)
	GetConcentration(ctx context.Context, owner string, drug string) ([]entity.Concentration, error)
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

	vals := []interface{}{}
	for _, row := range concentration {
		vals = append(vals, row.Owner, row.Id_injection, row.Drug, row.Dt, row.C, row.CC, row.CCT, row.CT)
	}

	sqlStr := `INSERT INTO concentration (owner, id_injection, drug, dt, c, cc, cct, ct) VALUES %s`
	sqlStr = ReplaceSQL(sqlStr, "(?, ?, ?, ?, ?, ?, ?, ?)", len(concentration))

	//Prepare and execute the statement
	stmt, _ := r.db2.Prepare(sqlStr)
	_, err := stmt.Exec(vals...)

	return err

}

func (r repository) GetConcentrationDrugs(ctx context.Context, owner string) ([]entity.Concentration, error) {
	var concentrationDrugs []entity.Concentration
	err := r.db.With(ctx).Select("drug").Where(dbx.HashExp{"owner": owner}).Distinct(true).All(&concentrationDrugs)
	return concentrationDrugs, err
}

func (r repository) GetConcentration(ctx context.Context, owner string, drug string) ([]entity.Concentration, error) {
	var concentration []entity.Concentration
	/*
		var concentration []struct {
			Drug string
			Dt int64
			CCT int64
		}*/
	/*q := r.db.With(ctx).NewQuery("select drug, CAST (round(dt/(1000*60*15))*(1000*60*15) AS BIGINT) as dt, max(CCT) as CCT " +
	"from concentration where owner = '3a56fd3a-e7ee-11eb-ba80-0242ac130004' " +
	"and drug = '00000001-0003-0000-0000-ff00ff00ff00' " +
	"group by drug,round(dt/(1000*60*5))*(1000*60*5) order by drug,dt")*/
	q := r.db.With(ctx).NewQuery("select drug, CAST (round(dt/(1000*60*5))*(1000*60*5) AS BIGINT) as dt from concentration where id_injection = '029ef562-c5f5-48a8-9ebb-eabd148e8c70'\n")
	err := q.All(&concentration)
	return concentration, err
	//return r.db.With(ctx).NewQuery("wqw")
	//.Delete("concentration",dbx.HashExp{"id_injection": id}).All(&concentration)
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
