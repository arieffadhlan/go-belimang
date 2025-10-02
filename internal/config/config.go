package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Host      string
	Port      string
	Envs      string
	JWTSecret string

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	AccessKey  string
	SecretKey  string
	UseSSL     bool
	Endpoint   string
	BucketName string
}

func LoadAllAppConfig() (Config, error) {
	_ = godotenv.Load()
	useSSL, _ := strconv.ParseBool(os.Getenv("MINIO_SSL"))

	return Config{
		Port:      os.Getenv("APP_PORT"),
		Envs:      os.Getenv("APP_ENVS"),
		JWTSecret: os.Getenv("JWT_SECRET"),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBName: os.Getenv("DB_NAME"),

		AccessKey:  os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:  os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:     useSSL,
		Endpoint:   os.Getenv("MINIO_ENDPOINT"),
		BucketName: os.Getenv("MINIO_BUCKET_NAME"),
	}, nil
}

func InitDBConnection(cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		 return nil, fmt.Errorf("failed to parse pgconfig: %w", err)
	}

	poolCfg.MaxConns = 50
	poolCfg.MinConns = 10
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 10 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		 return nil, fmt.Errorf("failed to create pgpool: %w", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		 return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbpool, nil
}

func InitMCConncection(cfg Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		 return nil, err
	}

	return minioClient, nil
}
