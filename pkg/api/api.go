package api

import (
	"leti/pkg/service"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

// одна из самах популярных библиотек в go для роутинга
// gorillaMux

type api struct {
	r         *mux.Router
	srv       *service.Service
	logger    *slog.Logger
	authToken string
}

func New(router *mux.Router, srv *service.Service, logger *slog.Logger, at string) *api {
	return &api{r: router, srv: srv, logger: logger, authToken: at}
}

func (api *api) RegistreRoutes() {
	api.HandleBooks()
	api.HandleAuthors()
	api.HandleGenres()
}

func (api *api) HandleBooks() {
	// Публичные GET-запросы — без middleware
	api.r.HandleFunc("/api/books", api.books).Methods(http.MethodGet)
	api.r.HandleFunc("/api/books", api.books).Methods(http.MethodGet).Queries("id", "{id}")
	api.r.HandleFunc("/api/books/withauthors", api.booksWithAuthor).Methods(http.MethodGet)

	// Приватные операции — с middleware
	privateBooks := api.r.PathPrefix("/api/books").Subrouter()
	privateBooks.Use(api.middleware)
	privateBooks.HandleFunc("", api.books).Methods(http.MethodPost)
	privateBooks.HandleFunc("", api.books).Methods(http.MethodDelete).Queries("id", "{id}")
	privateBooks.HandleFunc("", api.books).Methods(http.MethodPatch).Queries("id", "{id}")
}

func (api *api) HandleAuthors() {
	api.r.HandleFunc("/api/authors", api.authors)
}

func (api *api) HandleGenres() {
	api.r.HandleFunc("/api/genres", api.genres)
}

func (api *api) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api.r)
}
