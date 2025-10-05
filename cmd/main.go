package main

import (
	"belimang/internal/config"
	"belimang/internal/handlers"
	"belimang/internal/repository"
	"belimang/internal/route"
	"belimang/internal/services"
	"belimang/internal/utils"
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadAllAppConfig()
	if err != nil {
		 log.Fatal().Err(err)
	}

	dbp, err := config.InitDBConnection(cfg)
	if err != nil {
		 log.Fatal().Err(err)
	}
	defer func() {
		if dbp != nil {
			 dbp.Close()
		}
	}()

	minioClient, err := config.InitMCConncection(cfg)
	if err != nil {
		 log.Fatal().Err(err)
	}

	v := validator.New()
	utils.RegisterCustomValidations(v)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	repository.SetPool(dbp)

	authRepository := repository.NewAuthRepository(dbp)
	merchantRepository := repository.NewMerchantRepository(dbp)
	// purchaseRepository := repository.NewPurchaseRepository(dbp)

	hashingPool := services.NewHashingPool(2, 40)
	authService := services.NewAuthService(authRepository, hashingPool)
	fileService := services.NewFileService(minioClient, cfg)
	merchantService := services.NewMerchantService(merchantRepository)
	// purchaseService := services.NewPurchaseService(purchaseRepository)

	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(authService, v)
	merchantHandler := handlers.NewMerchantHandler(merchantService, v)
	// purchaseHandler := handlers.NewPurchaseHandler(purchaseService, v)

	r.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	route.RegisterAuthRoutes(r, authHandler)
	route.RegisterFileRoutes(r, fileHandler)
	route.RegisterMerchantRoutes(r, merchantHandler)

	server := http.Server{
		Addr:        cfg.Host + ":" + cfg.Port,
		Handler:     r,
		ReadTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Str("port", cfg.Port).Msg("server running")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("could not start server")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("server shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited")
}
