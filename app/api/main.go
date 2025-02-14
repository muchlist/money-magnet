package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/muchlist/moneymagnet/cfg"
	"github.com/muchlist/moneymagnet/pkg/cache"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/global"
	"github.com/muchlist/moneymagnet/pkg/mfirebase"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/observ/mmetric"
	"github.com/redis/go-redis/v9"

	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/jackc/pgx/v5/pgxpool"
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
	firebase  *firebase.App
	redis     *redis.Client
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
	contextField := map[string]any{"request_id": global.RequestIDKey}
	if config.Toggle.TraceON {
		contextField["trace_id"] = global.TraceIDKey
	}
	log := mlogger.New(mlogger.Options{
		Level:        mlogger.LevelInfo,
		Output:       config.App.LoggerOutput,
		ContextField: contextField,
	})

	// Set Tracer and Metrics Open Telemetry
	ctxWC, cancle := context.WithCancel(ctx)
	defer cancle()
	otelCfg := observ.Option{
		ServiceName:  config.App.Name,
		CollectorURL: config.Telemetry.URL,
		Headers:      map[string]string{"api-key": config.Telemetry.Key},
		Insecure:     config.Telemetry.Insecure,
	}
	if config.Toggle.TraceON {
		cleanUp := observ.InitTracer(ctx, otelCfg, log)
		defer cleanUp(ctx)
	}
	if config.Toggle.MetricON {
		cleanUpMetric := observ.InitMeter(ctx, otelCfg, log)
		defer cleanUpMetric(ctx)
		// register monitor metrics
		mmetric.RegisterMonitorMetric(ctxWC)
	}

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

	// init redis
	redis := cache.InitRedis(config)
	defer func() {
		if redis != nil {
			redis.Close()
		}
	}()

	// init firebase app
	firebaseApp, err := mfirebase.InitFirebase(mfirebase.Config{
		CredLocation: config.Google.CredentialLocation,
	})
	if err != nil {
		log.Error("init firebase", err)
		panic(err.Error())
	}

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
		firebase:  firebaseApp,
		redis:     redis,
	}

	// start debug server
	if config.App.DebugPort != 0 {
		debugMux := debugMux(database)
		go func(mux *http.ServeMux) {
			if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", config.App.DebugPort), mux); err != nil {
				log.Error("serve debug api", err)
			}
		}(debugMux)
	}

	// create and start api server
	webApi := web.New(app.logger, config.App.Port, config.App.Env, config.App.Name)
	routes, err := app.routes()
	if err != nil {
		log.Error("generate routes", err)
		panic(err.Error())
	}

	err = webApi.Serve(routes)
	if err != nil {
		log.Error("serve web api", err)
	}
}
