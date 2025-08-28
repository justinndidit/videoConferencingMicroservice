package store

import (
	databaseconnection "github.com/justinndidit/videoConferencingMicroservice/pkg/databaseConnection"
	"github.com/justinndidit/videoConferencingMicroservice/services/auth/internal/application"
)

func InitializeDatabaseConfig(app *application.Application) *databaseconnection.DatabaseConfig {
	return &databaseconnection.DatabaseConfig{
		Username:     app.EnvironmentVariables.DB_Username,
		Password:     app.EnvironmentVariables.DB_Password,
		Host:         app.EnvironmentVariables.DB_Host,
		DatabaseName: app.EnvironmentVariables.DB_Name,
		Port:         app.EnvironmentVariables.DB_Port,
	}
}

func InitializeDatabaseDriver(app *application.Application) (databaseconnection.Database, error) {
	dbDriver := databaseconnection.NewPosgres(InitializeDatabaseConfig(app))

	dbDriver.Logger = app.Logger
	return dbDriver, nil
}

func OpenDatabaseConnection(app *application.Application) error {
	dbConn, err := InitializeDatabaseDriver(app)
	if err != nil {
		app.Logger.Error("Error initializing database driver")
		return err
	}
	db, err := dbConn.OpenDatabaseConnection()

	if err != nil {
		app.Logger.Error("Error Opening database connection")
		return err
	}

	app.DB = db
	return nil
}
