package postgres

import (
	"fmt"
	"time"

	"github.com/phpCoder88/url-shortener-observable/internal/config"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

const (
	maxOpenConnections = 60
	connMaxLifetime    = 120
	maxIdleConnections = 30
	connMaxIdleTime    = 20
)

func NewPgConnection(dbConf *config.DBConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		dbConf.Host,
		dbConf.Port,
		dbConf.User,
		dbConf.Name,
		dbConf.Password,
		dbConf.SSLMode,
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func GetConnectionString(dbConf *config.DBConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConf.User,
		dbConf.Password,
		dbConf.Host,
		dbConf.Port,
		dbConf.Name,
		dbConf.SSLMode,
	)
}
