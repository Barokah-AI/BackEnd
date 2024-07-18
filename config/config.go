package config

import
// "github.com/joho/godotenv"
"os"

// func GetEnv(envName string) string {
// 	// envFile, _ := godotenv.Read("../.env")
// 	// return envFile[envName]
// 	return os.Getenv(envName)
// }

func Getenv(envName string) string {
	// envFile, _ := godotenv.Read("../.env")
	// return envFile[envName]
	return os.Getenv(envName)
}