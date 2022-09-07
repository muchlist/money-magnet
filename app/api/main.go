package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

const version = "1.0.0"

type config struct {
	port      int
	debugPort int
	env       string
	db        struct {
		dsn         string
		maxOpenCons int
		minOpenCons int
	}
	secret string
}

type application struct {
	config    config
	logger    mlogger.Logger
	validator validate.Validator
	db        *pgxpool.Pool
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "Api server port")
	flag.IntVar(&cfg.debugPort, "debug-port", 4000, "Debug server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/money_magnet?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max", 100, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minOpenCons, "db-min", 1, "PostgreSQL min open connections")
	flag.StringVar(&cfg.secret, "secret", "xoxoxoxo", "jwt secret")

	flag.Parse()

	// init log
	log := mlogger.New("info", "stdout")

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
	webApi := web.New(app.logger, app.config.port, app.config.env)
	err = webApi.Serve(app.routes())
	if err != nil {
		log.Error("serve web api", err)
	}
}
