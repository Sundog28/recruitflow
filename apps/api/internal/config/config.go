package config

import (
"log"
"os"
)

type Config struct {
Port      string
DBURL     string
JWTSecret string
}

func Load() Config {
cfg := Config{
Port:      getenv("PORT", "8080"),
DBURL:     getenv("DATABASE_URL", "postgres://recruitflow:recruitflow@localhost:5432/recruitflow?sslmode=disable"),
JWTSecret: getenv("JWT_SECRET", "dev_secret_change_me"),
}

if cfg.JWTSecret == "dev_secret_change_me" {
log.Println("[warn] JWT_SECRET is using default dev value; set env var in production.")
}
return cfg
}

func getenv(key, fallback string) string {
v := os.Getenv(key)
if v == "" {
return fallback
}
return v
}
