package injections

import (
	"context"
	"fmt"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"math"
	"sort"

	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	GetAllDose(ctx context.Context, owner string, sDate string, fDate string) ([]InjectionModel, error)
	Getinj(ctx context.Context, id string, owner string, save bool) (Points, error)
	GetinjOne(ctx context.Context, id string) (InjectionModel, error)
	GetinjReort(ctx context.Context, owner string, sDate string, fDate string) (Points, error)
	Delete(ctx context.Context, id string, owner string) (InjectionModel, error)
	DeleteDose(ctx context.Context, id string, idDose string, owner string) (InjectionModel, error)
	Create(ctx context.Context, input CreateInjectionsRequest, owner string) (InjectionModel, error)
	GetForBloodVolume(ctx context.Context, owner string) ([]entity.Antro, error)
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
func (s service) GetAllDose(ctx context.Context, owner string, sDate string, fDate string) ([]InjectionModel, error) {
	//logger := s.logger.With(ctx, "owner", owner)

	date1, err := time.Parse("2006-01-02", sDate)
	//fmt.Println(date1)
	if err != nil {
		return nil, err
	}

	date2, err := time.Parse("2006-01-02", fDate)
	if err != nil {
		return nil, err
	}
	// добавляем + 1 день
	date2 = date2.AddDate(0, 0, 1)
	//fmt.Println(date2)

	items, err := s.repo.GetByDate(ctx, owner, date1, date2)
	if err != nil {
		return nil, err
	}
	result := []InjectionModel{}

	itemsAllDose, errAllDose := s.repo.GetAllDose(ctx, owner)
	if errAllDose != nil {
		return nil, errAllDose
	}

	for _, item := range items {
		resultModel := InjectionModel{}
		resultModel.Injection = item
		resultModel.Injection_Dose = []entity.Injection_Dose{}

		for _, itemDose := range itemsAllDose {
			//logger.Infof("injection.Injection_Dose " + itemDose.Drug)
			if itemDose.Id_injection == item.ID {
				resultModel.Injection_Dose = append(resultModel.Injection_Dose, itemDose)
			}
		}
		result = append(result, resultModel)
	}
	return result, nil
}

func (s service) GetinjReort(ctx context.Context, owner string, sDate string, fDate string) (Points, error) {
	result := Points{}

	date1, err := time.Parse("2006-01-02", sDate)
	//fmt.Println(date1)
	if err != nil {
		return Points{}, err
	}

	date2, err := time.Parse("2006-01-02", fDate)
	if err != nil {
		return Points{}, err
	}
	// добавляем + 1 день
	date2 = date2.AddDate(0, 0, 1)

	ConcentrationDrugs, errConcentrationDrugs := s.repo.GetConcentrationDrugs(ctx, owner, date1, date2)

	if errConcentrationDrugs != nil {
		return Points{}, errConcentrationDrugs
	}

	for _, itemConcentrationDrugs := range ConcentrationDrugs {
		result.Drugs = append(result.Drugs, itemConcentrationDrugs.Drug)
	}

	Concentration, errConcentration := s.repo.GetConcentration(ctx, owner, "", date1, date2)

	if errConcentration != nil {
		return Points{}, errConcentration
	}

	var preDt int64 = 0
	var point entity.Point

	for i := 0; i < len(Concentration); i++ {
		//for i := 0; i < 10; i++ {
		//fmt.Println("start")
		//fmt.Println(i)

		if preDt != Concentration[i].Dt {
			//fmt.Println("preDt != Concentration[i].Dt ")
			//новая дата или первый раз

			if preDt != 0 {
				//если не первый раз сюда попали - последний блок всегда теряем ;)
				//fmt.Println("append(result.Points, point)")
				if len(point.PointValues) != len(result.Drugs)+1 {
					//fmt.Println("len(point.PointValues) != len(result.Drugs)")

					for _, Drug := range result.Drugs {
						var existsDrug = false
						for _, pV := range point.PointValues {
							if pV.Drug == Drug {
								existsDrug = true
							}
						}
						if !existsDrug {
							//fmt.Println(Drug)
							pValue := entity.PointValue{}

							pValue.Drug = Drug
							pValue.CCT = 0
							point.PointValues = append(point.PointValues, pValue)
						}
					}
				}
				result.Points = append(result.Points, point)
			}
			point = entity.Point{}
			point.Dt = Concentration[i].Dt

			point.PointValues = append(point.PointValues, entity.PointValue{})

			point.PointValues[0].Drug = "SUMM"
			point.PointValues[0].CCT = 0

			//fmt.Println(len(point.PointValues))
		}
		//fmt.Println(len(point.PointValues))
		//fmt.Println("pValue := entity.PointValue{}")

		pValue := entity.PointValue{}

		pValue.Drug = Concentration[i].Drug
		pValue.CCT = Concentration[i].CCT

		point.PointValues[0].CCT += pValue.CCT

		point.PointValues = append(point.PointValues, pValue)

		//result.Points = append(result.Points, point)

		preDt = Concentration[i].Dt
	}
	/*
		for _, itemConcentrationDrugs := range ConcentrationDrugs {

			if errConcentration != nil {
				return Points{}, errConcentration
			}
			for _, itemConcentration := range Concentration {

					point := entity.Point{}
					point.Dt = itemConcentration.Dt

					point.PointValues = append(point.PointValues, entity.PointValue{})

					point.PointValues[0].Drug = "SUMM"
					point.PointValues[0].CCT = 0

					pValue := entity.PointValue{}

					pValue.Drug = itemConcentration.Drug
					pValue.CCT = itemConcentration.CCT

					//point.PointValues[0].CCT += pValue.CCT

					point.PointValues = append(point.PointValues, pValue)

					result.Points = append(result.Points, point)

				//fmt.Println(itemConcentration.Drug)
			}

		}
	*/
	//sort.Sort(result.Points)
	//подпорка для фронта чтобы не падал когда в отчёте пусто
	if len(result.Drugs) == 0 {
		result.Drugs = []string{}
		result.Points = []entity.Point{}
	}
	return result, nil
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func (s service) Getinj(ctx context.Context, id string, owner string, save bool) (Points, error) {
	logger := s.logger.With(ctx, "id", id)
	logger.Infof("Getinj start")

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
		//logger.Infof("injection.Injection_Dose " + injection.Injection.Dt.String())
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

			//if item.Drug != "" {
			//logger.Infof("injection.Injection_Dose " + item.ID)
			//}

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
		const ZERO = 1e-6   //что считать нулем

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
					pValue.C = item.Dose * 0.7 * (drugs[item.Drug].Out / 100.0)
					//this.routt[idx]=inj.dose[idx]*0.7*(drug.out/100)*(drug.outt/100);
					pValue.CT = item.Dose * 0.7 * (drugs[item.Drug].Out / 100.0) * drugs[item.Drug].Out / 100.0

				} else {
					pValue.C = result.Points[count-1].PointValues[count_Injection_Dose].C
					pValue.CT = result.Points[count-1].PointValues[count_Injection_Dose].CT
				}

				pValue.C = pValue.C * math.Exp(-math.Ln2/float64(drugs[item.Drug].Halflife))
				if pValue.C < ZERO {
					pValue.C = 0
				} else {
					condition = true
				}

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

		//иньекция
		point := entity.Point{}
		//logger.Infof("injection.Injection_Dose " + injection.Injection.Dt.String())
		point.Dt = injection.Injection.Dt.Unix() * 1000 // приведение в формат старого кода

		//резерв под сумму
		point.PointValues = append(point.PointValues, entity.PointValue{})
		point.PointValues[0].Drug = "SUMM"
		point.PointValues[0].C = 0
		point.PointValues[0].CC = 0
		point.PointValues[0].CCT = 0
		point.PointValues[0].CT = 0

		//проходим по всем дозам

		for _, item := range injection.Injection_Dose {

			//if item.Drug != "" {
			//	logger.Infof("injection.Injection_Dose " + item.ID)
			//}

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
		const ZERO = 1e-6   //что считать нулем
		const ZERO2 = 1e-4  //что считать нулем
		injection.Injection.SkinSumm = 0
		injection.Injection.TotalV = 0

		for {

			//создаём объект для графика
			point := entity.Point{}
			//logger.Infof("injection.Injection_Dose " + injection.Injection.Dt.String())

			//иньекция работает сразу
			point.Dt = result.Points[count-1].Dt + cStep

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
					condition = true

					pValue.OutK = 0.95 * (drugs[item.Drug].Out / 100.0)
					pValue.OutKT = 0.95 * (drugs[item.Drug].Out / 100.0) * (drugs[item.Drug].Outt / 100.0)

					pValue.Dose = item.Dose * 1
					pValue.Volume = item.Volume * 1

					//logger.Infof("item.OutK=" + fmt.Sprintf("%f",item.OutK))
					//logger.Infof("item.OutKT=" + fmt.Sprintf("%f",item.OutKT))

					pValue.C = 0
					pValue.CT = 0
					//this.R[idx]=Math.pow((3*inj.volume[idx])/(4*Math.PI),1/3);
					pValue.R = math.Pow((3.0*pValue.Volume)/(4.0*math.Pi), 1.0/3.0)

					var skin = utils.SkinStep(item.Solvent + injection.Injection.What)

					//logger.Infof("SkinStep=" + fmt.Sprintf("%f", skin))

					injection.Injection.TotalV = pValue.Volume * 1
					injection.Injection.SkinSumm = skin * (pValue.Volume * 1)

				} else {
					pValue.C = result.Points[count-1].PointValues[count_Injection_Dose].C
					pValue.CT = result.Points[count-1].PointValues[count_Injection_Dose].CT
					pValue.R = result.Points[count-1].PointValues[count_Injection_Dose].R
					pValue.OutK = result.Points[count-1].PointValues[count_Injection_Dose].OutK
					pValue.OutKT = result.Points[count-1].PointValues[count_Injection_Dose].OutKT
					pValue.Dose = result.Points[count-1].PointValues[count_Injection_Dose].Dose
					pValue.Volume = result.Points[count-1].PointValues[count_Injection_Dose].Volume

					//var Cout=float64(0)
					//var Coutt=float64(0)

					//if (R>ZERO) [Cout,Coutt]=this.ballOut(i);
					//var r=pValue.R
					//logger.Infof("pValue.R=" + fmt.Sprintf("%f",r))
					if pValue.R < ZERO {
						pValue.Cout = 0.0
						pValue.Coutt = 0.0
					} else {
						//logger.Infof("injection.Injection.Skin=" + fmt.Sprintf("%f", injection.Injection.Skin))
						pValue.Ri = pValue.R - injection.Injection.Skin
						if pValue.Ri < 0.0 {
							pValue.Ri = 0.0
						}

						pValue.Depo = pValue.Dose
						//var depoi=depo*((4/3*Math.PI*Math.pow(ri,3))/this.volume[idx]); //считаем новый объем
						//(4/3*math.Pi*math.Pow(pValue.Ri,3))  объём шара
						pValue.Depoi = pValue.Depo * ((4.0 / 3.0 * math.Pi * math.Pow(pValue.Ri, 3.0)) / pValue.Volume) //считаем новый объем

						pValue.Dv = pValue.Depo - pValue.Depoi
						pValue.R = pValue.Ri
						if pValue.Depoi < ZERO {
							pValue.Depoi = 0
						}

						pValue.Dose = pValue.Depoi

						if pValue.Dv < ZERO {
							pValue.Cout = 0.0
							pValue.Coutt = 0.0
						} else {
							pValue.Cout = pValue.Dv * pValue.OutK
							pValue.Coutt = pValue.Dv * pValue.OutKT
						}

					}

					//var a=(this.CO[i]*Math.exp(-Math.LN2/this.halflife[i]))+Cout;
					pValue.C = (pValue.C * math.Exp(-math.Ln2/float64(drugs[item.Drug].Halflife))) + pValue.Cout
					if pValue.C < ZERO {
						pValue.C = 0.0
					} else {
						condition = true
					}

					//this.COT[i]=(this.COT[i]*Math.exp(-Math.LN2/this.halflife[i]))+Coutt;
					pValue.CT = (pValue.CT * math.Exp(-math.Ln2/float64(drugs[item.Drug].Halflife))) + pValue.Coutt

					pValue.CC = pValue.C
					pValue.CCT = pValue.CT

				}
				point.PointValues[0].C += pValue.C
				point.PointValues[0].CC += pValue.CC
				point.PointValues[0].CCT += pValue.CCT
				point.PointValues[0].CT += pValue.CT

				point.PointValues = append(point.PointValues, pValue)
			}

			if count == 1 {
				injection.Injection.Skin = injection.Injection.SkinSumm / injection.Injection.TotalV
				//logger.Infof("injection.Injection.Skin=" + fmt.Sprintf("%f", injection.Injection.Skin))
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
	}

	//чистим данные по прошлым расчётам
	//errDelete := s.repo.DeleteConcentration(ctx, owner, id)
	//if errDelete != nil {
	//	return Points{}, errDelete
	//}
	//сохроняем в БД если ранее не делали
	//fmt.Println(injection.Injection.ID)
	//fmt.Println(injection.Injection.Calc)
	if !injection.Injection.Calc && save {
		injection.Injection.Calc = true
		fmt.Println(1111111)
		errUpdateInjection := s.repo.UpdateInjection(ctx, injection.Injection)
		if errUpdateInjection != nil {
			return Points{}, errUpdateInjection
		}
		//errGrp, _ := errgroup.WithContext(context.Background())

		var arryaConcentration []entity.Concentration
		for _, Drug := range result.Drugs {
			logger.Infof("save injection Drugs = " + Drug + ", result.Points=" + fmt.Sprintf(":%v", len(result.Points)))
			var count = 15
			arryaConcentration = []entity.Concentration{}

			for _, itemPoint := range result.Points {
				//fmt.Println(itemPoint.Dt)
				for _, itemPoints := range itemPoint.PointValues {

					if itemPoints.Drug == Drug {
						//каждый 15й элемент
						if count%15 == 0 {
							arryaConcentration = append(arryaConcentration, entity.Concentration{
								Owner:        owner,
								Id_injection: id,
								Drug:         itemPoints.Drug,
								Dt:           itemPoint.Dt,
								C:            itemPoints.C,
								CC:           itemPoints.CC,
								CCT:          itemPoints.CCT,
								CT:           itemPoints.CT,
							})
						}
						count++
					}

				}
			}
			logger.Infof("arryaConcentration lengt = " + fmt.Sprintf("%v", len(arryaConcentration)))
			s.repo.SaveConcentration2(ctx, arryaConcentration)
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
	//удаляем прошлые расчёты
	errDelete := s.repo.DeleteConcentration(ctx, owner, id)
	if errDelete != nil {
		return InjectionModel{}, err
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

	//logger.Infof("Injection_Dose=" + fmt.Sprintf("%f", len(injection.Injection_Dose)))
	if len(injection.Injection_Dose) == 0 {
		if err = s.repo.Delete(ctx, id); err != nil {
			return InjectionModel{}, err
		}
		//чистим данные по прошлым расчётам
		errDelete := s.repo.DeleteConcentration(ctx, owner, id)
		if errDelete != nil {
			return InjectionModel{}, errDelete
		}
		return InjectionModel{}, nil
	}

	return injection, nil
}

func (s service) GetinjOne(ctx context.Context, id string) (InjectionModel, error) {

	injection, err := s.GetOne(ctx, id)

	if err != nil {
		return InjectionModel{}, err
	}

	return injection, nil
}

func (s service) GetForBloodVolume(ctx context.Context, owner string) ([]entity.Antro, error) {
	antro, err := s.repo.GetForBloodVolume(ctx, owner)
	if err != nil {
		return []entity.Antro{}, err
	}
	return antro, nil
}
