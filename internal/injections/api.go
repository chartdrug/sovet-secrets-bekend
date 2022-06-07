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
	r.Get("/api/injectionsCall/<id>", res.asyncCall)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/api/injections/<sDate>/<fDate>/", res.getAllDose)
	//r.Get("/api/injections/", res.getAllDose)
	//r.Get("/api/injections/report", res.getReort)
	r.Get("/api/injections/report/<sDate>/<fDate>/", res.getReort)
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

		injection, err := r.service.Create(c.Request.Context(), input, claims["id"].(string))
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

		return c.WriteWithStatus(injection, http.StatusCreated)

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

		var arrayUid []string

		for _, d := range input.Injection.Date {
			for _, t := range input.Injection.Times {
				inj := CreateInjectionsRequest{}

				inj.Injection_Dose = input.Injection_Dose

				inj.Injection.Dt, _ = time.Parse("2006-01-02", d)
				fmt.Println(inj.Injection.Dt)
				inj.Injection.Dt = inj.Injection.Dt.Add(time.Minute * time.Duration(t/(900000/15)))
				fmt.Println(inj.Injection.Dt)

				inj.Injection.What = input.Injection.What

				injection, err := r.service.Create(c.Request.Context(), inj, claims["id"].(string))
				if err != nil {
					return err
				}
				arrayUid = append(arrayUid, injection.Injection.ID)
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

		return c.WriteWithStatus(arrayUid, http.StatusCreated)

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
