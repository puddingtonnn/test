package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/puddingtonnn/test/internal/config"
	"github.com/puddingtonnn/test/internal/handler"
	"github.com/puddingtonnn/test/internal/repository"
	"github.com/puddingtonnn/test/internal/service"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	logger.Info("starting application", "port", cfg.AppPort, "env", "production")

	sqlDB, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer sqlDB.Close()

	goose.SetDialect("postgres")
	goose.SetBaseFS(embedMigrations)
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to init gorm: %v", err)
	}

	repo := repository.NewChatRepository(gormDB)
	svc := service.NewChatService(repo)
	h := handler.NewHandler(svc)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /chats/", h.CreateChat)
	mux.HandleFunc("POST /chats/{id}/messages/", h.CreateMessage)
	mux.HandleFunc("GET /chats/{id}", h.GetChat)
	mux.HandleFunc("DELETE /chats/{id}", h.DeleteChat)

	addr := fmt.Sprintf(":%s", cfg.AppPort)

	handlerWithLogging := handler.LoggingMiddleware(mux)

	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handlerWithLogging); err != nil { // <-- Сюда
		log.Fatal(err)
	}
}
