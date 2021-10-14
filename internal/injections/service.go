package injections

import (
	"context"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"math"

	//"fmt"
	//"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	//"math"

	//"fmt"
	//"github.com/qiangxue/sovet-secrets-bekend/internal/utils"

	//"fmt"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	//"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	//"math"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	Get(ctx context.Context, owner string) ([]InjectionModel, error)
	Getinj(ctx context.Context, id string, owner string) (Points, error)
	Delete(ctx context.Context, id string, owner string) (InjectionModel, error)
	DeleteDose(ctx context.Context, id string, idDose string, owner string) (InjectionModel, error)
	Create(ctx context.Context, input CreateInjectionsRequest, owner string) (InjectionModel, error)
}

// Album represents the data about an album.
type InjectionModel struct {
	entity.InjectionModel
}
type Injection struct {
	entity.Injection
}

type Injection_Dose struct {
	entity.Injection_Dose
}

type CreateInjectionsRequest struct {
	entity.InjectionModel
}

type CreatePoints struct {
	entity.PointsArray
}

type PointValue struct {
	entity.PointValue
}
type Point struct {
	entity.Point
}
type Points struct {
	entity.PointsArray
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
func (s service) Get(ctx context.Context, owner string) ([]InjectionModel, error) {
	items, err := s.repo.Get(ctx, owner)
	if err != nil {
		return nil, err
	}
	result := []InjectionModel{}
	for _, item := range items {
		itemsDose, errDose := s.repo.GetDose(ctx, item.ID)
		if errDose != nil {
			return nil, errDose
		}
		//resultDose := []Injection_Dose{}
		//item.Injection_Dose = itemsDose
		resultModel := InjectionModel{}
		resultModel.Injection = item
		resultModel.Injection_Dose = itemsDose
		result = append(result, resultModel)
	}
	return result, nil
}

func (s service) Getinj(ctx context.Context, id string, owner string) (Points, error) {
	logger := s.logger.With(ctx, "id", id)

	injection, err := s.GetOne(ctx, id)
	if err != nil {
		return Points{}, err
	}

	result := Points{}

	for _, item := range injection.Injection_Dose {
		result.Drugs = append(result.Drugs, item.Drug)
	}

	//TODO формула
	if injection.Injection.What == "W" {
		//табл
		point := entity.Point{}
		logger.Infof("injection.Injection_Dose " + injection.Injection.Dt.String())
		point.Dt = injection.Injection.Dt.Unix() * 1000 // приведение в формат старого кода
		//в исходнике 1633244400000
		// у меня     1633255200
		//point.Dt = 1625079780000

		//резерв под сумму
		point.PointValues = append(point.PointValues, entity.PointValue{})
		point.PointValues[0].Drug = "SUMM"
		point.PointValues[0].C = 0
		point.PointValues[0].CC = 0
		point.PointValues[0].CCT = 0
		point.PointValues[0].CT = 0

		//проходим по всем дозам

		for _, item := range injection.Injection_Dose {

			if item.Drug != "" {
				logger.Infof("injection.Injection_Dose " + item.ID)
			}

			pValue := entity.PointValue{}

			pValue.Drug = item.Drug
			pValue.C = 0
			pValue.CC = 0
			pValue.CCT = 0
			pValue.CT = 0

			point.PointValues[0].C += pValue.C
			point.PointValues[0].CC += pValue.CC
			point.PointValues[0].CCT += pValue.CCT
			point.PointValues[0].CT += pValue.CT

			point.PointValues = append(point.PointValues, pValue)
		}

		result.Points = append(result.Points, point)

		var condition = true
		var count = 1

		drugs := utils.GetDrugs()

		const cStep = 60000 // шаг в расчетах 1 min
		//const rStepDefault=300000 // шаг в результате
		//const ZERO=1e-6 //что считать нулем
		const ZERO = 0.003 //что считать нулем

		// таблетка только через 30 мин растварится
		/* почемуто ломает всю логику
		for i := 1; i <= 30; i++ {
			point1 := point
			point1.Dt += cStep * int64(i)
			result.Points = append(result.Points, point1)
		}*/

		for {

			//создаём объект для графика
			point := entity.Point{}
			//logger.Infof("injection.Injection_Dose " + injection.Injection.Dt.String())

			if count == 1 {
				// колесо работае только через 30 мин
				point.Dt = result.Points[count-1].Dt + cStep*30
			} else {
				// добавляем время
				point.Dt = result.Points[count-1].Dt + cStep
			}

			// создаём под сумму
			point.PointValues = append(point.PointValues, entity.PointValue{})
			point.PointValues[0].Drug = "SUMM"
			point.PointValues[0].C = 0
			point.PointValues[0].CC = 0
			point.PointValues[0].CCT = 0
			point.PointValues[0].CT = 0

			// проходим по всем иньякциям
			//this.CO=state.CO; // концентрация лекарства
			//this.COT=state.COT; // концентрация в пересчете на тестостерон
			var count_Injection_Dose = 0
			condition = false
			for _, item := range injection.Injection_Dose {
				count_Injection_Dose += 1

				pValue := entity.PointValue{}

				pValue.Drug = item.Drug
				//первый проход
				if count == 1 {
					//this.rout[idx]=inj.dose[idx]*0.7*(drug.out/100);
					pValue.C = item.Dose * 0.7 * (drugs[item.Drug].Out / 100)
					//this.routt[idx]=inj.dose[idx]*0.7*(drug.out/100)*(drug.outt/100);
					pValue.CT = item.Dose * 0.7 * (drugs[item.Drug].Out / 100) * drugs[item.Drug].Out / 100

				} else {
					pValue.C = result.Points[count-1].PointValues[count_Injection_Dose].C
					pValue.CT = result.Points[count-1].PointValues[count_Injection_Dose].CT
				}

				//.Infof("drugs[item.ID].Halflife " + string(drugs[item.Drug].Halflife))
				//logger.Infof("pValue.C " + fmt.Sprintf("%f",pValue.C))
				//logger.Infof("item.Dose " + fmt.Sprintf("%f",item.Dose))
				//logger.Infof("drugs[item.ID].Out " + fmt.Sprintf("%f",drugs[item.Drug].Out))
				//logger.Infof("drugs[item.ID].ID " + drugs[item.Drug].ID)
				//logger.Infof("item.ID " + item.Drug)

				//var a=(this.CO[i]*Math.exp(-Math.LN2/this.halflife[i]));
				pValue.C = pValue.C * math.Exp(-math.Ln2/float64(drugs[item.Drug].Halflife))
				if pValue.C < ZERO {
					//logger.Infof("pValue.C " + fmt.Sprintf("%f",pValue.C))
					pValue.C = 0
				} else {
					condition = true
				}

				//this.COT[i]=(this.COT[i]*Math.exp(-Math.LN2/this.halflife[i]));
				pValue.CT = pValue.CT * math.Exp(-math.Ln2/float64(drugs[item.Drug].Halflife))

				pValue.CC = pValue.C
				pValue.CCT = pValue.CT

				point.PointValues[0].C += pValue.C
				point.PointValues[0].CC += pValue.CC
				point.PointValues[0].CCT += pValue.CCT
				point.PointValues[0].CT += pValue.CT

				point.PointValues = append(point.PointValues, pValue)
			}

			result.Points = append(result.Points, point)

			//защита от зацикливания
			if count == 200000 {
				condition = false
			}

			if !condition {
				break
			}
			count += 1

		}

		//цик расчёта

	} else {

	}

	//dtStart := injection.Injection.Dt
	if 1 != 1 {
		for i := 0; i < 3; i++ {

			point := entity.Point{}

			point.Dt = 1625079780000 //+ (i * 10000000)

			//point.Dt = int(dtStart + (i * 10000))
			//var plussI int = i * 10
			//point.Dt = int(dtStart.Add(time.Minute * plussI).Unix())

			//next_time:= cur_time.Add(time.Hour * 2 + time.Minute * 1+ time.Second * 21)
			//fmt.Printf("current time is :%s\n", cur_time )
			//fmt.Printf("calculated time is :%s", next_time)

			point.PointValues = append(point.PointValues, entity.PointValue{})

			point.PointValues[0].Drug = "SUMM"
			point.PointValues[0].C = 0
			point.PointValues[0].CC = 0
			point.PointValues[0].CCT = 0
			point.PointValues[0].CT = 0

			for _, item := range injection.Injection_Dose {

				if item.Drug != "" {
					logger.Infof("injection.Injection_Dose " + item.ID)
				}

				pValue := entity.PointValue{}

				pValue.Drug = item.Drug
				pValue.C = 0.0028981088472945313
				pValue.CC = 0.0028981088472945313
				pValue.CCT = 0.0009991573645662046
				pValue.CT = 0.0028981088472945313

				point.PointValues[0].C += pValue.C
				point.PointValues[0].CC += pValue.CC
				point.PointValues[0].CCT += pValue.CCT
				point.PointValues[0].CT += pValue.CT

				point.PointValues = append(point.PointValues, pValue)
			}

			result.Points = append(result.Points, point)

		}
	}

	return result, nil
}

func (s service) GetOne(ctx context.Context, id string) (InjectionModel, error) {
	logger := s.logger.With(ctx, "id", id)

	injection, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return InjectionModel{}, err
	}

	logger.Infof("GetOne sucsess " + injection.ID)

	itemsDose, errDose := s.repo.GetDose(ctx, injection.ID)
	if errDose != nil {
		return InjectionModel{}, errDose
	}

	resultModel := InjectionModel{}
	resultModel.Injection = injection
	resultModel.Injection_Dose = itemsDose

	return resultModel, nil
}

func (s service) Create(ctx context.Context, req CreateInjectionsRequest, owner string) (InjectionModel, error) {

	id := entity.GenerateID()
	// записфваем сначало применение доз
	for _, item := range req.Injection_Dose {
		err := s.repo.CreateInjectionDose(ctx, entity.Injection_Dose{
			ID:           entity.GenerateID(),
			Id_injection: id,
			Dose:         item.Dose,
			Drug:         item.Drug,
			Volume:       item.Volume,
			Solvent:      item.Solvent,
		})
		if err != nil {
			return InjectionModel{}, err
		}
	}
	// создаём сущность самого применения
	err := s.repo.CreateInjection(ctx, entity.Injection{
		ID:    id,
		Owner: owner,
		Dt:    req.Injection.Dt,
		What:  req.Injection.What,
	})
	if err != nil {
		return InjectionModel{}, err
	}
	return s.GetOne(ctx, id)
}

func (s service) Delete(ctx context.Context, id string, owner string) (InjectionModel, error) {
	injection, err := s.GetOne(ctx, id)
	if err != nil {
		return InjectionModel{}, err
	}

	if injection.Injection.Owner != owner {
		return InjectionModel{}, errors.NotFound("The requested Injection was not found")
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return InjectionModel{}, err
	}

	for _, item := range injection.Injection_Dose {
		if err = s.repo.DeleteDose(ctx, item.ID); err != nil {
			return InjectionModel{}, err
		}
	}

	return injection, nil
}

func (s service) DeleteDose(ctx context.Context, id string, idDose string, owner string) (InjectionModel, error) {
	//logger := s.logger.With(ctx, "id", id)
	injection, err := s.GetOne(ctx, id)
	if err != nil {
		return InjectionModel{}, err
	}

	if injection.Injection.Owner != owner {
		return InjectionModel{}, errors.NotFound("The requested Injection was not found")
	}

	var existsDose bool
	existsDose = false

	for _, item := range injection.Injection_Dose {
		if item.ID == idDose {
			existsDose = true
		}
	}

	if !existsDose {
		return InjectionModel{}, errors.NotFound("The requested Injection dose was not found")
	}

	if err = s.repo.DeleteDose(ctx, idDose); err != nil {
		return InjectionModel{}, err
	}

	injection, err = s.GetOne(ctx, id)
	if err != nil {
		return InjectionModel{}, err
	}
	// если нет больше доз то удаляем всё связку

	//logger.Infof(" len(injection.Injection_Dose)= " +  len(injection.Injection_Dose))
	if len(injection.Injection_Dose) < 0 {
		if err = s.repo.Delete(ctx, id); err != nil {
			return InjectionModel{}, err
		}
		return InjectionModel{}, nil
	}

	return injection, nil
}
