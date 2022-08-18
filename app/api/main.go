package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
	"github.com/muchlist/moneymagnet/bussines/sys/validate"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/web"
)

const version = "1.0.0"

type config struct {
	port      int
	debugPort int
	env       string
	db        struct {
		dsn          string
		maxOpenConns int
		minOpenConns int
	}
	secret string
}

type application struct {
	config config
	logger mlogger.Logger
	db     *pgxpool.Pool
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "Api server port")
	flag.IntVar(&cfg.debugPort, "debug-port", 4000, "Debug server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@localhost:5432/money_magnet?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max", 100, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minOpenConns, "db-min", 1, "PostgreSQL min open connections")
	flag.StringVar(&cfg.secret, "secret", "xoxoxoxo", "jwt secret")

	flag.Parse()

	// init log
	log := mlogger.New("info", "stdout")

	// init database
	database, err := db.OpenDB(db.Config{
		DSN:          cfg.db.dsn,
		MaxOpenConns: int32(cfg.db.maxOpenConns),
		MinOpenConns: int32(cfg.db.minOpenConns),
	})
	if err != nil {
		log.Error("connection to database", err)
		panic(err.Error())
	}
	defer database.Close()

	// init validator
	validate.Init()

	// init application
	app := application{
		config: cfg,
		logger: log,
		db:     database,
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
