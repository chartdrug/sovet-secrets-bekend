package payment

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

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	//сюда пробрасывается уведомление о результату
	r.Post("/api/crypto-result-pay", res.cryptoPaymentResult)

	r.Use(authHandler)
	// создать платеж
	r.Get("/api/crypto-payment-create", res.cryptoPaymentCreate)
	//получить все платежи и обновить по ним статус
	r.Get("/api/crypto-payment-all", res.cryptoPaymentAll)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) cryptoPaymentResult(c *routing.Context) error {

	var input entity.CryptocloudPostback
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	err := r.service.cryptoPaymentResult(c.Request.Context(), c.Request.Header, input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus("", http.StatusOK)
}

func (r resource) cryptoPaymentAll(c *routing.Context) error {
	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		paym, err2 := r.service.cryptoPaymentAll(c.Request.Context(), claims["id"].(string))
		if err2 != nil {
			return err2
		}

		return c.WriteWithStatus(paym, http.StatusOK)

	} else {
		return err
	}
}

func (r resource) cryptoPaymentCreate(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		//paym, err := r.service.cryptoPaymentCreate(c.Request.Context(), claims["id"].(string))
		paym, err := r.service.cryptoPaymentCreate(c.Request.Context(), claims)
		if err != nil {
			return err
		}

		return c.WriteWithStatus(paym, http.StatusCreated)

	} else {
		return err
	}

}
