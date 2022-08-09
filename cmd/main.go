package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Movie struct {
	ID               int64
	Title            string
	Release          string
	StreamingService string
	SavedAt          time.Time
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", "0.0.0.0", 5432, "postgres", "1marvin2mode3", "postgres")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select * from movies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	//fmt.Println(rows.Columns())

	movies := make([]Movie, 0)
	for rows.Next() {
		m := Movie{}
		err := rows.Scan(&m.ID, &m.Title, &m.Release, &m.StreamingService, &m.SavedAt)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(movies)
}
