package main

import (
	"belimang/internal/config"
	"belimang/internal/repository"
	"belimang/internal/handlers"
	"belimang/internal/services"
	"belimang/internal/route"
	"context"
	"syscall"
	"os/signal"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadsAllAppConfig()
	if err != nil {
		 log.Fatal().Err(err)
	}

	dbp, err := config.InitsDBConnection(cfg)
	if err != nil {
		 log.Fatal().Err(err)
	}
	defer func() {
		if dbp != nil {
			 dbp.Close()
		}
	}()

	v := validator.New()
	
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	repository.SetPool(dbp)
	
	authRepository := repository.NewAuthRepository(dbp)
	// merchantRepository := repository.NewMerchantRepository(db)
	// purchaseRepository := repository.NewPurchaseRepository(db)

	authService := services.NewAuthService(authRepository)
	// fileService := services.NewFileService(cfg)
	// merchantService := services.NewMerchantService(merchantRepository)
	// purchaseService := services.NewPurchaseService(purchaseRepository)

	authHandler := handlers.NewAuthHandler(v, authService)
	// fileHandler := handlers.NewFileHandler(v, fileService)
	// merchantHandler := handlers.NewMerchantHandler(v, merchantService)
	// purchaseHandler := handlers.NewPurchaseHandler(v, purchaseService)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok"}`))
		})
		
		route.RegisterAuthRoutes(v1, authHandler)
		// route.RegisterFileRoutes(v1, fileHandler)
		// route.RegisterMerchantRoutes(v1, merchantHandler)
		// route.RegisterPurchaseRoutes(v1, purchaseHandler)
	})

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
