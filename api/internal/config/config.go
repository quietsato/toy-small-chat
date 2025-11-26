package config

import (
	"fmt"
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT"`
	User string `envconfig:"USER"`
	Pass string `envconfig:"PASS"`
	Name string `envconfig:"NAME"`
}

func (d *Database) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.User, d.Pass, d.Host, d.Port, d.Name)
}

type Config struct {
	Database     Database `envconfig:"DATABASE"`
	OtlpEndpoint string   `envconfig:"OTLP_ENDPOINT"`
	JWTSecretKey string   `envconfig:"JWT_SECRET_KEY"`
}

func Load() Config {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		slog.Error("failed to load config from env", slog.Any("err", err))
		panic(err) // 回復手段がないため panic
	}
	return c
}
