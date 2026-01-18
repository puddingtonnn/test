package app

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/puddingtonnnn/test/internal/config"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	http.ListenAndServe(":9090", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
