package transport

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"

	"net/http"
	"strconv"

	_ "github.com/BalamutDiana/crud_movie_manager/docs"
	"github.com/BalamutDiana/crud_movie_manager/internal/domain"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	// JSONDocumentationPath is the path of the swagger documentation in json format.
	JSONDocumentationPath = "/documentation/json"
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

type statusResponse struct {
	Message string `json:"status"`
}

func NewHandler(movies Movies) *Handler {
	return &Handler{
		movieService: movies,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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

// GetMovies godoc
// @Summary     Get movies
// @Description Get all movies list
// @Accept      json
// @Produce     json
// @Success     200 {object} []domain.Movie
// @Router      /movies [get]
func (h *Handler) getMovies(w http.ResponseWriter, r *http.Request) {

	m, err := h.movieService.GetMovies(context.TODO())
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getMovies",
			"problem": "service problem",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(m)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getMovies",
			"problem": "marshaling error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(resp)
}

// GetMoviesByID godoc
// @Summary     Get movies by ID
// @Description Get movies by ID
// @Accept      json
// @Produce     json
// @Param       id      path     string true "account id"
// @Success     200     {object} domain.Movie
// @Router      /movies/{id} [get]
func (h *Handler) getMovieByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getMovieByID",
			"problem": "getIdFromRequest problem",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var movie interface{}

	movie, err = h.movieService.GetMovieByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, errors.New("Movie not found")) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.WithFields(log.Fields{
			"handler": "getMovieByID",
			"problem": "service problem",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(movie)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "getMovieByID",
			"problem": "marshaling error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// InsertMovie doc
// @Summary     Add new movie
// @Description Add new movie
// @Accept      json
// @Produce     json
// @Param       input body     domain.MovieMainInfo true "Add movie to list, 'id' and 'savedAt' not necessary params"
// @Success     200,201      {object}             domain.MovieMainInfo
// @Router      /movies [post]
func (h *Handler) insertMovie(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "insertMovie",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var movie domain.Movie

	if err = json.Unmarshal(reqBytes, &movie); err != nil {
		log.WithFields(log.Fields{
			"handler": "insertMovie",
			"problem": "unmarchaling request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.movieService.InsertMovie(context.TODO(), movie)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "insertMovie",
			"problem": "service problem",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteMovie doc
// @Summary     DeleteMovie from list
// @Description DeleteMovie from list
// @Accept      json
// @Produce     json
// @Param       id  body     id path string true "account id"
// @Success     200 {object} statusResponse
// @Router      /movies/{id} [delete]
func (h *Handler) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteMovie",
			"problem": "getIdFromRequest problem",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.movieService.DeleteMovie(context.TODO(), id); err != nil {
		log.WithFields(log.Fields{
			"handler": "deleteMovie",
			"problem": "service problem",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateMovie doc
// @Summary     Update movie info
// @Description Update movie info
// @Accept      json
// @Produce     json
// @Param       input body     id                   path string true "account id"
// @Param       input body     domain.MovieMainInfo true "Add movie to list, 'id' and 'savedAt' not necessary params"
// @Success     200   {object} domain.MovieMainInfo
// @Router      /movies/{id} [put]
func (h *Handler) updateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateMovie",
			"problem": "getIdFromRequest problem",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateMovie",
			"problem": "reading request body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upd domain.Movie
	if err = json.Unmarshal(reqBytes, &upd); err != nil {
		log.WithFields(log.Fields{
			"handler": "updateMovie",
			"problem": "unmarshaling error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.movieService.UpdateMovie(context.TODO(), id, upd)
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "updateMovie",
			"problem": "service problem",
		}).Error(err)
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
