package dto

import "leti/pkg/models"

type CreateGenreRequest struct {
	Name string `json:"genre" validate:"required,min=1"`
}

func (cgr CreateGenreRequest) ToGenreModel() models.Genre {
	return models.Genre{
		Genre: cgr.Name,
	}
}

// GenreResponse — ответ с информацией о жанре
type GenreResponse struct {
	ID   int    `json:"id"`
	Name string `json:"genre"`
}

func GenreFromModels(genre models.Genre) GenreResponse {
	return GenreResponse{
		ID:   genre.ID,
		Name: genre.Genre,
	}
}

func FromGenreModelsArray(genres []models.Genre) []GenreResponse {
	resp := make([]GenreResponse, len(genres))
	for i, genre := range genres {
		resp[i] = GenreFromModels(genre)
	}
	return resp
}
