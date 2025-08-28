package config

import (
	"fmt"

	"github.com/lpernett/godotenv"
)

type EnvironmentVariables struct {
}

func NewEnvironmentVariables() *EnvironmentVariables {
	return &EnvironmentVariables{}
}

func LoadEnvironmentVariables(envName string) error {
	err := godotenv.Load(envName)

	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return err
	}
	return nil
}
