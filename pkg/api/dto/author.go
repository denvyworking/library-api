package dto

import "leti/pkg/models"

type CreateAuthorRequest struct {
	Name string `json:"author" validate:"required,min=1"`
}

func (crq CreateAuthorRequest) ToAuthorModel() models.Author {
	return models.Author{
		Author: crq.Name,
	}
}

// AuthorResponse — ответ с информацией об авторе
type AuthorResponse struct {
	ID   int    `json:"id"`
	Name string `json:"author"`
}

func FromAuthorModels(author models.Author) AuthorResponse {
	return AuthorResponse{
		ID:   author.ID,
		Name: author.Author,
	}
}

func FromAuthorModelsArray(authors []models.Author) []AuthorResponse {
	resp := make([]AuthorResponse, len(authors))
	for i, author := range authors {
		resp[i] = FromAuthorModels(author)
	}
	return resp
}
