package transport

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/BalamutDiana/crud_movie_manager/internal/domain"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// GetMovies godoc
// @Summary     Get movies
// @Description Get all movies list
// @Accept      json
// @Produce     json
// @Success     200 {object} []domain.Movie
// @Router      /movies [get]
func (h *Handler) getMovies(w http.ResponseWriter, r *http.Request) {
	m, err := h.movieService.List(r.Context())
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

	movie, err = h.movieService.GetMovieByID(r.Context(), id)
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

	err = h.movieService.Create(r.Context(), movie)

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

	if err = h.movieService.DeleteMovie(r.Context(), id); err != nil {
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

	err = h.movieService.UpdateMovie(r.Context(), id, upd)

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
