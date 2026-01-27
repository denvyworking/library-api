////go:build e2e

package api

import (
	"leti/pkg/service"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_FullGenreFlow(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	_ = createGenre(t, ts.URL, "Рассказ")

	genres := getGenres(t, ts.URL)

	require.Len(t, genres, 1)

}
