package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBAddress              string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Load env error: ", err)
	}

	return Config{
		PublicHost: getEnv("PUBLIC_HOST", "http://localhost"),
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "DBPassword"),
		DBName:     getEnv("DB_NAME", "ecom"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),

		JWTSecret:              getEnv("JWT_SECRET", "Dev123456"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRE", 3600*24*7),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			return fallback
		}
		return int64(valueInt)
	}
	return fallback
}
