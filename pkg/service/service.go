package service

import "leti/pkg/repository"

type Service struct {
	db repository.DataBase
}

func NewService(db repository.DataBase) *Service {
	return &Service{db: db}
}
