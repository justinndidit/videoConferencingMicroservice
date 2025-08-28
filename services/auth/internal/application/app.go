package application

import (
	"database/sql"

	// "github.com/justinndidit/videoConferencingMicroservice/internal/config"

	"github.com/justinndidit/videoConferencingMicroservice/pkg/logger"
)

type EnvironmentVariable struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Name     string
	DB_Host     string
	Version     string
	ServiceName string
	Environment string
}

func NewEnvironmentVariable() *EnvironmentVariable {
	return &EnvironmentVariable{}
}

type Application struct {
	Logger               *logger.ContextLogger
	EnvironmentVariables EnvironmentVariable
	DB                   *sql.DB
}

func NewApplication() *Application {
	return &Application{}
}
