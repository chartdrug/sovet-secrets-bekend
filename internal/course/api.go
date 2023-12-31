package course

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"net/http"
	"strings"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Get("/api/course", res.get)
	r.Delete("/api/course/<id>", res.delete)
	r.Post("/api/course", res.create)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		profile, err := r.service.Get(c.Request.Context(), claims["id"].(string))
		if err != nil {
			return err
		}

		return c.Write(profile)
	} else {
		return err
	}

}

func (r resource) delete(c *routing.Context) error {
	course, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(course)
}

func (r resource) create(c *routing.Context) error {
	var input entity.Course
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

		course, err := r.service.Create(c.Request.Context(), input, claims["id"].(string))
		if err != nil {
			return err
		}

		return c.WriteWithStatus(course, http.StatusCreated)

	} else {
		return err
	}

}
