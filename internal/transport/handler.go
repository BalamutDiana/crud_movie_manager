package transport

import (
	"context"

	"net/http"

	_ "github.com/BalamutDiana/crud_movie_manager/docs"
	"github.com/BalamutDiana/crud_movie_manager/internal/domain"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	// JSONDocumentationPath is the path of the swagger documentation in json format.
	JSONDocumentationPath = "/documentation/json"
)

type Movies interface {
	Create(ctx context.Context, movie domain.Movie) error
	GetMovieByID(ctx context.Context, id int64) (domain.Movie, error)
	List(ctx context.Context) ([]domain.Movie, error)
	DeleteMovie(ctx context.Context, id int64) error
	UpdateMovie(ctx context.Context, id int64, newMovie domain.Movie) error
}

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	movieService Movies
	usersService User
}

type statusResponse struct {
	Message string `json:"status"`
}

func NewHandler(movies Movies, users User) *Handler {
	return &Handler{
		movieService: movies,
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
	}

	books := r.PathPrefix("/movies").Subrouter()
	{
		books.Use(h.authMiddleware)

		books.HandleFunc("", h.insertMovie).Methods(http.MethodPost)
		books.HandleFunc("", h.getMovies).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.getMovieByID).Methods(http.MethodGet)
		books.HandleFunc("/{id:[0-9]+}", h.deleteMovie).Methods(http.MethodDelete)
		books.HandleFunc("/{id:[0-9]+}", h.updateMovie).Methods(http.MethodPut)
	}

	return r
}
