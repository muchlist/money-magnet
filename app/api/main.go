package main

import (
	"flag"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
	"github.com/muchlist/moneymagnet/bussines/sys/validate"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/web"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		minOpenConns int
	}
}

type application struct {
	config config
	logger mlogger.Logger
	db     *pgxpool.Pool
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:postgres@172.24.48.1:5432/test_db?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max", 100, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.minOpenConns, "db-min", 1, "PostgreSQL min open connections")

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

	// new api server
	webApi := web.New(app.logger, app.config.port, app.config.env)
	err = webApi.Serve(app.routes())
	if err != nil {
		log.Error("serve web api", err)
	}
}
