package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string, ip string) (string, error)
}

type User struct {
	entity.Users
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetName() string
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	repo            Repository
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, logger log.Logger, repo Repository) Service {
	return service{signingKey, tokenExpiration, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string, ip string) (string, error) {

	user, err := s.authenticate(ctx, username, password, ip)
	if err != nil {
		return "", errors.Unauthorized(err.Error())
	}
	if user.Login != "" {
		return s.generateJWT(user)
	} else {
		return "", errors.Unauthorized("")

	}

}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, username, password string, ip string) (User, error) {
	logger := s.logger.With(ctx, "user", username)

	user, err := s.repo.Get(ctx, username, password)
	if user.GetLogn() != "" && err == nil {
		logger.Infof("authentication OK username=" + username)
		//сохранить дату логина
		errU := s.repo.UpdateTimeLastLogin(ctx, user.ID)
		if err != errU {
			return User{}, errU
		}
		//+ запись истории
		errIp, Country_, Region_, City_ := utils.GetInfo(ip)
		if errIp != nil {
			logger.Error("problem with utils.GetInfo ip=" + ip + ", err=" + errIp.Error())
		}
		errH := s.repo.saveHistoryLogin(ctx, entity.HistoryLogin{IdUser: user.ID, DateEvent: time.Now(), IpAddress: ip, Country: Country_, Region: Region_, City: City_})
		if errH != nil {
			return User{}, errH
		}

		return User{user}, nil
	}

	logger.Error("authentication failed username=" + username + ",err=" + err.Error())
	return User{}, errors.Forbidden("The login or password you entered is incorrect. Please try again.")
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(user User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.GetID(),
		"login": user.GetLogn(),
		"exp":   time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
