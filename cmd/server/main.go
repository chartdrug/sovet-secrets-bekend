package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/lib/pq"
	"github.com/qiangxue/sovet-secrets-bekend/internal/admin"
	"github.com/qiangxue/sovet-secrets-bekend/internal/album"
	"github.com/qiangxue/sovet-secrets-bekend/internal/antros"
	"github.com/qiangxue/sovet-secrets-bekend/internal/auth"
	"github.com/qiangxue/sovet-secrets-bekend/internal/config"
	"github.com/qiangxue/sovet-secrets-bekend/internal/course"
	"github.com/qiangxue/sovet-secrets-bekend/internal/errors"
	"github.com/qiangxue/sovet-secrets-bekend/internal/healthcheck"
	"github.com/qiangxue/sovet-secrets-bekend/internal/injections"
	"github.com/qiangxue/sovet-secrets-bekend/internal/profile"
	"github.com/qiangxue/sovet-secrets-bekend/internal/spr"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/accesslog"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/dbcontext"
	"github.com/qiangxue/sovet-secrets-bekend/pkg/log"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"time"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	//db, err := dbx.MustOpen("postgres", cfg.DSN)
	db, err := dbx.Open("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	db2, err1 := sql.Open("postgres", cfg.DSN)
	if err1 != nil {
		logger.Error(err1)
		os.Exit(-1)
	}

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), db2, cfg),
	}

	c := cron.New()
	/*
		c.AddFunc("@every 10s", func() {
			//fmt.Println(time.Now())
			response, err := http.Get("http://localhost:8080/v1/api/injectionsCall")

			if err != nil {
				logger.Error(err.Error())
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logger.Error(err)
			}
			logger.Info(string(responseData))
		})
	*/

	//body := "<div style=\"background-color:#f4f4f4;margin:0;padding:0\">\n    <div style=\"Margin:0;background-color:#f4f4f4;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;table-layout:fixed;text-align:center\">\n        <table class=\"m_7451640406457273268container\" cellpadding=\"0\" align=\"center\" style=\"border-collapse:collapse;margin:0 auto;max-width:600px;width:100%\">\n            <tbody>\n            <tr>\n                <td style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:0;line-height:18px;padding-bottom:26px;text-align:right\">\n                    <div class=\"m_7451640406457273268accessibilitylink\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;text-align:right;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#585858;font-family:helvetica,arial;font-size:11px;line-height:18px;padding:25px 0 0 20px;text-align:right\">\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;background-color:#fff;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:600px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"30\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;border-bottom:#f4f4f4 1px solid;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:left\">\n                                    <h1 style=\"Margin:0 0 8px 0;color:#121212;font-family:helvetica,arial;font-size:23px;font-weight:700;line-height:1.3;text-align:left\">Спасибо за <span class=\"il\">регистрацию</span> на ChartDrug </h1>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Ваша <span class=\"il\">регистрация</span> успешно завершена, данные для доступа в личный кабинет: </p>\n                                    <p style=\"Margin:30px 0 0;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Логин: <b>+user.Login+</b>\n                                    </p>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Ваш пароль: <b>+user.Passwd+</b>\n                                    </p>\n                                    <div style=\"Margin:30px 0 30px;color:#000;font-family:helvetica,arial;font-size:13px;line-height:18px;text-align:center\">\n                                        <a href=\"https://chartdrug.com\" style=\"background-color:#2b87db;border:#2b87db 11px solid;border-radius:25px;color:#fff!important;display:inline-block;font-family:helvetica;font-size:14px;font-weight:700;padding:0 15px;text-decoration:none\" target=\"_blank\">\n                                            <span style=\"color:#fff\">Перейти в личный кабинет</span>\n                                        </a>\n                                    </div>\n                                    <p style=\"Margin:0 0 15px;color:#121212;font-family:helvetica,arial;font-size:15px;line-height:1.47;text-align:left\">Если вы не регистрировались на <a href=\"https://chartdrug.com\" style=\"color:#2b87db!important\" target=\"_blank\">\n                                        <span style=\"color:#2b87db\">ChartDrug</span>\n                                    </a>, просто удалите это письмо. </p>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            <tr>\n                <td style=\"Margin:0;color:#000;font-family:helvetica,arial;font-size:0;line-height:18px;padding-top:10px;text-align:left!important\">\n                    <div style=\"Margin:0;color:#000;display:inline-block;font-family:helvetica,arial;font-size:14px;line-height:18px;max-width:580px;text-align:left;vertical-align:top;width:100%\">\n                        <table width=\"100%\" cellpadding=\"20\" style=\"border-collapse:collapse\">\n                            <tbody>\n                            <tr>\n                                <td class=\"m_7451640406457273268cell\" style=\"Margin:0;color:#999999;font-family:helvetica,arial;font-size:11px;line-height:14px;text-align:center\">\n                                    <b>© 2022 Chart Drug</b>\n                                </td>\n                            </tr>\n                            </tbody>\n                        </table>\n                    </div>\n                </td>\n            </tr>\n            </tbody>\n        </table>\n    </div>\n</div>"
	//err = utils.SendMailGmail("k.dereshev@gmail.com", "Успешная регистрация на ChartDrug", body)
	c.Start()

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}

}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB, db2 *sql.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	authHandler := auth.Handler(cfg.JWTSigningKey)

	album.RegisterHandlers(rg.Group(""),
		album.NewService(album.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger, auth.NewRepository(db, logger)),
		logger,
	)

	profile.RegisterHandlers(rg.Group(""),
		profile.NewService(profile.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	antros.RegisterHandlers(rg.Group(""),
		antros.NewService(antros.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	injections.RegisterHandlers(rg.Group(""),
		injections.NewService(injections.NewRepository(db, db2, logger), logger),
		authHandler, logger,
	)

	spr.RegisterHandlers(rg.Group(""),
		spr.NewService(spr.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	admin.RegisterHandlers(rg.Group(""),
		admin.NewService(admin.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	course.RegisterHandlers(rg.Group(""),
		course.NewService(course.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
