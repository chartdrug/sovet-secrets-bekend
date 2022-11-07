package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/qiangxue/sovet-secrets-bekend/internal/payment"
	"github.com/segmentio/kafka-go"
	"io/ioutil"
	"strings"

	//"github.com/Shopify/sarama"
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
	//"gopkg.in/alecthomas/kingpin.v2"
	//"io/ioutil"
	"net/http"
	"os"
	//"os/signal"
	"time"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
var brokerList = flag.String("brokerList", "89.208.219.91:9092", "broker list")

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

	// start the HTTP server with graceful shutdown

	go func() {
		brokers := *brokerList
		config := kafka.ReaderConfig{
			//Brokers:         []string{"89.208.219.91:9092"},
			Brokers: strings.Split(brokers, ","),
			//GroupID:         kafkaClientId,
			Topic:           "calc_injection",
			MinBytes:        10e3,            // 10KB
			MaxBytes:        10e6,            // 10MB
			MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
			ReadLagInterval: -1,
		}

		reader := kafka.NewReader(config)
		defer reader.Close()

		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				logger.Error("error while receiving message: %s", err.Error())
				continue
			}

			value := m.Value

			if err != nil {
				logger.Error("error while receiving message: %s", err.Error())
				continue
			}

			logger.Infof("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(value))

			response, err := http.Get("http://localhost:8080/v1/api/injectionsCall/" + string(value))

			if err != nil {
				logger.Error(err.Error())
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				logger.Error(err)
			}
			logger.Info(string(responseData))
		}
	}()

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

	payment.RegisterHandlers(rg.Group(""),
		payment.NewService(payment.NewRepository(db, logger), logger),
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
