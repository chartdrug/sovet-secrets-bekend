package injections

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"math"
	"net/http"
	"strings"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/api/injections/<sDate>/<fDate>/", res.getAllDose)
	//r.Get("/api/injections/", res.getAllDose)
	//r.Get("/api/injections/report", res.getReort)
	r.Get("/api/injections/report/<sDate>/<fDate>/", res.getReort)
	r.Get("/api/injections/inj/<id>", res.getinj)
	//r.Delete("/api/antros/<id>", res.delete)
	r.Post("/api/injection", res.create)
	r.Delete("/api/injection/<id>", res.delete)
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
		injection, err := r.service.Getinj(c.Request.Context(), c.Param("id"), claims["id"].(string))
		if err != nil {
			return err
		}

		//оставляем только каждые 15 минут
		for i := len(injection.Points) - 1; i >= 0; i-- {
			if i%15 != 0 {
				injection.Points = append(injection.Points[:i], injection.Points[i+1:]...)
			}
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
				for y := 0; y < len(application.PointValues); y++ {
					application.PointValues[y].CCT = math.Round(application.PointValues[y].CCT*1000) / 1000
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
					application.PointValues[y].CCT = math.Round(application.PointValues[y].CCT*1000) / 1000
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

		injection, err := r.service.Create(c.Request.Context(), input, claims["id"].(string))
		if err != nil {
			return err
		}

		return c.WriteWithStatus(injection, http.StatusCreated)

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
