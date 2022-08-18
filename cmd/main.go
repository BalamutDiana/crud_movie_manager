package main

import (
	"log"
	"net/http"
	"os"
	"time"

	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	rest "github.com/BalamutDiana/crud_movie_manager/internal/transport"
	"github.com/BalamutDiana/crud_movie_manager/pkg/database"

	_ "github.com/lib/pq"
)

func main() {

	password := os.Getenv("DB_PASSWORD")
	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "host.docker.internal",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	booksRepo := repo.NewMovies(db)
	handler := rest.NewHandler(booksRepo)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
