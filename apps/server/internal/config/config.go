package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrEnvFileNotFound = errors.New("env file not found")
	ErrInvalidAppEnv   = errors.New("invalid app environment")
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func (e Environment) isValidEnv() bool {
	switch e {
	case Development, Production:
		return true
	default:
		return false
	}
}

const envFilePath = ".env"

type Config struct {
	Environment           Environment
	Port                  string
	DBUrl                 string
	DBMaxConn             string
	DBMinConn             string
	DBConnMaxLifetime     string
	DBConnMaxIdleLifetime string
	DBHealthCheckPeriod   string
	ConnectTimeout        string
	EmailServiceBaseURL   string
	EmailServiceToken     string
	EmailProvider         string
	JwtSecret             string
	AccessTokenExpiry     string // e.g., "15m"
	RefreshTokenExpiry    string // e.g., "168h" (7 days)
	LogLevel              string
	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleOIDCIssuer   string
	FrontendOrigin     string
	// AWS S3 Config
	AWSRegion          string
	AWSS3Bucket        string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	PresignedURLExpiry int
	UploadMaxFileSize  int64
}

func Load() (*Config, error) {
	err := godotenv.Load(envFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrEnvFileNotFound
		} else {
			return nil, fmt.Errorf("failed to load the env file : %v\n", err)
		}
	}
	cfg := &Config{
		Port:                  getEnv("PORT"),
		DBUrl:                 getEnv("DB_URL"),
		DBMaxConn:             getEnvOrDefault("DB_MAX_CONN", "10"),
		DBMinConn:             getEnvOrDefault("DB_MIN_CONN", "1"),
		DBConnMaxLifetime:     getEnvOrDefault("DB_CONN_MAX_LIFETIME", "0"),
		DBConnMaxIdleLifetime: getEnvOrDefault("DB_CONN_MAX_IDLE_LIFETIME", "0"),
		DBHealthCheckPeriod:   getEnvOrDefault("DB_HEALTH_CHECK_PERIOD", "0"),
		ConnectTimeout:        getEnvOrDefault("DB_CONNECT_TIMEOUT", "0"),
		EmailServiceBaseURL:   getEnv("EMAIL_SERVICE_BASE_URL"),
		EmailServiceToken:     getEnv("EMAIL_SERVICE_TOKEN"),
		EmailProvider:         getEnv("EMAIL_PROVIDER"),
		JwtSecret:             getEnv("JWT_SECRET"),
		AccessTokenExpiry:     getEnv("ACCESS_TOKEN_EXPIRY"),
		RefreshTokenExpiry:    getEnv("REFRESH_TOKEN_EXPIRY"),
		LogLevel:              getEnv("LOG_LEVEL"),
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:     getEnv("GOOGLE_REDIRECT_URL"),
		GoogleOIDCIssuer:      getEnv("GOOGLE_OIDC_ISSUER"),
		FrontendOrigin:        getEnv("FRONTEND_ORIGIN"),
		AWSRegion:             getEnv("AWS_REGION"),
		AWSS3Bucket:           getEnv("AWS_S3_BUCKET"),
		AWSAccessKeyID:        getEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:    getEnv("AWS_SECRET_ACCESS_KEY"),
		PresignedURLExpiry:    parseInteger(getEnv("PRESIGNED_URL_EXPIRY")),
		UploadMaxFileSize:     parseInteger64(getEnv("UPLOAD_MAX_FILE_SIZE")),
	}

	appEnv := Environment(getEnv("APP_ENV"))
	if !appEnv.isValidEnv() {
		return nil, ErrInvalidAppEnv
	}
	cfg.Environment = appEnv

	return cfg, nil
}

// getEnv returns the value of the environment variable
func getEnv(key string) string {
	val := os.Getenv(key)
	val = strings.TrimSpace(val)
	if val == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	} else {
		return val
	}
}

// getEnvOrDefault returns the value of the environment variable if it is set, otherwise it returns the default value
func getEnvOrDefault(key string, val string) string {
	if val := os.Getenv(key); val == "" {
		return val
	} else {
		return val
	}
}

// parseInteger converts a string to an int
func parseInteger(val string) int {
	val = strings.TrimSpace(val)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return i
}

// parseInteger64 converts a string to an int64
func parseInteger64(val string) int64 {
	val = strings.TrimSpace(val)
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
