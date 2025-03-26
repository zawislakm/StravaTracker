package handlers

import "app/internal/database"

type Handler struct {
	db database.Service
}

func NewHandler(db database.Service) *Handler {
	return &Handler{
		db: db,
	}
}
