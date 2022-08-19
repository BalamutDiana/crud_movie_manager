package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BalamutDiana/crud_movie_manager/internal/config"
	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	rest "github.com/BalamutDiana/crud_movie_manager/internal/transport"
	"github.com/BalamutDiana/crud_movie_manager/pkg/database"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"method":  "config.New",
			"problem": "creating config",
		}).Fatal(err)
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
		logrus.WithFields(logrus.Fields{
			"method":  "database.NewPostgresConnection",
			"problem": "creating connection",
		}).Fatal(err)
	}
	defer db.Close()

	booksRepo := repo.NewMovies(db)
	handler := rest.NewHandler(booksRepo)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	logrus.Info("SERVER STARTED")

	if err := srv.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"method":  "srv.ListenAndServe",
			"problem": "http server problem",
		}).Fatal(err)
	}
}
