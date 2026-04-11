package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type PsychosocialSupportController struct {
	svc service.PsychosocialSupportService
}

func NewPsychosocialSupportController(svc service.PsychosocialSupportService) *PsychosocialSupportController {
	return &PsychosocialSupportController{svc: svc}
}

func (c *PsychosocialSupportController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	ps, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrPsychosocialSupportNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "apoyo no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, ps)
}

func (c *PsychosocialSupportController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *PsychosocialSupportController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var ps models.PsychosocialSupport
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if ps.CaseID == "" || ps.FollowUpID == "" || ps.Type == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos requeridos: case_id, follow_up_id, type"})
		return
	}
	if err := c.svc.Create(r.Context(), &ps); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, ps)
}

func (c *PsychosocialSupportController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var ps models.PsychosocialSupport
	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	ps.ID = id
	if err := c.svc.Update(r.Context(), &ps); err != nil {
		if errors.Is(err, service.ErrPsychosocialSupportNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "apoyo no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, ps)
}

func (c *PsychosocialSupportController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrPsychosocialSupportNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "apoyo no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
