package databaseconnection

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/justinndidit/videoConferencingMicroservice/pkg/logger"
)

type postgres struct {
	dbConfig *DatabaseConfig
	db       *sql.DB
	logger   *logger.ContextLogger
}

func NewPosgres(dbConfig *DatabaseConfig) Database {
	loggerConfig := logger.NewLogConfig()
	loggerConfig.Environment = "development"
	loggerConfig.Level = "debug"
	loggerConfig.ServiceName = "Postgres database package"
	loggerConfig.Version = "1.0"

	logger, err := logger.NewLogger(*loggerConfig)
	if err != nil {
		fmt.Printf("Error initializing logger: %v", err)
		return nil
	}
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
	return &postgres{
		dbConfig: dbConfig,
		logger:   logger,
	}
}

func (pg *postgres) BuildDatabaseUrl() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pg.dbConfig.Host, pg.dbConfig.Username, pg.dbConfig.Password, pg.dbConfig.DatabaseName, pg.dbConfig.Port)
}

func (pg *postgres) OpenDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", pg.BuildDatabaseUrl())

	if err != nil {
		pg.logger.Error("Failed to Open database connection!")
		return nil, fmt.Errorf("Failed to Open database connectio: %w", err)
	}

	db.SetMaxOpenConns(pg.dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(pg.dbConfig.MaxIdleConns)
	db.SetConnMaxLifetime(pg.dbConfig.ConnMaxLife)

	err = db.Ping()

	if err != nil {
		pg.logger.Error("Error pinging database")
		if pg.db != nil {
			pg.logger.Info("Closing database...")
			err = pg.db.Close()
		}
		if err != nil {
			pg.logger.Fatal("Error closing database connection!")
			panic(err)
		}
		return nil, fmt.Errorf("Failed to Ping postgres database: %w", err)
	}

	pg.db = db
	return db, nil
}

func (pg postgres) Close() error {
	if pg.db != nil {
		return pg.db.Close()
	}
	pg.logger.Info("Database connection not established!")
	return nil
}

func (pg postgres) Ping() error {
	if pg.db == nil {
		return fmt.Errorf("Postgres connection not established yet")
	}
	pg.logger.Info("Pinging database...")
	return pg.db.Ping()
}
