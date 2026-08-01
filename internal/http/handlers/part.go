package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rcovery/go-stock-control/internal/part"
)

type PartHandler struct {
	service *part.Service
}

func NewPartHandler(service *part.Service) *PartHandler {
	return &PartHandler{
		service: service,
	}
}

func (h *PartHandler) HandlePart() {
	http.HandleFunc("POST /parts", h.createPart)
	http.HandleFunc("GET /parts", h.listParts)
	http.HandleFunc("PUT /parts/{id}", h.updatePart)
	http.HandleFunc("DELETE /parts/{id}", h.deletePart)
	http.HandleFunc("GET /restock/priorities", h.restockPriorities)
}

func (h *PartHandler) createPart(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotImplemented, "not_implemented")
}

func (h *PartHandler) listParts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	_ = category
	writeJSONError(w, http.StatusNotImplemented, "not_implemented")
}

func (h *PartHandler) updatePart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_ = id
	writeJSONError(w, http.StatusNotImplemented, "not_implemented")
}

func (h *PartHandler) deletePart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	_ = id
	writeJSONError(w, http.StatusNotImplemented, "not_implemented")
}

func (h *PartHandler) restockPriorities(w http.ResponseWriter, r *http.Request) {
	writeJSONError(w, http.StatusNotImplemented, "not_implemented")
}

func writeJSONError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"error": code,
	})
	if encodeErr != nil {
		log.Println("failed encoding json error:", encodeErr)
	}
}
