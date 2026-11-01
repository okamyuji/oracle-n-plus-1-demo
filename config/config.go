package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/sijms/go-ora/v3"
)

// Config - アプリケーション設定
type Config struct {
	DBHost        string
	DBPort        int
	DBServiceName string
	DBUsername    string
	DBPassword    string

	// Redis設定（オプション）
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

// LoadConfig - 設定を読み込む
func LoadConfig() (*Config, error) {
	// .envファイルを読み込む（存在する場合）
	_ = godotenv.Load()

	config := &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBServiceName: getEnv("DB_SERVICE_NAME", "ORCLPDB1"),
		DBUsername:    getEnv("DB_USERNAME", ""),
		DBPassword:    getEnv("DB_PASSWORD", ""),

		// Redis設定（オプション）
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
	}

	// DBポート番号の解析
	portStr := getEnv("DB_PORT", "1521")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	config.DBPort = port

	// Redisポート番号の解析
	redisPortStr := getEnv("REDIS_PORT", "6379")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
	}
	config.RedisPort = redisPort

	// 必須項目のチェック
	if config.DBUsername == "" {
		return nil, fmt.Errorf("DB_USERNAME is required")
	}
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return config, nil
}

// OracleDSN - go-ora の接続文字列を返す
func OracleDSN(config *Config) string {
	// go-ora v3.0.1 の FastLogin は、2 本目以降の接続で Oracle 23c に EOF で切られる（sijms/go-ora#732）
	return fmt.Sprintf("oracle://%s:%s@%s:%d/%s?FAST+LOGIN=false",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBServiceName,
	)
}

// ConnectDatabase - データベースに接続する
func ConnectDatabase(config *Config) (*sql.DB, error) {
	db, err := sql.Open("oracle", OracleDSN(config))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}

// getEnv - 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
