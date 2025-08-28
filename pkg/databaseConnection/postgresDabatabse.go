package databaseconnection

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/justinndidit/videoConferencingMicroservice/pkg/logger"
)

type Postgres struct {
	dbConfig *DatabaseConfig
	db       *sql.DB
	Logger   *logger.ContextLogger
}

func NewPosgres(dbConfig *DatabaseConfig) *Postgres {

	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}
	if dbConfig.Port == "" {
		dbConfig.Port = "5432"
	}
	if dbConfig.MaxOpenConns == 0 {
		dbConfig.MaxOpenConns = 25
	}
	if dbConfig.MaxIdleConns == 0 {
		dbConfig.MaxIdleConns = 5
	}
	if dbConfig.ConnMaxLife == 0 {
		dbConfig.ConnMaxLife = 5 * time.Minute
	}
	return &Postgres{
		dbConfig: dbConfig,
	}
}

func (pg *Postgres) BuildDatabaseUrl() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pg.dbConfig.Host, pg.dbConfig.Username, pg.dbConfig.Password, pg.dbConfig.DatabaseName, pg.dbConfig.Port)
}

func (pg *Postgres) OpenDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", pg.BuildDatabaseUrl())

	if err != nil {
		pg.Logger.Error("Failed to Open database connection!")
		return nil, fmt.Errorf("Failed to Open database connectio: %w", err)
	}
	pg.db = db
	defer pg.Close()

	pg.db.SetMaxOpenConns(pg.dbConfig.MaxOpenConns)
	pg.db.SetMaxIdleConns(pg.dbConfig.MaxIdleConns)
	pg.db.SetConnMaxLifetime(pg.dbConfig.ConnMaxLife)

	err = pg.db.Ping()

	if err != nil {
		pg.Logger.Error("Error pinging database")
		if pg.db != nil {
			pg.Logger.Info("Closing database...")
			err = pg.db.Close()
		}
		if err != nil {
			pg.Logger.Fatal("Error closing database connection!")
			panic(err)
		}
		return nil, fmt.Errorf("Failed to Ping Postgres database: %w", err)
	}

	return pg.db, nil
}

func (pg Postgres) Close() error {
	if pg.db != nil {
		return pg.db.Close()
	}
	pg.Logger.Info("Database connection not established!")
	return nil
}

func (pg Postgres) Ping() error {
	if pg.db == nil {
		return fmt.Errorf("Postgres connection not established yet")
	}
	pg.Logger.Info("Pinging database...")
	return pg.db.Ping()
}
