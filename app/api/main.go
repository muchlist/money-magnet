package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muchlist/moneymagnet/cfg"
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

type application struct {
	config    *cfg.Config
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
	config := cfg.Load()
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
		ServiceName:  config.App.Name,
		CollectorURL: config.Telemetry.URL,
		ApiKey:       config.Telemetry.Key,
		Insecure:     config.Telemetry.Insecure,
	}
	cleanUp := observ.InitTracer(ctx, otelCfg, log)
	defer cleanUp(ctx)
	cleanUpMetric := observ.InitMeter(ctx, otelCfg, log)
	defer cleanUpMetric(ctx)

	// init database
	database, err := db.OpenDB(db.Config{
		DSN:          config.DB.DSN,
		MaxOpenConns: config.DB.MaxOpenCons,
		MinOpenConns: config.DB.MinOpenCons,
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
		config:    config,
		logger:    log,
		validator: validatorInst,
		db:        database,
	}

	// start debug server
	debugMux := debugMux(database)
	go func(mux *http.ServeMux) {
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", config.App.DebugPort), mux); err != nil {
			log.Error("serve debug api", err)
		}
	}(debugMux)

	// create and start api server
	webApi := web.New(app.logger, config.App.Port, config.App.Env, config.App.Name)
	err = webApi.Serve(app.routes())
	if err != nil {
		log.Error("serve web api", err)
	}
}
