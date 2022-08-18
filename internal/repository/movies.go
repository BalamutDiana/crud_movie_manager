package internal

import (
	"context"
	"database/sql"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
)

type Movies struct {
	db *sql.DB
}

func NewMovies(db *sql.DB) *Movies {
	return &Movies{db}
}

func (m *Movies) GetMovies(ctx context.Context) ([]domain.Movie, error) {
	rows, err := m.db.Query("select * from movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make([]domain.Movie, 0)
	for rows.Next() {
		var m domain.Movie
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

func (m *Movies) GetMovieByID(ctx context.Context, id int64) (domain.Movie, error) {
	var movie domain.Movie

	err := m.db.QueryRow("select * from movies where id = $1", id).
		Scan(&movie.ID, &movie.Title, &movie.Release, &movie.StreamingService, &movie.SavedAt)

	return movie, err
}

func (m *Movies) InsertMovie(ctx context.Context, movie domain.Movie) error {
	_, err := m.db.Exec("insert into movies (title, release, streaming_service) values ($1, $2, $3)",
		movie.Title, movie.Release, movie.StreamingService)

	return err
}

func (m *Movies) DeleteMovie(ctx context.Context, id int64) error {
	_, err := m.db.Exec("delete from movies where id = $1", id)
	return err
}

func (m *Movies) UpdateMovie(ctx context.Context, id int64, newMovie domain.Movie) error {
	_, err := m.db.Exec("update movies set title=$1, release = $2, streaming_service = $3 where id = $4",
		newMovie.Title, newMovie.Release, newMovie.StreamingService, id)
	return err
}
