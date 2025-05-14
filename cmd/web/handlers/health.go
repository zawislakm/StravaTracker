package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("KURWA KURWA KURWA CHECK")
	dbHealth := h.db.Health()

	if dbHealth["status"] == "ok" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(h.db.Health())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(fmt.Printf("Error encoding health response: %v", err))
	}

}
