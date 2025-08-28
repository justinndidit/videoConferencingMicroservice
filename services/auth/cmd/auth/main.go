package main

import (
	"os"

	configPkg "github.com/justinndidit/videoConferencingMicroservice/internal/config"
	"github.com/justinndidit/videoConferencingMicroservice/pkg/logger"
	"github.com/justinndidit/videoConferencingMicroservice/services/auth/internal/application"
	configAuth "github.com/justinndidit/videoConferencingMicroservice/services/auth/internal/config"
	"github.com/justinndidit/videoConferencingMicroservice/services/auth/internal/store"
)

func main() {

	err := configPkg.LoadEnvironmentVariables("./services/auth/.env.auth")

	if err != nil {
		panic(err)
	}

	app := application.NewApplication()
	envVar := application.NewEnvironmentVariable()

	envVar.DB_Host = os.Getenv("DB_HOST")
	envVar.DB_Name = os.Getenv("DB_DATABASE_NAME")
	envVar.DB_Password = os.Getenv("DB_PASSWORD")
	envVar.DB_Username = os.Getenv("DB_USERNAME")
	envVar.DB_Port = os.Getenv("DB_PORT")
	envVar.Version = os.Getenv("VERSION")
	envVar.ServiceName = os.Getenv("SERVICE_NAME")
	envVar.Environment = os.Getenv("ENVIRONMENT")

	var lv string

	if envVar.Environment == "development" {
		lv = "debug"
	} else {
		lv = "info"
	}

	envcfg := logger.LogConfig{
		Environment: envVar.Environment,
		ServiceName: envVar.ServiceName,
		Version:     envVar.Version,
		Level:       lv,
	}

	app.EnvironmentVariables = *envVar

	authLogger, err := configAuth.InitializeAuthLogger(envcfg)

	if err != nil {
		panic(err)
	}

	app.Logger = authLogger

	err = store.OpenDatabaseConnection(app)

	if err != nil {
		authLogger.Error("Error opening database")
		panic(err)
	}
	app.Logger.Sugar().Infof("Successfully Initialized app: %v", *app)
}
