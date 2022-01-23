package profile

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Post("/api/createprofile", res.create)
	r.Post("/api/restorepassword", res.restorepassword)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/api/updateprofile", res.create)
	r.Get("/api/profile", res.get)
	r.Get("/api/historylogin", res.getHistoryLogin)

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
		//fmt.Println(profile.TypeSports)
		if err != nil {
			return err
		}

		//подпорка масива из Бд на С
		strs := strings.Split(strings.ReplaceAll(strings.ReplaceAll(profile.TypeSports, "{", ""), "}", ""), ",")

		ary := make([]int, len(strs))
		for i := range ary {
			ary[i], _ = strconv.Atoi(strs[i])
		}

		if profile.TypeSports == "{}" {
			ary = []int{}
		}

		profileRest := entity.UsersRest{ID: profile.ID,
			Login:          profile.Login,
			Email:          profile.Email,
			DateRegistered: profile.DateRegistered,
			DateLastlogin:  profile.DateLastlogin,
			Sex:            profile.Sex,
			Birthday:       profile.Birthday,
			TypeSports:     ary,
		}

		//TypeSports     []string  `json:"typesports"`
		return c.Write(profileRest)
	} else {
		return err
	}

}

func (r resource) getHistoryLogin(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		HistoryLogin, err := r.service.GetHistoryLogin(c.Request.Context(), claims["id"].(string))
		if err != nil {
			return err
		}

		return c.Write(HistoryLogin)
	} else {
		return err
	}

}

func GetValue(service Service, logger log.Logger, c *routing.Context) string {
	res := resource{service, logger}
	a := res.get(c)
	return a.Error()
}

func (r resource) create(c *routing.Context) error {
	var input CreateUser
	//fmt.Println("111")
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	//fmt.Println(22222)
	user1, err1 := r.service.GetByLogin(c.Request.Context(), input.Login)

	if err1 == nil && user1.Login == input.Login {
		return errors.BadRequest("login exists")
	}

	user1, err1 = r.service.GetByEmail(c.Request.Context(), input.Email)
	//fmt.Println(input.Email)
	//fmt.Println(user1.Email)

	if err1 == nil && user1.Email == input.Email {
		return errors.BadRequest("Email exists")
	}

	user, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}
	var arr []string

	json.Unmarshal([]byte(user.TypeSports), &arr)
	fmt.Println(len(arr))

	return c.WriteWithStatus(user, http.StatusCreated)
	/*return c.WriteWithStatçus(entity.UsersRest{
	ID:             user.ID,
	Login:          user.Login,
	Email:          user.Email,
	DateRegistered: user.DateLastlogin,
	Sex:            user.Sex,
	Birthday:       user.Birthday,
	TypeSports:		"{"+ strings.Join(req.TypeSports,",") + "}", http.StatusCreated)*/
}

func (r resource) restorepassword(c *routing.Context) error {
	var input CreateUser
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	fmt.Println(input.Email)

	user, err := r.service.GetByEmail(c.Request.Context(), input.Email)

	if err != nil {
		r.logger.Error("restorepassword error:" + err.Error())
		return c.WriteWithStatus("", http.StatusAccepted)
	}

	fmt.Println(user.Email)
	fmt.Println(user.Passwd)

	utils.SendMail("smtp.chartdrug.com:587", (&mail.Address{"info@chartdrug.com", "info@chartdrug.com"}).String(),
		"Restore password", "message body", []string{(&mail.Address{"", "k.dereshev@ya.ru"}).String()})

	return c.WriteWithStatus("", http.StatusBadGateway)
}
