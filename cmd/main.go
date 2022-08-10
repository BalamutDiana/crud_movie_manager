package main

import (
	"fmt"
	"log"

	movies "github.com/BalamutDiana/crud_movie_manager/internal"
	db "github.com/BalamutDiana/crud_movie_manager/pkg/database"
	_ "github.com/lib/pq"
)

func main() {
	db, err := db.NewPostgresConnection(db.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "1marvin2mode3",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	m, err := movies.GetMovieByID(db, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)
}
