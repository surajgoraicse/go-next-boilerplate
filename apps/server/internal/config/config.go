package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
	// system
	AppName       string
	AppEnv        Environment
	Port          string
	LogLevel      string
	LogFile       string
	LogFileLevel  string
	LogMaxSizeMB  int
	LogMaxBackups int
	LogMaxAgeDays int

	// Database
	DBHost                string
	DBPort                string
	DBUser                string
	DBPassword            string
	DBName                string
	SSLMode               string
	DBMaxConn             int32
	DBMinConn             int32
	DBConnMaxLifetime     time.Duration
	DBConnMaxIdleLifetime time.Duration
	DBHealthCheckPeriod   time.Duration
	ConnectTimeout        time.Duration

	// Email
	EmailServiceBaseURL        string
	EmailServiceToken          string
	EmailProvider              string
	VerificationEmailExpiry    string
	VerificationEmailRateLimit int
	// Auth
	JwtSecret          string
	AccessTokenExpiry  time.Duration // e.g., "15m"
	RefreshTokenExpiry time.Duration // e.g., "168h" (7 days)

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

	MaxActiveDevices int
}

func Load() (*Config, error) {
	err := godotenv.Overload(envFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrEnvFileNotFound
		} else {
			return nil, fmt.Errorf("failed to load the env file : %v\n", err)
		}
	}
	cfg := &Config{
		// system
		AppName:       getEnv("APP_NAME"),
		Port:          getEnv("PORT"),
		LogLevel:      getEnv("LOG_LEVEL"),
		LogFile:       getEnv("LOG_FILE"),
		LogFileLevel:  getEnv("FILE_LOG_LEVEL"),
		LogMaxSizeMB:  parseInteger(getEnv("LOG_MAX_SIZE_MB")),
		LogMaxBackups: parseInteger(getEnv("LOG_MAX_BACKUPS")),
		LogMaxAgeDays: parseInteger(getEnv("LOG_MAX_AGE_DAYS")),

		// database
		DBHost:                getEnv("DB_HOST"),
		DBPort:                getEnv("DB_PORT"),
		DBUser:                getEnv("DB_USER"),
		DBPassword:            getEnv("DB_PASSWORD"),
		DBName:                getEnv("DB_NAME"),
		SSLMode:               getEnv("DB_SSL_MODE"),
		DBMaxConn:             int32(parseInteger(getEnvOrDefault("DB_MAX_CONN", "10"))),
		DBMinConn:             int32(parseInteger(getEnvOrDefault("DB_MIN_CONN", "1"))),
		DBConnMaxLifetime:     parseDuration(getEnvOrDefault("DB_CONN_MAX_LIFETIME", "0")),
		DBConnMaxIdleLifetime: parseDuration(getEnvOrDefault("DB_CONN_MAX_IDLE_LIFETIME", "0")),
		DBHealthCheckPeriod:   parseDuration(getEnvOrDefault("DB_HEALTH_CHECK_PERIOD", "0")),
		ConnectTimeout:        parseDuration(getEnvOrDefault("DB_CONNECT_TIMEOUT", "0")),

		// email
		EmailServiceBaseURL:        getEnv("EMAIL_SERVICE_BASE_URL"),
		EmailServiceToken:          getEnv("EMAIL_SERVICE"),
		EmailProvider:              getEnv("EMAIL_PROVIDER"),
		VerificationEmailExpiry:    getEnvOrDefault("VERIFICATION_EMAIL_EXPIRY", "15m"),
		VerificationEmailRateLimit: parseInteger(getEnvOrDefault("VERIFICATION_EMAIL_RATE_LIMIT", "5")),

		// Auth
		JwtSecret:          getEnv("JWT_SECRET"),
		AccessTokenExpiry:  parseDuration(getEnvOrDefault("ACCESS_TOKEN_EXPIRY", "15m")),
		RefreshTokenExpiry: parseDuration(getEnvOrDefault("REFRESH_TOKEN_EXPIRY", "168h")),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL"),
		GoogleOIDCIssuer:   getEnv("GOOGLE_OIDC_ISSUER"),
		FrontendOrigin:     getEnv("FRONTEND_ORIGIN"),

		// AWS credentials
		AWSRegion:          getEnv("AWS_REGION"),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY"),
		PresignedURLExpiry: parseInteger(getEnv("PRESIGNED_URL_EXPIRY")),
		UploadMaxFileSize:  parseInteger64(getEnv("UPLOAD_MAX_FILE_SIZE")),
		MaxActiveDevices:   parseInteger(getEnvOrDefault("MAX_ACTIVE_DEVICES", "3")),
	}

	appEnv := Environment(getEnv("APP_ENV"))
	if !appEnv.isValidEnv() {
		return nil, ErrInvalidAppEnv
	}
	cfg.AppEnv = appEnv
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
func getEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)
	val = strings.TrimSpace(val)
	if val == "" {
		return defaultValue
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

// parseDuration converts a string to a time.Duration
func parseDuration(val string) time.Duration {
	val = strings.TrimSpace(val)
	d, err := time.ParseDuration(val)
	if err != nil {
		panic(err)
	}
	return d
}
