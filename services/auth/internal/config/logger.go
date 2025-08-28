package config

import (
	"fmt"

	"github.com/justinndidit/videoConferencingMicroservice/pkg/logger"
)

func InitializeAuthLogger(lc logger.LogConfig) (*logger.ContextLogger, error) {

	lc.InitializeLogConfig("Authentication Service")

	authLogger, err := logger.NewLogger(lc)

	if err != nil {
		fmt.Printf("Error initializing Auth logger: %v\n", err)
		return nil, err
	}
	return authLogger, nil
}
