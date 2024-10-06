package database

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/nurcahyaari/ecommerce/config"
	"github.com/rs/zerolog/log"
)

type SQLDatabase struct {
	DB *sqlx.DB
	Tx *sqlx.Tx
}

type SQLDatabaseTrx[T any] interface {
	BeginTx(ctx context.Context) (T, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func (db *SQLDatabase) RunMigrations(cfg config.Config) error {
	driver, err := mysql.WithInstance(db.DB.DB, &mysql.Config{
		DatabaseName: cfg.DB.MySQL.Name,
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migration",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func NewMysql(cfg config.Config) *SQLDatabase {
	log.Info().Msg("Initialize Mysql connection")
	var err error

	dbHost := cfg.DB.MySQL.Host
	dbPort := cfg.DB.MySQL.Port
	dbName := cfg.DB.MySQL.Name
	dbUser := cfg.DB.MySQL.User
	dbPass := cfg.DB.MySQL.Pass
	maxPoolSize := 10
	if cfg.DB.MongoDB.MaxPoolSize != 0 {
		maxPoolSize = cfg.DB.MongoDB.MaxPoolSize
	}

	sHost := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sqlx.Connect("mysql", sHost)
	if err != nil {
		log.Fatal().Msgf("Error to loading Database %s", err)
	}

	db.SetMaxOpenConns(maxPoolSize)

	log.Info().Str("Name", dbName).Msg("Success connect to DB")
	return &SQLDatabase{
		DB: db,
	}
}
