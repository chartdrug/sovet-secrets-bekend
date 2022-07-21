package injections

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"math"
	"net/http"
	"strings"
	"time"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	//расчёт раз в минуту для отчёта
	r.Get("/api/injectionsCall", res.asyncCall)
	r.Get("/api/injectionsCall/<id>", res.asyncCallId)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/api/injections/<sDate>/<fDate>/", res.getAllDose)
	//r.Get("/api/injections/", res.getAllDose)
	//r.Get("/api/injections/report", res.getReort)
	r.Get("/api/injections/report/<sDate>/<fDate>/", res.getReort)
	r.Get("/api/injections/course/report/<id>", res.getReort2)
	r.Get("/api/injections/inj/<id>", res.getinj)
	//r.Delete("/api/antros/<id>", res.delete)
	r.Post("/api/injection", res.create)
	r.Post("/api/injection2", res.create2)
	r.Delete("/api/injection/<id>", res.delete)
	r.Delete("/api/injections", res.deleteArray)
	r.Delete("/api/injection/<id>/dose/<id_dose>", res.deleteDose)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getAllDose(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		injection, err := r.service.GetAllDose(c.Request.Context(), claims["id"].(string), c.Param("sDate"), c.Param("fDate"))
		if err != nil {
			return err
		}

		return c.Write(injection)
	} else {
		return err
	}

}
func (r resource) getinj(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		injection, err := r.service.Getinj(c.Request.Context(), c.Param("id"), claims["id"].(string), false)
		if err != nil {
			return err
		}

		//оставляем только каждые 15 минут
		for i := len(injection.Points) - 1; i >= 0; i-- {
			if i%15 != 0 {
				injection.Points = append(injection.Points[:i], injection.Points[i+1:]...)
			}
		}
		//получаем обём крови
		antro, errAntro := r.service.GetForBloodVolume(c.Request.Context(), claims["id"].(string))
		if errAntro != nil {
			return errAntro
		}

		var BloodVolume = utils.GetBloodVolume(claims["sex"].(string), antro)

		//кол-во доступных вариантов крови
		var b = 0
		var lb = len(BloodVolume)

		if lb == 0 {
			return errors.NotFound("Вам нужно заполнить информацию по aнтропометрии")
		}

		//удаляем всё что меньше 0.2 и делаем округелние
		for i := len(injection.Points) - 1; i >= 0; i-- {
			//fmt.Println("test2")
			application := injection.Points[i]
			// Condition to decide if current element has to be deleted:
			if application.PointValues[0].CCT < 0.2 {
				//fmt.Println("server %v\n", application.PointValues[0].CCT)
				injection.Points = append(injection.Points[:i],
					injection.Points[i+1:]...)
			} else {
				// если расчёт крови боль даты инькции и есть ещё древние то смещаемся пока не найдём
				if BloodVolume[b].Dt > application.Dt && b != lb-1 {
					fmt.Println("----")
					for y := b; y < lb; y++ {
						if BloodVolume[y].Dt <= application.Dt {
							b = y
							break
						}
					}
				}
				for y := 0; y < len(application.PointValues); y++ {
					//application.PointValues[y].CCT = math.Round(application.PointValues[y].CCT*1000) / 1000
					application.PointValues[y].CCT = math.Round((((application.PointValues[y].CCT/288431)*1000000000)/BloodVolume[b].V)*1000) / 1000
				}

			}
		}

		return c.Write(injection)
	} else {
		return err
	}

}

func (r resource) getinjRep(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		injection, err := r.service.Getinj(c.Request.Context(), c.Param("id"), claims["id"].(string), false)
		if err != nil {
			return err
		}

		//оставляем только каждые 15 минут
		for i := len(injection.Points) - 1; i >= 0; i-- {
			if i%15 != 0 {
				injection.Points = append(injection.Points[:i], injection.Points[i+1:]...)
			}
		}
		//получаем обём крови
		antro, errAntro := r.service.GetForBloodVolume(c.Request.Context(), claims["id"].(string))
		if errAntro != nil {
			return errAntro
		}

		var BloodVolume = utils.GetBloodVolume(claims["sex"].(string), antro)

		//кол-во доступных вариантов крови
		var b = 0
		var lb = len(BloodVolume)

		if lb == 0 {
			return errors.NotFound("Вам нужно заполнить информацию по aнтропометрии")
		}

		//удаляем всё что меньше 0.2 и делаем округелние
		for i := len(injection.Points) - 1; i >= 0; i-- {
			//fmt.Println("test2")
			application := injection.Points[i]
			// Condition to decide if current element has to be deleted:
			if application.PointValues[0].CCT < 0.2 {
				//fmt.Println("server %v\n", application.PointValues[0].CCT)
				injection.Points = append(injection.Points[:i],
					injection.Points[i+1:]...)
			} else {
				// если расчёт крови боль даты инькции и есть ещё древние то смещаемся пока не найдём
				if BloodVolume[b].Dt > application.Dt && b != lb-1 {
					fmt.Println("----")
					for y := b; y < lb; y++ {
						if BloodVolume[y].Dt <= application.Dt {
							b = y
							break
						}
					}
				}
				for y := 0; y < len(application.PointValues); y++ {
					//application.PointValues[y].CCT = math.Round(application.PointValues[y].CCT*1000) / 1000
					application.PointValues[y].CCT = math.Round((((application.PointValues[y].CCT/288431)*1000000000)/BloodVolume[b].V)*1000) / 1000
				}

			}
		}

		return c.Write(injection)
	} else {
		return err
	}

}

func (r resource) getReort2(c *routing.Context) error {
	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		injection, err := r.service.GetinjReort2(c.Request.Context(), claims["id"].(string), c.Param("id"))
		if err != nil {
			return err
		}

		//получаем обём крови
		antro, errAntro := r.service.GetForBloodVolume(c.Request.Context(), claims["id"].(string))
		if errAntro != nil {
			return errAntro
		}

		var BloodVolume = utils.GetBloodVolume(claims["sex"].(string), antro)

		//кол-во доступных вариантов крови
		var b = 0
		var lb = len(BloodVolume)
		var countDaysAndro = 0.0

		if lb == 0 {
			return errors.NotFound("Вам нужно заполнить информацию по aнтропометрии")
		}

		//Расчет даты начала ПКТ (посткурсовая терапия)
		//Надо найти дату > последнего применения в которой общая
		//концентрация снижается до 10 нмоль. (если юзер использовал препараты так что не
		//поднимался выше 10 нмоль за весь курс (что весьма маловероятно, ибо теряется всякий
		//смысл использования анаболиков), то такая дата найдена быть не может и в этом случае
		//надо определить дату начала ПКТ как дата последнего применения + 7 дней.

		for _, item := range injection.Injections {
			if item.Injection.Dt.IsZero() ||
				item.Injection.Dt.After(injection.Pkt) {
				injection.Pkt = item.Injection.Dt
			}
		}

		if !injection.Pkt.IsZero() {
			injection.Pkt = injection.Pkt.AddDate(0, 0, 7)
		}

		//Анаболический индекс
		/*
			Это отношение площади фигуры общего графика выше 30нмоль к площади
			прямоугольника ниже 30 нмоль. (отношение красной площади к черной).
			Считать видимо так: каждую единицу времени (каждый шаг цикла расчета
			концентрации) делить общую концентрацию на 30 (если концентрация <30, то
			присвоить результату деления значение 0). Сумма этих частных за курс и будет
			анаболическим индексом курса. Для красоты умножить на 100 и округлить до целого
			(либо откинуть дробную часть).
		*/
		injection.AnabolicIndex = 0.0
		//Андрогенный индекс
		/*
			Это отношение средней концентрации (считается также по общей) к 30.
			Среднее значение вычисляется из ранее полученных частных. Для красоты умножить
			на 100 и округлить до целого (либо откинуть дробную часть).
		*/
		injection.AndrogenicIndex = 0.0
		cAndrogenicIndex := 1.0

		first := true
		lastpkt := true

		if (len(injection.Points) - 1) > 0 {
			injection.Control = time.Unix(injection.Points[len(injection.Points)-1].Dt/1000, 0)
		}

		for i := len(injection.Points) - 1; i >= 0; i-- {

			//fmt.Println("test2")
			application := injection.Points[i]
			//PKT
			if math.Round((((application.PointValues[0].CCT/288431)*1000000000)/BloodVolume[b].V)*1000)/1000 <= 10 && lastpkt {
				injection.Pkt = time.Unix(application.Dt/1000, 0)
			} else {
				lastpkt = false
			}
			//injection.AnabolicIndex
			if ((application.PointValues[0].CCT/288431)*1000000000)/BloodVolume[b].V >= 30.0 {
				injection.AnabolicIndex = injection.AnabolicIndex + (((application.PointValues[0].CCT / 288431) * 1000000000) / BloodVolume[b].V / 30.0)
				//injection.AndrogenicIndex
				//считаем от даты ПКТ
				if !lastpkt {
					injection.AndrogenicIndex = injection.AndrogenicIndex + (((application.PointValues[0].CCT / 288431) * 1000000000) / BloodVolume[b].V / 30.0)
					cAndrogenicIndex = cAndrogenicIndex + 1.0
				}
			}

			if math.Round((((application.PointValues[0].CCT/288431)*1000000000)/BloodVolume[b].V)*1000)/1000 < 0.2 && first {
				injection.Points = append(injection.Points[:i],
					injection.Points[i+1:]...)
			} else {
				first = false
				for y := 0; y < len(application.PointValues); y++ {
					application.PointValues[y].CCT = math.Round((((application.PointValues[y].CCT/288431)*1000000000)/BloodVolume[b].V)*1000) / 1000
				}
			}
		}

		if !injection.Control.IsZero() {
			injection.Control = injection.Control.AddDate(0, 0, 7)
		}

		if len(injection.Points) > 0 {
			countDaysAndro = injection.Pkt.Sub(time.Unix(injection.Points[0].Dt/1000, 0)).Hours() / 24
		}
		//андрогенный множить количество дней курса( от первого дня до дня ПКТ)
		injection.AnabolicIndex = math.Round(injection.AnabolicIndex*100) / 100
		injection.AndrogenicIndex = math.Round((injection.AndrogenicIndex/cAndrogenicIndex)*countDaysAndro*100) / 100

		return c.Write(injection)
	} else {
		return err
	}
}

func (r resource) getReort(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		injection, err := r.service.GetinjReort(c.Request.Context(), claims["id"].(string), c.Param("sDate"), c.Param("fDate"))
		if err != nil {
			return err
		}

		first := true

		//получаем обём крови
		antro, errAntro := r.service.GetForBloodVolume(c.Request.Context(), claims["id"].(string))
		if errAntro != nil {
			return errAntro
		}

		var BloodVolume = utils.GetBloodVolume(claims["sex"].(string), antro)

		//кол-во доступных вариантов крови
		var b = 0
		var lb = len(BloodVolume)

		if lb == 0 {
			return errors.NotFound("Вам нужно заполнить информацию по aнтропометрии")
		}

		for i := len(injection.Points) - 1; i >= 0; i-- {

			//fmt.Println("test2")
			application := injection.Points[i]
			//fmt.Println(application.PointValues[0].CCT )
			// Condition to decide if current element has to be deleted:
			if application.PointValues[0].CCT < 0.2 && first {
				//fmt.Println("server %v\n", application.PointValues[0].CCT)
				injection.Points = append(injection.Points[:i],
					injection.Points[i+1:]...)
			} else {
				first = false
				for y := 0; y < len(application.PointValues); y++ {
					//application.PointValues[y].CCT = math.Round(application.PointValues[y].CCT*1000) / 1000
					application.PointValues[y].CCT = math.Round((((application.PointValues[y].CCT/288431)*1000000000)/BloodVolume[b].V)*1000) / 1000
				}
			}
		}

		return c.Write(injection)
	} else {
		return err
	}

}

func (r resource) create(c *routing.Context) error {
	var input CreateInjectionsRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		//проверяем что объм не больше 10
		var summ = 0.0
		for _, item := range input.Injection_Dose {
			summ += item.Volume
		}

		Course, errCourse := r.service.GetCourse(c.Request.Context(), input.Injection.Course)
		if err != nil {
			return errCourse
		}

		//округляем дату
		d := 24 * time.Hour
		if input.Injection.Dt.Before(Course.Course_start) || input.Injection.Dt.Truncate(d).After(Course.Course_end) {
			return errors.BadRequest("Применение фармакологии должно быть в период курса")
		}

		if summ > 10 {
			return errors.BadRequest("All values must be less than or equal to 10")
		}
		//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		//fmt.Println(input.Injection.ID)
		if len(input.Injection.ID) > 1 {
			fmt.Println("service.Delete")
			_, errDell := r.service.Delete(c.Request.Context(), input.Injection.ID, claims["id"].(string))
			if err != nil {
				return errDell
			}
		}

		err := r.service.Create(c.Request.Context(), input, claims["id"].(string))
		if err != nil {
			return err
		}
		//делаем расчёт для сохранения в БД
		/* расчёт идёт по крону
		go r.service.Getinj(c.Request.Context(), injection.Injection.ID, claims["id"].(string), true)

		fmt.Println("Sleep 1 second")
		time.Sleep(1 * time.Second)

		for i := 0; i <= 3; i++ {
			fmt.Println("Start check status")

			injOne, errinjOne := r.service.GetinjOne(c.Request.Context(), injection.Injection.ID)
			if errinjOne != nil {
				return errinjOne
			}
			if injOne.Injection.Calc {
				break
			}
			fmt.Println("retray GetinjOne")
			fmt.Println("Sleep 1 second")
		}

		*/

		return c.WriteWithStatus("injection", http.StatusCreated)

	} else {
		return err
	}

}

func (r resource) create2(c *routing.Context) error {
	var input entity.InjectionModelv2
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		//проверяем что объм не больше 10
		var summ = 0.0
		for _, item := range input.Injection_Dose {
			summ += item.Volume
		}

		if summ > 10 {
			return errors.BadRequest("All values must be less than or equal to 10")
		}
		//fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		//fmt.Println(input.Injection.ID)

		/*
			//седлать удаление потом
			if len(input.Injection.ID) > 1 {
				fmt.Println("service.Delete")
				_, errDell := r.service.Delete(c.Request.Context(), input.Injection.ID, claims["id"].(string))
				if err != nil {
					return errDell
				}
			}*/

		//var arrayUid []string

		//проверка что даты входя в период курса

		Course, errCourse := r.service.GetCourse(c.Request.Context(), input.Injection.Course)
		if err != nil {
			return errCourse
		}

		for _, d := range input.Injection.Date {
			dt, _ := time.Parse("2006-01-02", d)
			if dt.Before(Course.Course_start) || dt.After(Course.Course_end) {
				return errors.BadRequest("Применение фармакологии должно быть в период курса")
			}
		}

		for _, d := range input.Injection.Date {
			for _, t := range input.Injection.Times {
				inj := CreateInjectionsRequest{}

				inj.Injection_Dose = input.Injection_Dose

				inj.Injection.Dt, _ = time.Parse("2006-01-02", d)
				fmt.Println(inj.Injection.Dt)
				inj.Injection.Dt = inj.Injection.Dt.Add(time.Minute * time.Duration(t/(900000/15)))
				fmt.Println(inj.Injection.Dt)

				inj.Injection.What = input.Injection.What
				inj.Injection.Course = input.Injection.Course

				err := r.service.Create(c.Request.Context(), inj, claims["id"].(string))
				if err != nil {
					return err
				}
				//arrayUid = append(arrayUid, injection.Injection.ID)
			}
		}

		//делаем расчёт для сохранения в БД
		//сделать позже по массиву
		/* расчёт работает по крону
		go r.service.GetinjArray(c.Request.Context(), arrayUid, claims["id"].(string), true)

		fmt.Println("Sleep 1 second")
		time.Sleep(1 * time.Second)

		for i := 0; i <= 3; i++ {
			fmt.Println("Start check status")

			for _, u := range arrayUid {
				injOne, errinjOne := r.service.GetinjOne(c.Request.Context(), u)
				if errinjOne != nil {
					return errinjOne
				}
				if injOne.Injection.Calc {
					break
				}
			}
			fmt.Println("retray GetinjOne")
			fmt.Println("Sleep 1 second")
		}
		*/

		return c.WriteWithStatus("arrayUid", http.StatusCreated)

	} else {
		return err
	}

}

func (r resource) delete(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		injection, err := r.service.Delete(c.Request.Context(), c.Param("id"), claims["id"].(string))
		if err != nil {
			return err
		}

		return c.Write(injection)

	} else {
		return err
	}

}

func (r resource) deleteArray(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	var input []string

	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		for _, item := range input {
			_, err := r.service.Delete(c.Request.Context(), item, claims["id"].(string))
			if err != nil {
				return err
			}
		}

		return c.WriteWithStatus(input, http.StatusAccepted)

	} else {
		return err
	}

}

func (r resource) deleteDose(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		injection, err := r.service.DeleteDose(c.Request.Context(), c.Param("id"), c.Param("id_dose"), claims["id"].(string))
		if err != nil {
			return err
		}

		return c.Write(injection)

	} else {
		return err
	}

}
func (r resource) asyncCall(c *routing.Context) error {

	err := r.service.AsyncCall(c.Request.Context())

	if err != nil {
		return err
	}
	return c.Write("Done asyncCall")

}
func (r resource) asyncCallId(c *routing.Context) error {

	err := r.service.AsyncCallID(c.Request.Context(), c.Param("id"))

	if err != nil {
		return err
	}
	return c.Write("Done AsyncCallID")

}
