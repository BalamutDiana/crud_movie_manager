package internal

import (
	"database/sql"
	"time"
)

type Movie struct {
	ID               int64
	Title            string
	Release          string
	StreamingService string
	SavedAt          time.Time
}

func GetMovies(db *sql.DB) ([]Movie, error) {
	rows, err := db.Query("select * from movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]Movie, 0)
	for rows.Next() {
		m := Movie{}
		err := rows.Scan(&m.ID, &m.Title, &m.Release, &m.StreamingService, &m.SavedAt)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func GetMovieByID(db *sql.DB, id int) (Movie, error) {
	var m Movie
	err := db.QueryRow("select * from movies where id = $1", id).
		Scan(&m.ID, &m.Title, &m.Release, &m.StreamingService, &m.SavedAt)

	return m, err
}

func InsertMovie(db *sql.DB, m Movie) error {
	_, err := db.Exec("insert into movies (title, release, streaming_service) values ($1, $2, $3)",
		m.Title, m.Release, m.StreamingService)

	return err
}

func DeleteMovie(db *sql.DB, id int) error {
	_, err := db.Exec("delete from movies where id = $1", id)
	return err
}

func UpdateMovie(db *sql.DB, id int, newMovie Movie) error {
	_, err := db.Exec("update movies set title=$1, release = $2, streaming_service = $3 where id = $4",
		newMovie.Title, newMovie.Release, newMovie.StreamingService, id)
	return err
}

// err = updateMovie(db, 2, Movie{
// 	Title:            "The Sandman",
// 	Release:          "2022",
// 	StreamingService: "Netflix",
// })
// if err != nil {
// 	log.Fatal(err)
// }

// err = insertMovie(db, Movie{
// 	Title:            "Wheel of time",
// 	Release:          "2021",
// 	StreamingService: "Amazon",
// })
// if err != nil {
// 	log.Fatal(err)
// }

// movies, err := internal.GetMovies(db)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println(movies)
