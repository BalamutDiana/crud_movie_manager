package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BalamutDiana/crud_movie_manager/internal/config"
	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	rest "github.com/BalamutDiana/crud_movie_manager/internal/transport"
	"github.com/BalamutDiana/crud_movie_manager/pkg/database"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	booksRepo := repo.NewMovies(db)
	handler := rest.NewHandler(booksRepo)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
