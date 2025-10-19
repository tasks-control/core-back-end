package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tasks-control/core-back-end/internal/models"
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
	MemberRepository
	TokenRepository
}

type MemberRepository interface {
	CreateMember(ctx context.Context, member *models.Member) error
	GetMemberByEmail(ctx context.Context, email string) (*models.Member, error)
	GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error)
	GetMemberByUsername(ctx context.Context, username string) (*models.Member, error)
	UpdateMember(ctx context.Context, member *models.Member) error
}

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
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
