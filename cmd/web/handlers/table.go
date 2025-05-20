package handlers

import (
	"app/cmd/web/templates"
	"log"
	"net/http"
	"time"
)

func (h *Handler) HandleTable(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	yearFilter := r.URL.Query().Get("year")
	if yearFilter == "" {
		yearFilter = time.Now().Format("2006")
	}

	athletesData := h.db.GetAthletesData(yearFilter)
	component := templates.TableData(athletesData)
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HandleTable: %e", err)
	}
}
