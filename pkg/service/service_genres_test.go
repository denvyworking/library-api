package service

import (
	"context"
	"leti/pkg/models"
	"leti/pkg/repository/fake"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_NewGenre_EmptyName(t *testing.T) {
	fakeDB := &fake.FakeRepo{}
	svc := NewService(fakeDB)

	_, err := svc.NewAuthor(context.Background(), models.Author{Author: "   "})
	require.Error(t, err)
	require.Contains(t, err.Error(), "cannot be empty")
}
