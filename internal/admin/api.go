package admin

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"strings"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)

	// the following endpoints require a valid JWT

	r.Get("/api/users", res.get)
	r.Get("/api/drugs", res.getDrugs)

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
		println("-------")
		//println(claims["admin"].(bool))
		if claims["admin"] == nil || !claims["admin"].(bool) {
			return errors.Forbidden("You are not an admin")
		}
		users, err := r.service.Get(c.Request.Context())
		//fmt.Println(profile.TypeSports)
		if err != nil {
			return err
		}
		return c.Write(users)

		/*
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
			}*/

		//TypeSports     []string  `json:"typesports"`
		//return c.Write(profileRest)
	} else {
		return err
	}

}

func (r resource) getDrugs(c *routing.Context) error {
	reqToken := strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1]

	token, _, err := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		println("-------")
		if claims["admin"] == nil || !claims["admin"].(bool) {
			return errors.Forbidden("You are not an admin")
		}

		drug := utils.GetDrugs()
		drugs := []utils.Drug{}

		for _, item := range drug {
			drugs = append(drugs, item)
		}

		return c.Write(drugs)
	} else {
		return err
	}
}
