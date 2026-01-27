////go:build e2e

package api

import (
	"leti/pkg/service"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_FullAuthorFlow(t *testing.T) {
	repo := setupTestDBWithMigrations(t)
	srv := service.NewService(repo)

	r := newTestAPI(srv)

	ts := httptest.NewServer(r)
	defer ts.Close()

	_ = createAuthor(t, ts.URL, "Александр Дюма")

	authors := getAuthors(t, ts.URL)

	require.Len(t, authors, 1)

}
