//go:build e2e

package api

import (
	"leti/pkg/service"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestE2E_FullAuthorFlow(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	authToken := "adminToken"

	r := mux.NewRouter()
	apiInst := New(r, srv, logger, authToken)
	apiInst.HandleAuthors()

	ts := httptest.NewServer(r)
	defer ts.Close()

	_ = createAuthor(t, ts.URL, "Александр Дюма")

	authors := getAuthors(t, ts.URL)

	require.Len(t, authors, 1)

}
