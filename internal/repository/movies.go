package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
	cc "github.com/BalamutDiana/custom_cache"
)

type Movies struct {
	db    *sql.DB
	cache *cc.Cache
}

func NewMovies(database *sql.DB, cache *cc.Cache) *Movies {
	return &Movies{
		db:    database,
		cache: cache,
	}
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

	for _, item := range movies {
		if _, err := m.cache.Get(fmt.Sprint(item.ID)); err != nil {
			m.cache.Set(fmt.Sprint(item.ID), item, time.Minute*2)
		}
	}

	return movies, nil
}

func (m *Movies) GetMovieByID(ctx context.Context, id int64) (domain.Movie, error) {
	if movie, err := m.cache.Get(fmt.Sprint(id)); err == nil {
		return movie.(domain.Movie), err
	}

	var movie domain.Movie

	err := m.db.QueryRow("select * from movies where id = $1", id).
		Scan(&movie.ID, &movie.Title, &movie.Release, &movie.StreamingService, &movie.SavedAt)

	return movie, err
}

func (m *Movies) InsertMovie(ctx context.Context, movie domain.Movie) error {
	if _, err := m.db.Exec("insert into movies (title, release, streaming_service) values ($1, $2, $3)",
		movie.Title, movie.Release, movie.StreamingService); err != nil {
		return err
	}

	m.cache.Set(fmt.Sprint(movie.ID), movie, time.Minute*2)
	return nil
}

func (m *Movies) DeleteMovie(ctx context.Context, id int64) error {
	if _, err := m.db.Exec("delete from movies where id = $1", id); err != nil {
		return err
	}
	if err := m.cache.Delete(fmt.Sprint(id)); err != nil {
		return err
	}
	return nil
}

func (m *Movies) UpdateMovie(ctx context.Context, id int64, newMovie domain.Movie) error {
	if _, err := m.db.Exec("update movies set title=$1, release = $2, streaming_service = $3 where id = $4",
		newMovie.Title, newMovie.Release, newMovie.StreamingService, id); err != nil {
		return err
	}

	m.cache.Set(fmt.Sprint(id), newMovie, time.Minute*2)
	return nil
}
