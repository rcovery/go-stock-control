package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/rcovery/go-stock-control/internal/part"
	"github.com/rcovery/go-stock-control/internal/part/errs"
	part_service "github.com/rcovery/go-stock-control/internal/part/service"
)

type PartHandler struct {
	partService    *part_service.PartService
	restockService *part_service.RestockService
}

func NewPartHandler(partService *part_service.PartService, restockService *part_service.RestockService) *PartHandler {
	return &PartHandler{
		partService:    partService,
		restockService: restockService,
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
	var p part.Part
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request_body")
		return
	}

	created, err := h.partService.Create(r.Context(), p)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *PartHandler) listParts(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var (
		parts []part.Part
		err   error
	)
	if category != "" {
		parts, err = h.partService.ListByCategory(r.Context(), category)
	} else {
		parts, err = h.partService.List(r.Context())
	}
	if err != nil {
		h.writeError(w, err)
		return
	}

	if parts == nil {
		parts = []part.Part{}
	}

	writeJSON(w, http.StatusOK, parts)
}

func (h *PartHandler) updatePart(w http.ResponseWriter, r *http.Request) {
	id, err := part.ParseID(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_id")
		return
	}

	var p part.Part
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_request_body")
		return
	}
	p.ID = id

	updated, err := h.partService.Update(r.Context(), p)
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *PartHandler) deletePart(w http.ResponseWriter, r *http.Request) {
	id, err := part.ParseID(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid_id")
		return
	}

	if err := h.partService.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PartHandler) restockPriorities(w http.ResponseWriter, r *http.Request) {
	priorities, err := h.restockService.ListPriorities(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string][]part.RestockPriority{
		"priorities": priorities,
	})
}

func (h *PartHandler) writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, errs.NotFoundError) {
		writeJSONError(w, http.StatusNotFound, "not_found")
		return
	}

	var validationErr errs.ValidationError
	if errors.As(err, &validationErr) {
		writeJSONErrorDetailed(w, http.StatusBadRequest, "validation_error", validationErr.Message)
		return
	}

	writeJSONError(w, http.StatusInternalServerError, "internal_server_error")
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(body)
	if encodeErr != nil {
		log.Println("failed encoding json response:", encodeErr)
	}
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

func writeJSONErrorDetailed(w http.ResponseWriter, status int, code, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encodeErr := json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": detail,
	})
	if encodeErr != nil {
		log.Println("failed encoding json error:", encodeErr)
	}
}
