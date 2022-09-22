package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/global"
	"github.com/muchlist/moneymagnet/pkg/observ"

	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/muchlist/moneymagnet/docs"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

const version = "1.0.0"

type config struct {
	applicationName string
	port            int
	debugPort       int
	env             string
	db              struct {
		dsn         string
		maxOpenCons int
		minOpenCons int
	}
	secret            string
	collectorURL      string
	collectorKey      string
	collectorInsecure bool
}

type application struct {
	config    config
	logger    mlogger.Logger
	validator validate.Validator
	db        *pgxpool.Pool
}

// @title Money Magnet API
// @version 1.0
// @description this is server for money magnet application.
// @termsOfService http://swagger.io/terms/

// @contact.name Muchlis
// @contact.url https://muchlis.dev
// @contact.email whois.muchlis@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
func main() {
	var cfg config

	flag.StringVar(&cfg.applicationName, "name", "money-magnet", "Application Name")
	flag.IntVar(&cfg.port, "port", 8081, "Api server port")
	flag.IntVar(&cfg.debugPort, "debug-port", 4000, "Debug server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/money_magnet?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max", 100, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minOpenCons, "db-min", 1, "PostgreSQL min open connections")
	flag.StringVar(&cfg.secret, "secret", "xoxoxoxo", "jwt secret")
	flag.StringVar(&cfg.collectorURL, "otel-url", "localhost:4317", "open telemetry collector url")
	flag.StringVar(&cfg.collectorKey, "otel-key", "example-api-key", "open telemetry api-key")
	flag.BoolVar(&cfg.collectorInsecure, "otel-insecure", true, "open telemetry insecure")

	flag.Parse()

	ctx := context.Background()

	// init log
	log := mlogger.New(mlogger.Options{
		Level:  mlogger.LevelInfo,
		Output: "stdout",
		ContextField: map[string]any{
			"request_id": global.RequestIDKey,
			"trace_id":   global.TraceIDKey,
		},
	})

	// Set Tracer and Metrics Open Telemetry
	otelCfg := observ.Option{
		ServiceName:  cfg.applicationName,
		CollectorURL: cfg.collectorURL,
		ApiKey:       cfg.collectorKey,
		Insecure:     cfg.collectorInsecure,
	}
	cleanUp := observ.InitTracer(ctx, otelCfg, log)
	defer cleanUp(ctx)
	cleanUpMetric := observ.InitMeter(ctx, otelCfg, log)
	defer cleanUpMetric(ctx)

	// init database
	database, err := db.OpenDB(db.Config{
		DSN:          cfg.db.dsn,
		MaxOpenConns: int32(cfg.db.maxOpenCons),
		MinOpenConns: int32(cfg.db.minOpenCons),
	})
	if err != nil {
		log.Error("connection to database", err)
		panic(err.Error())
	}
	defer database.Close()

	// init validator
	validateRegists := []validate.Register{
		{
			Key:       "custom_date",
			Translate: "{0} must be valid date format",
			ValidFunc: func(fl validator.FieldLevel) bool {
				str := fl.Field().String()
				layout := "2006-01-02 15:04:05"
				_, err := time.Parse(layout, str)
				return err == nil
			},
		},
	}
	validatorInst := validate.New(validateRegists...)

	// init application
	app := application{
		config:    cfg,
		logger:    log,
		validator: validatorInst,
		db:        database,
	}

	// start debug server
	debugMux := debugMux(database)
	go func(mux *http.ServeMux) {
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", cfg.debugPort), mux); err != nil {
			log.Error("serve debug api", err)
		}
	}(debugMux)

	// create and start api server
	webApi := web.New(app.logger, app.config.port, app.config.env, app.config.applicationName)
	err = webApi.Serve(app.routes())
	if err != nil {
		log.Error("serve web api", err)
	}
}
