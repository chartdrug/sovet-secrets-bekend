package tbot

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Get("/api/chatsQuestion", res.getChatsQuestion)
	r.Post("/api/chatsQuestion", res.CreateChatsQuestion)
	r.Put("/api/chatsQuestion", res.UpdateChatsQuestion)
	r.Delete("/api/chatsQuestion/<id>", res.DeleteChatsQuestion)

	r.Get("/api/chatsAnswer", res.getChatsAnswer)
	r.Post("/api/chatsAnswer", res.CreateChatsAnswer)
	r.Put("/api/chatsAnswer", res.UpdateChatsAnswer)
	r.Delete("/api/chatsAnswer/<id>", res.DeleteChatsAnswer)

	r.Get("/api/UsersQuestion", res.getUsersQuestion)

	r.Use(authHandler)

	// the following endpoints require a valid JWT

}

type resource struct {
	service Service
	logger  log.Logger
}

// блок вопросов
func (r resource) getChatsQuestion(c *routing.Context) error {

	chatsQuestion, err := r.service.GetChatsQuestion(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) CreateChatsQuestion(c *routing.Context) error {

	var input entity.ChatsQuestion
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	err := r.service.CreateChatsQuestion(c.Request.Context(), input)
	if err != nil {
		return err
	}

	chatsQuestion, err := r.service.GetChatsQuestion(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) UpdateChatsQuestion(c *routing.Context) error {

	var input entity.ChatsQuestion
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	err := r.service.UpdateChatsQuestion(c.Request.Context(), input)
	if err != nil {
		return err
	}

	chatsQuestion, err := r.service.GetChatsQuestion(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) DeleteChatsQuestion(c *routing.Context) error {
	err := r.service.DeleteChatsQuestion(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write("{}")
}

// блок ответов
func (r resource) getChatsAnswer(c *routing.Context) error {

	chatsQuestion, err := r.service.GetChatsAnswer(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) CreateChatsAnswer(c *routing.Context) error {

	var input entity.ChatsAnswer
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	err := r.service.CreateChatsAnswer(c.Request.Context(), input)
	if err != nil {
		return err
	}

	chatsQuestion, err := r.service.GetChatsAnswer(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) UpdateChatsAnswer(c *routing.Context) error {

	var input entity.ChatsAnswer
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	err := r.service.UpdateChatsAnswer(c.Request.Context(), input)
	if err != nil {
		return err
	}

	chatsQuestion, err := r.service.GetChatsAnswer(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(chatsQuestion)
}

func (r resource) DeleteChatsAnswer(c *routing.Context) error {
	err := r.service.DeleteChatsAnswer(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write("{}")
}

func (r resource) getUsersQuestion(c *routing.Context) error {

	rq, err := r.service.GetUsersQuestion(c.Request.Context())
	if err != nil {
		return err
	}
	return c.Write(rq)
}
