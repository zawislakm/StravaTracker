package handlers

import "app/internal/database"

type Handler struct {
	db *database.MongoDBClient
}

func NewHandler(db *database.MongoDBClient) *Handler {
	return &Handler{
		db: db,
	}
}
