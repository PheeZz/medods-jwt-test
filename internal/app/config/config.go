package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var Conf = New()

type Config struct {
	JwtSecretKey             []byte
	MongoURI                 string
	DatabaseName             string
	AccessTokenDuration      time.Duration
	RefreshTokenDuration     time.Duration
	JwtAlgorithm             jwt.SigningMethod
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		JwtSecretKey:             []byte(getEnvStrictString("JWT_SECRET_KEY")),
		MongoURI:                 getEnvStrictString("MONGO_URI"),
		DatabaseName:             getEnvStrictString("DATABASE_NAME"),
		AccessTokenDuration:      getEnvStrictDuration("ACCESS_TOKEN_DURATION_HOURS", time.Hour),
		RefreshTokenDuration:     getEnvStrictDuration("REFRESH_TOKEN_DURATION_HOURS", time.Hour),
		JwtAlgorithm:             jwt.SigningMethodHS512,
	}
}

func getEnvStrictString(envVarName string) string {
	value, exists := os.LookupEnv(envVarName)
	if !exists {
		log.Panic("Environment variable not found: ", envVarName)
	}
	return value
}

func getEnvStrictDuration(envVarName string, duration time.Duration) time.Duration {
	value, exists := os.LookupEnv(envVarName)
	if !exists {
		log.Panic("Environment variable not found: ", envVarName)
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panic("Environment variable is not an integer: ", envVarName)
	}
	return time.Duration(intValue) * duration
}
