package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BalamutDiana/crud_movie_manager/internal/config"
	repo "github.com/BalamutDiana/crud_movie_manager/internal/repository"
	"github.com/BalamutDiana/crud_movie_manager/internal/service"
	rest "github.com/BalamutDiana/crud_movie_manager/internal/transport"
	"github.com/BalamutDiana/crud_movie_manager/pkg/database"
	"github.com/BalamutDiana/crud_movie_manager/pkg/hash"
	"github.com/BalamutDiana/custom_cache"
	"github.com/sirupsen/logrus"

	_ "github.com/BalamutDiana/crud_movie_manager/docs"
	_ "github.com/lib/pq"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// @title       CRUD movie manager API
// @version     1.0
// @description API server for saving movies

// @host     localhost:8080
// @BasePath /

func main() {
	cfg, err := config.New("configs", "config")
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

	hasher := hash.NewSHA1Hasher("salt")

	cache := custom_cache.New()
	booksRepo := repo.NewMovies(db, cache)
	usersRepo := repo.NewUsers(db)
	usersService := service.NewUsers(usersRepo, hasher, []byte("dFscwEtgdS"), cfg.Auth.TokenTTL)

	handler := rest.NewHandler(booksRepo, usersService)

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
