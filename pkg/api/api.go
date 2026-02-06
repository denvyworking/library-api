package api

import (
	"leti/pkg/auth"
	"leti/pkg/service"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// одна из самых популярных библиотек в go для роутинга
// gorillaMux

type api struct {
	r          *mux.Router
	srv        *service.Service
	logger     *slog.Logger
	jwtService *auth.JWTService
}

func New(router *mux.Router, srv *service.Service, logger *slog.Logger, at *auth.JWTService) *api {
	return &api{r: router, srv: srv, logger: logger, jwtService: at}
}

func (api *api) RegistreRoutes() {
	api.HandleAuth()
	api.HandleBooks()
	api.HandleAuthors()
	api.HandleGenres()
}

func (api *api) HandleBooks() {
	api.r.HandleFunc("/api/v1/books", api.getBooks).Methods(http.MethodGet)
	api.r.HandleFunc("/api/v1/book", api.getBookById).Methods(http.MethodGet).Queries("id", "{id}")
	api.r.HandleFunc("/api/v1/books/withauthors", api.booksWithAuthor).Methods(http.MethodGet)

	privateBooksV1 := api.r.PathPrefix("/api/v1/books").Subrouter()
	privateBooksV1.Use(api.middleware)
	privateBooksV1.HandleFunc("", api.createBook).Methods(http.MethodPost)
	privateBooksV1.HandleFunc("", api.deleteBook).Methods(http.MethodDelete).Queries("id", "{id}")
	privateBooksV1.HandleFunc("", api.updateBook).Methods(http.MethodPatch).Queries("id", "{id}")

	api.r.HandleFunc("/api/books", api.getBooks).Methods(http.MethodGet)
	api.r.HandleFunc("/api/book", api.getBookById).Methods(http.MethodGet).Queries("id", "{id}")
	api.r.HandleFunc("/api/books/withauthors", api.booksWithAuthor).Methods(http.MethodGet)

	privateBooks := api.r.PathPrefix("/api/books").Subrouter()
	privateBooks.Use(api.middleware)
	privateBooks.HandleFunc("", api.createBook).Methods(http.MethodPost)
	privateBooks.HandleFunc("", api.deleteBook).Methods(http.MethodDelete).Queries("id", "{id}")
	privateBooks.HandleFunc("", api.updateBook).Methods(http.MethodPatch).Queries("id", "{id}")
}

func (api *api) HandleAuthors() {
	api.r.HandleFunc("/api/v1/authors", api.getAuthors).Methods(http.MethodGet)
	privateAuthorsV1 := api.r.PathPrefix("/api/v1/authors").Subrouter()
	privateAuthorsV1.Use(api.middleware)
	privateAuthorsV1.HandleFunc("", api.postAuthors).Methods(http.MethodPost)

	api.r.HandleFunc("/api/authors", api.getAuthors).Methods(http.MethodGet)
	api.r.HandleFunc("/api/authors", api.postAuthors).Methods(http.MethodPost)
}

func (api *api) HandleGenres() {
	api.r.HandleFunc("/api/v1/genres", api.getGenres).Methods(http.MethodGet)
	privateGenresV1 := api.r.PathPrefix("/api/v1/genres").Subrouter()
	privateGenresV1.Use(api.middleware)
	privateGenresV1.HandleFunc("", api.postGenres).Methods(http.MethodPost)

	api.r.HandleFunc("/api/genres", api.getGenres).Methods(http.MethodGet)
	api.r.HandleFunc("/api/genres", api.postGenres).Methods(http.MethodPost)
}

func (api *api) HandleAuth() {
	api.r.HandleFunc("/api/v1/auth/login", api.login).Methods(http.MethodPost)
	api.r.HandleFunc("/api/v1/auth/register", api.register).Methods(http.MethodPost)
	api.r.HandleFunc("/api/v1/auth/refresh", api.refresh).Methods(http.MethodPost)
	api.r.HandleFunc("/api/auth/login", api.login).Methods(http.MethodPost)
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
