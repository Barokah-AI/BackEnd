package config

import (
	// "log"
	"os"
	// "github.com/joho/godotenv"
)

func GetEnv(envName string) string {
	// envFile, _ := godotenv.Read("../.env")
	// return envFile[envName]
	return os.Getenv(envName)
}

// // LoadEnv loads the environment variables from a .env file
// func LoadEnv() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}
// }

// // GetEnv gets the environment variable by name
// func GetEnv(envName string, defaultValue string) string {
// 	value, exists := os.LookupEnv(envName)
// 	if !exists {
// 		return defaultValue
// 	}
// 	return value
// }

// func init() {
// 	LoadEnv()
// }
