package handlers

import (
	"app/cmd/web/templates"
	"app/internal/model"
	"log"
	"net/http"
)

func (h *Handler) HandleTable(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	yearFilter := r.URL.Query().Get("year")
	if yearFilter == "" {
		yearFilter = "2025" // TODO make it current year
	}

	sortField := r.URL.Query().Get("sortField")
	if sortField == "" {
		sortField = "Distance"
	}

	athletesData := h.db.GetAthletesData(yearFilter)
	model.SortAthletesData(athletesData, sortField)
	component := templates.Table(athletesData)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HandleTable: %e", err)
	}
}
