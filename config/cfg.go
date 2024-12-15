package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	SSLMode    string

	RedisHost     string
	RedisPort     string
	RedisDB       string
	RedisProtocol string
	RedisPassword string

	MongoURI      string
	MongoUsername string
	MongoPassword string
	MongoDatabase string

	HMACSecret string
}

func ReadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DBUsername:    os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASS"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		SSLMode:       os.Getenv("SSL_MODE"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisDB:       Default(os.Getenv("REDIS_DB"), "0"),
		RedisProtocol: Default(os.Getenv("REDIS_PROTOCOL"), "3"),
		MongoURI:      os.Getenv("MONGO_URI"),
		MongoUsername: os.Getenv("MONGO_USERNAME"),
		MongoPassword: os.Getenv("MONGO_PASSWORD"),
		MongoDatabase: os.Getenv("MONGO_DB"),
		HMACSecret:    os.Getenv("HMAC_SECRET"),
	}
}

func Default(env string, defaultValue string) string {
	if env == "" {
		return defaultValue
	}
	return env
}
