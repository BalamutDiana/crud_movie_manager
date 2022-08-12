package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
	"github.com/gorilla/mux"
)

type Movies interface {
	InsertMovie(ctx context.Context, movie domain.Movie) error
	GetMovieByID(ctx context.Context, id int64) (domain.Movie, error)
	GetMovies(ctx context.Context) ([]domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
	UpdateMovie(ctx context.Context, id int64, newMovie domain.Movie) error
}

type Handler struct {
	movieService Movies
}

func NewHandler(movies Movies) *Handler {
	return &Handler{
		movieService: movies,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	books := r.PathPrefix("/movies").Subrouter()
	{
		books.HandleFunc("", h.insertMovie).Methods(http.MethodPost)
		books.HandleFunc("", h.getMovies).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getMovieByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteMovie).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateMovie).Methods(http.MethodPut)
	}

	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) getMovies(w http.ResponseWriter, r *http.Request) {

	m, err := h.movieService.GetMovies(context.TODO())

	if err != nil {
		log.Println("getMovies() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(m)
	if err != nil {
		log.Println("getMovies() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handler) getMovieByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)

	if err != nil {
		log.Println("getMovieByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.GetMovieByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, errors.New("Movie not found")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("getMoviesByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(movie)
	if err != nil {
		log.Println("getMoviesByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) insertMovie(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("insertMovie() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var movie domain.Movie

	if err = json.Unmarshal(reqBytes, &movie); err != nil {
		log.Println("unmarshaling error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.movieService.InsertMovie(context.TODO(), movie)
	if err != nil {
		log.Println("createMovie() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteMovie() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.movieService.DeleteMovie(context.TODO(), id)
	if err != nil {
		log.Println("deleteMovie() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) updateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("updateMovie() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("updateMovie() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upd domain.Movie
	if err = json.Unmarshal(reqBytes, &upd); err != nil {
		log.Println("unmarshaling error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.movieService.UpdateMovie(context.TODO(), id, upd)
	if err != nil {
		log.Println("UpdateMovie() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
