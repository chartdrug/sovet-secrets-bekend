package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/utils"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"io/ioutil"
	"net/http"
	"time"
)

// Service encapsulates usecase logic for albums.
type Service interface {
	cryptoPaymentResult(ctx context.Context, header map[string][]string, cryptocloud entity.CryptocloudPostback) error
	cryptoPaymentAll(ctx context.Context, owner string) ([]entity.Cryptocloud, error)
	cryptoPaymentCreate(ctx context.Context, owner jwt.MapClaims) (entity.Cryptocloud, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) cryptoPaymentResult(ctx context.Context, header map[string][]string, cryptocloud entity.CryptocloudPostback) error {
	logger := s.logger.With(ctx, "InvoiceId", cryptocloud.InvoiceId)
	logger.Info(cryptocloud)
	logger.Info(header)

	errCreate := s.repo.CreateCryptocloudPostback(ctx, cryptocloud)
	if errCreate != nil {
		utils.SendMailError("Сryptocloud Create", errCreate.Error())
		return errCreate
	} else {
		j, _ := json.Marshal(cryptocloud)
		errU := s.repo.UpdateCryproInvoice(ctx, cryptocloud.OrderId, "paid", string(j))
		if errU != nil {
			logger.Error("error UpdateCryproInvoice message: %s", errU.Error())
			utils.SendMailError("UpdateCryproInvoice", errU.Error())
		}
	}

	return nil
}

func (s service) cryptoPaymentAll(ctx context.Context, owner string) ([]entity.Cryptocloud, error) {
	logger := s.logger.With(ctx, "owner", owner)

	items, err := s.repo.GetAll(ctx, owner)
	if err != nil {
		return nil, err
	}
	for y := 0; y < len(items); y++ {
		//  по всем созданым пробегаемся и проверяем статус
		if items[y].Statusinvoice == "created" {
			req, errHttp := http.NewRequest("GET", "https://api.cryptocloud.plus/v1/invoice/info?uuid="+items[y].Invoiceid, nil)
			req.Header.Set("Authorization", "Token eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MjUxNiwiZXhwIjo4ODA2NzU1ODU4Nn0.de3C_eld1kCh0Ww2VSUAYd17cIhlpQ1ZpJHoOzOPAO8")

			cli := &http.Client{}
			resp, errHttp := cli.Do(req)

			if errHttp != nil {
				logger.Error("error cryptoPaymentAll http message: %s", errHttp.Error())
				utils.SendMailError("cryptoPaymentAll http", errHttp.Error())
				return nil, errHttp
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			logger.Info("cryptoPaymentAll ststus resp message: %s", resp)
			logger.Info("cryptoPaymentAll ststus body message: %s", string(body))

			if resp.StatusCode == 200 {
				var dat map[string]interface{}
				if err := json.Unmarshal(body, &dat); err != nil {
					logger.Error("error cryptoPaymentAll Unmarshal message: %s", err.Error())
					utils.SendMailError("cryptoPaymentAll Unmarshal", err.Error())

				} else {
					items[y].Statusinvoice = dat["status_invoice"].(string)
					items[y].Dtpaym = time.Now()
					errU := s.repo.UpdateCryproInvoice(ctx, items[y].Id, items[y].Statusinvoice, resp.Status+", "+string(body))
					if errU != nil {
						logger.Error("error UpdateCryproInvoice message: %s", errU.Error())
						utils.SendMailError("UpdateCryproInvoice", errU.Error())
					}
				}

			}

			items[y].Resthttpstatus = resp.Status + ", " + string(body)
		}
	}
	return items, nil
}

func (s service) cryptoPaymentCreate(ctx context.Context, owner jwt.MapClaims) (entity.Cryptocloud, error) {

	logger := s.logger.With(ctx, "owner", owner["id"].(string))

	result := entity.Cryptocloud{}
	result.Id = uuid.New().String()
	result.Owner = owner["id"].(string)
	result.Dt = time.Now()
	result.Shopid = "8owDW46PDIEPfZZg"
	result.Amount = 5.0
	result.Currency = "USD"
	result.Statusinvoice = "created"
	result.Email = owner["email"].(string)

	values := map[string]string{"shop_id": result.Shopid, "amount": fmt.Sprintf("%f", result.Amount),
		"order_id": result.Id, "currency": result.Currency, "email": result.Email}

	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", "https://api.cryptocloud.plus/v1/invoice/create", bytes.NewBuffer(jsonValue))
	//sercre key GC8uWj1MAFYGTQ6ahi15r8puD4qvzxXI3vbb
	req.Header.Set("Authorization", "Token eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MjUxNiwiZXhwIjo4ODA2NzU1ODU4Nn0.de3C_eld1kCh0Ww2VSUAYd17cIhlpQ1ZpJHoOzOPAO8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)

	if err != nil {
		logger.Error("error Сryptocloud http message: %s", err.Error())
		utils.SendMailError("Сryptocloud http", err.Error())
		return result, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	result.Reshtpp = resp.Status + ", " + string(body)

	if resp.StatusCode != 200 {
		logger.Error("error Сryptocloud StatusCode message: %s", resp.Status+", "+string(body))
		utils.SendMailError("Сryptocloud StatusCode", resp.Status+", "+string(body))
		return result, errors.BadRequest(resp.Status + ", " + string(body))
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(body, &dat); err != nil {
		logger.Error("error Сryptocloud Unmarshal message: %s", err.Error())
		utils.SendMailError("Сryptocloud Unmarshal", err.Error())
		return result, err
	}

	result.Status = dat["status"].(string)
	result.Payurl = dat["pay_url"].(string)
	result.Invoiceid = dat["invoice_id"].(string)

	logger.Info("Сryptocloud resp message: %s", resp)
	logger.Info("Сryptocloud result message: %s", result)

	errCreate := s.repo.CreateCryptocloud(ctx, result)
	if errCreate != nil {
		utils.SendMailError("Сryptocloud Create", errCreate.Error())
		return result, errCreate
	}

	return result, nil
}
