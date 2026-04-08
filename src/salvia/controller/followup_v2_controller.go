// Package controller expone los endpoints HTTP de la capa salvia (Fase 2).
package controller

import (
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// FollowUpV2Controller maneja las peticiones HTTP para FollowUpV2.
type FollowUpV2Controller struct {
	svc service.FollowUpV2Service
}

// NewFollowUpV2Controller construye el controlador inyectando el servicio.
func NewFollowUpV2Controller(svc service.FollowUpV2Service) *FollowUpV2Controller {
	return &FollowUpV2Controller{svc: svc}
}

// GetByIDHandler godoc
// GET /followups/v2?id={uuid}
// Respuestas: 200 OK | 400 Bad Request | 404 Not Found | 500 Internal Server Error
func (c *FollowUpV2Controller) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}

	fu, err := c.svc.GetFollowUpByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFollowUpNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "registro no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno del servidor"})
		return
	}

	writeJSON(w, http.StatusOK, fu)
}

// ListHandler godoc
// GET /followups/v2/list?page={int}&limit={int}
// page base-0, limit por defecto 20.
// Respuestas: 200 OK | 500 Internal Server Error
func (c *FollowUpV2Controller) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)

	result, err := c.svc.GetPaginatedFollowUps(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno del servidor"})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// --- helpers privados ---

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func queryInt(r *http.Request, key string, defaultVal int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return defaultVal
	}
	return v
}
