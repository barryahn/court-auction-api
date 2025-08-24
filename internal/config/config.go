package config

import (
	"fmt"
	"os"
)

// Config는 애플리케이션의 런타임 설정을 보관합니다.
type Config struct {
    // API 서버 포트 (예: ":8080")
    APIServerAddr string

    // PostgreSQL DSN (예: "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
    PostgresDSN string

    // Redis 주소 (예: "localhost:6379")
    RedisAddr string
}

// Load는 환경변수에서 설정을 읽어와 Config를 생성합니다.
// 필수값 누락 시 에러를 반환합니다.
func Load() (*Config, error) {
    cfg := &Config{
        APIServerAddr: getEnv("API_ADDR", ":8080"),
        PostgresDSN:   os.Getenv("POSTGRES_DSN"),
        RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
    }

    if cfg.PostgresDSN == "" {
        return nil, fmt.Errorf("POSTGRES_DSN is required")
    }

    return cfg, nil
}

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}


