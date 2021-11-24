package spr

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/api/spr/sports/<language>", res.getSports)
	r.Use(authHandler)

	// the following endpoints require a valid JWT
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getSports(c *routing.Context) error {

	spr, err := r.service.GetSports(c.Request.Context(), c.Param("language"))
	if err != nil {
		return err
	}

	return c.Write(spr)

}
