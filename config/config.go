package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort          string
	Postgres          PostgresConfig
	Smtp              Smtp
	Redis             Redis
	TokenSymmetricKey string
	ServerDomain      string
	Environment       string
	WsPort            string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Smtp struct {
	Sender   string
	Password string
}

type Redis struct {
	Addr string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env") // load .env file if it exists

	v := viper.New()
	v.AutomaticEnv()

	cfg := Config{
		Environment: cast.ToString(returnEnv(v, "ENVIRONMENT", "dev")),
		HttpPort:    cast.ToString(returnEnv(v, "HTTP_PORT", "8000")),
		Postgres: PostgresConfig{
			Host:     cast.ToString(returnEnv(v, "POSTGRES_HOST", "localhost")),
			Port:     cast.ToString(returnEnv(v, "POSTGRES_PORT", "5432")),
			User:     cast.ToString(returnEnv(v, "POSTGRES_USER", "postgres")),
			Password: cast.ToString(returnEnv(v, "POSTGRES_PASSWORD", "password")),
			Database: cast.ToString(returnEnv(v, "POSTGRES_DATABASE", "postgres")),
		},
		Smtp: Smtp{
			Sender:   cast.ToString(returnEnv(v, "SMTP_SENDER", "smth@gmail")),
			Password: cast.ToString(returnEnv(v, "SMTP_PASSWORD", "pass")),
		},
		Redis: Redis{
			Addr: cast.ToString(returnEnv(v, "REDIS_ADDR", "localhost:6379")),
		},
		TokenSymmetricKey: cast.ToString(returnEnv(v, "TOKEN_SYMMETRIC_KEY", "key")),
		WsPort:            cast.ToString(returnEnv(v, "WS_PORT", ":5001")),
	}

	return cfg
}

func returnEnv(v *viper.Viper, key string, defaultValue any) any {
	value := v.Get(key)
	if value != nil {
		return value
	}

	return defaultValue
}
