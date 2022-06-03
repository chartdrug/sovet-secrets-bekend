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
	r.Post("/api/updateprofile", res.update)
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
			Admin:          profile.Admin,
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

	//err = utils.SendMailGmail(user.Email, "Успешная регистрация на ChartDrug", user.Passwd, user.Login)
	body := "<div style=\"background-color:#f4f4f4;margin:0;padding:0\">\n    <div style=\"Margin:0;background-color:#f4f4f4;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;table-layout:fixed;text-align:center\">\n        <table class=\"m_7451640406457273268container\" cellpadding=\"0\" align=\"center\" style=\"border-collapse:collapse;margin:0 auto;max-width:600px;width:100%\">\n            <tbody>\n            <tr>\n                <td style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:0;line-height:18px;padding-bottom:26px;text-align:right\">\n                    <div class=\"m_7451640406457273268accessibilitylink\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;text-align:right;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;padding:25px 0 0 20px;text-align:right\">\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;background-color:#fff;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:600px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"30\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;border-bottom:#f4f4f4 1px solid;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:left\">\n                                    <h1 style=\"Margin:0 0 8px 0;color:#121212;font-family:helvetica,arial;font-size:23px;font-weight:700;line-height:1.3;text-align:left\">Спасибо за <span class=\"il\">регистрацию</span> на ChartDrug </h1>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Ваша <span class=\"il\">регистрация</span> успешно завершена, данные для доступа в личный кабинет: </p>\n                                    <p style=\"Margin:30px 0 0;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Логин: <b>{login}</b>\n                                    </p>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Ваш пароль: <b>{password}</b>\n                                    </p>\n                                    <div style=\"Margin:30px 0 30px;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:center\">\n                                        <a href=\"https://chartdrug.com\" style=\"background-color:#2b87db;border:#2b87db 11px solid;border-radius:25px;color:#fff!important;display:inline-block;font-family:helvetica;font-size:14px;font-weight:700;padding:0 15px;text-decoration:none\" target=\"_blank\">\n                                            <span style=\"color:#fff\">Перейти в личный кабинет</span>\n                                        </a>\n                                    </div>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Если вы не регистрировались на <a href=\"https://chartdrug.com\" style=\"color:#2b87db!important\" target=\"_blank\">\n                                        <span style=\"color:#2b87db\">ChartDrug</span>\n                                    </a>, просто удалите это письмо. </p>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;padding-top:10px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:580px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#999999;font-family:helvetica,arial;font-size:11px;line-height:14px;text-align:center\">\n                                    <b>© 2022 Chart Drug</b>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            </tbody>\n        </table>\n    </div>\n</div>"

	body = strings.Replace(body, "{login}", user.Login, 1)
	body = strings.Replace(body, "{password}", user.Passwd, 1)
	err = utils.SendMailGmail(user.Email, "Успешная регистрация на ChartDrug", body)

	if err != nil {
		fmt.Println(err.Error())
	}
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

func (r resource) update(c *routing.Context) error {

	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//profile, err := r.service.Get(c.Request.Context(), claims["id"].(string))
		var input CreateUser
		if err := c.Read(&input); err != nil {
			r.logger.With(c.Request.Context()).Info(err)
			return errors.BadRequest("")
		}

		user1, err1 := r.service.GetByLogin(c.Request.Context(), input.Login)

		//  смена логина и проверка что он не занят
		if err1 == nil && user1.Login == input.Login && user1.ID != claims["id"].(string) {
			return errors.BadRequest("login exists")
		}

		user1, err1 = r.service.GetByEmail(c.Request.Context(), input.Email)

		// смена емейл и проверка что он не занят
		if err1 == nil && user1.Email == input.Email && user1.ID != claims["id"].(string) {
			return errors.BadRequest("Email exists")
		}

		//user :=entity.Users{}
		user, err := r.service.Update(c.Request.Context(), claims["id"].(string), input)
		if err != nil {
			return err
		}

		//подпорка масива из Бд на С
		strs := strings.Split(strings.ReplaceAll(strings.ReplaceAll(user.TypeSports, "{", ""), "}", ""), ",")

		ary := make([]int, len(strs))
		for i := range ary {
			ary[i], _ = strconv.Atoi(strs[i])
		}

		if user.TypeSports == "{}" {
			ary = []int{}
		}

		profileRest := entity.UsersRest{ID: user.ID,
			Login:          user.Login,
			Email:          user.Email,
			DateRegistered: user.DateRegistered,
			DateLastlogin:  user.DateLastlogin,
			Sex:            user.Sex,
			Birthday:       user.Birthday,
			TypeSports:     ary,
		}

		return c.WriteWithStatus(profileRest, http.StatusAccepted)
	} else {
		return err
	}
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

	//fmt.Println(user.Email)
	//fmt.Println(user.Passwd)
	/*
		utils.SendMail("smtp.gmail.com:587", (&mail.Address{"chartdrug@gmail.com", "chartdrug@gmail.com"}).String(),
			"Restore password", "message body", []string{(&mail.Address{"", "k.dereshev@ya.ru"}).String()})*/
	//err = utils.SendMailGmail(user.Email, "ChartDrug restore password", "password "+user.Passwd+"for login "+user.Login)
	body := "<div style=\"background-color:#f4f4f4;margin:0;padding:0\">\n    <div style=\"Margin:0;background-color:#f4f4f4;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;table-layout:fixed;text-align:center\">\n        <table class=\"m_7451640406457273268container\" cellpadding=\"0\" align=\"center\" style=\"border-collapse:collapse;margin:0 auto;max-width:600px;width:100%\">\n            <tbody>\n            <tr>\n                <td style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:0;line-height:18px;padding-bottom:26px;text-align:right\">\n                    <div class=\"m_7451640406457273268accessibilitylink\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;text-align:right;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;padding:25px 0 0 20px;text-align:right\">\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;background-color:#fff;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:600px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"30\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;border-bottom:#f4f4f4 1px solid;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:left\">\n                                    <h1 style=\"Margin:0 0 8px 0;color:#121212;font-family:helvetica,arial;font-size:23px;font-weight:700;line-height:1.3;text-align:left\">Восстановление пароля на ChartDrug </h1>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Вы запросили данные для восстановлению пароля, данные для доступа в личный кабинет: </p>\n                                    <p style=\"Margin:30px 0 0;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Логин: <b>{login}</b>\n                                    </p>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Ваш пароль: <b>{password}</b>\n                                    </p>\n                                    <div style=\"Margin:30px 0 30px;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:center\">\n                                        <a href=\"https://chartdrug.com\" style=\"background-color:#2b87db;border:#2b87db 11px solid;border-radius:25px;color:#fff!important;display:inline-block;font-family:helvetica;font-size:14px;font-weight:700;padding:0 15px;text-decoration:none\" target=\"_blank\">\n                                            <span style=\"color:#fff\">Перейти в личный кабинет</span>\n                                        </a>\n                                    </div>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Если вы не регистрировались на <a href=\"https://chartdrug.com\" style=\"color:#2b87db!important\" target=\"_blank\">\n                                        <span style=\"color:#2b87db\">ChartDrug</span>\n                                    </a>, просто удалите это письмо. </p>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;padding-top:10px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:580px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#999999;font-family:helvetica,arial;font-size:11px;line-height:14px;text-align:center\">\n                                    <b>© 2022 Chart Drug</b>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            </tbody>\n        </table>\n    </div>\n</div>"

	body = strings.Replace(body, "{login}", user.Login, 1)
	body = strings.Replace(body, "{password}", user.Passwd, 1)
	err = utils.SendMailGmail(user.Email, "Восстановление пароля ChartDrug", body)

	if err != nil {
		return c.WriteWithStatus(err.Error(), http.StatusInternalServerError)
	}

	return c.WriteWithStatus("", http.StatusAccepted)
}
