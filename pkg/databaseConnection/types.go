package databaseconnection

import (
	"database/sql"
	"time"
)

type DatabaseConfig struct {
	Username     string
	Password     string
	Host         string
	DatabaseName string
	Port         string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
	SSLMode      string
}

type Database interface {
	OpenDatabaseConnection() (*sql.DB, error)
	BuildDatabaseUrl() string
	Close() error
	Ping() error
}
