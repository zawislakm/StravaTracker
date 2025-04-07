package handlers

import (
	"app/cmd/web/templates"
	"log"
	"net/http"
)

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	component := templates.Index()
	err = component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatalf("Error rendering in HandleIndex: %e", err)
	}
}
