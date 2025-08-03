package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/wrestler094/launchpad/internal/api"
	"github.com/wrestler094/launchpad/internal/contracts"
	"github.com/wrestler094/launchpad/internal/services"
	"github.com/wrestler094/launchpad/internal/storage"
)

func main() {
	// Initialize database
	db, err := storage.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := storage.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize blockchain client
	client, err := contracts.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService()
	tokenService := services.NewTokenService(client, db)
	presaleService := services.NewPresaleService(client, db)

	// Initialize API handlers
	apiHandlers := api.NewHandlers(authService, tokenService, presaleService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Get("/nonce", apiHandlers.GenerateNonce)
			r.Post("/login", apiHandlers.Login)
			r.Post("/verify", apiHandlers.VerifyToken)
		})

		// Protected routes
		r.Route("/", func(r chi.Router) {
			r.Use(apiHandlers.AuthMiddleware)

			// Token routes
			r.Route("/token", func(r chi.Router) {
				r.Post("/create", apiHandlers.CreateToken)
				r.Get("/list", apiHandlers.ListTokens)
				r.Get("/{address}", apiHandlers.GetToken)
			})

			// Presale routes
			r.Route("/presale", func(r chi.Router) {
				r.Post("/create", apiHandlers.CreatePresale)
				r.Get("/list", apiHandlers.ListPresales)
				r.Get("/{id}", apiHandlers.GetPresale)
				r.Post("/{id}/participate", apiHandlers.ParticipateInPresale)
			})
		})

		// Public presale routes (for landing pages)
		r.Route("/public/presale", func(r chi.Router) {
			r.Get("/{id}", apiHandlers.GetPublicPresale)
		})
	})

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}