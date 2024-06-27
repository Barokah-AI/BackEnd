package config

import (
	"github.com/joho/godotenv"
)

func GetEnv(envName string) string {
	envFile, _ := godotenv.Read("../.env")
	return envFile[envName]
}