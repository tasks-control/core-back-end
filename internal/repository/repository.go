package repository

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const maxOpenConns = 10

type Config struct {
	DBUserEnv     string `validate:"required" yaml:"dbUserEnv"`
	DBPasswordEnv string `validate:"required" yaml:"dbPasswordEnv"`
	DBNameEnv     string `validate:"required" yaml:"dbNameEnv"`
	DBHostEnv     string `validate:"required" yaml:"dbHostEnv"`
	DBPortEnv     string `validate:"required" yaml:"dbPortEnv"`
}

type Repository interface {
}

type repository struct {
	conn *sqlx.DB
}

func New(cfg Config) (Repository, error) {
	DBUser := os.Getenv(cfg.DBUserEnv)
	DBPassword := os.Getenv(cfg.DBPasswordEnv)
	DBName := os.Getenv(cfg.DBNameEnv)
	DBHost := os.Getenv(cfg.DBHostEnv)
	DBPort := os.Getenv(cfg.DBPortEnv)

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUser, DBPassword, DBHost, DBPort, DBName)
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(maxOpenConns)
	postgres := sqlx.NewDb(dbConn, "postgres")

	if err := postgres.Ping(); err != nil {
		return nil, err
	}

	return &repository{
		conn: postgres,
	}, nil
}
