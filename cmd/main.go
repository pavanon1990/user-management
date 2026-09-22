package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/adapter/handler"
	"user-management-api/internal/adapter/job"
	"user-management-api/internal/adapter/repository"
	"user-management-api/internal/adapter/router"
	"user-management-api/internal/core/service"
	"user-management-api/internal/infra/mongodb"
	"user-management-api/internal/infra/token"
	"user-management-api/pkg/config"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	appCfg := config.LoadAppConfig()
	mgCfg := config.LoadMongoConfig()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mongoClient, err := mongodb.NewClient(ctx, mgCfg)
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(mgCfg.Database)

	tokenProvider := token.NewJWTProvider(appCfg.SecretKey, appCfg.TokenTTL)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, tokenProvider)
	userHandler := handler.NewUserHandler(userService)

	r := router.NewRouter(userHandler, tokenProvider)

	srv := &http.Server{
		Addr:    ":" + appCfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on :%s", appCfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	go job.StartUserCountLogger(ctx, userService, 10*time.Second)

	<-ctx.Done()
	log.Println("shutdown signal received, shutting down gracefully...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		log.Printf("mongo disconnect error: %v", err)
	}

	log.Println("server stopped")
}
