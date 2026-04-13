package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type EmergencyMeasureController struct {
	svc service.EmergencyMeasureService
}

func NewEmergencyMeasureController(svc service.EmergencyMeasureService) *EmergencyMeasureController {
	return &EmergencyMeasureController{svc: svc}
}

func (c *EmergencyMeasureController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	em, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEmergencyMeasureNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "medida no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, em)
}

func (c *EmergencyMeasureController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *EmergencyMeasureController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var em models.EmergencyMeasure
	if err := json.NewDecoder(r.Body).Decode(&em); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if em.CaseID == "" || em.FollowUpID == "" || em.Type == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos requeridos: case_id, follow_up_id, type"})
		return
	}
	if err := c.svc.Create(r.Context(), &em); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, em)
}

func (c *EmergencyMeasureController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var em models.EmergencyMeasure
	if err := json.NewDecoder(r.Body).Decode(&em); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	em.ID = id
	if err := c.svc.Update(r.Context(), &em); err != nil {
		if errors.Is(err, service.ErrEmergencyMeasureNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "medida no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, em)
}

func (c *EmergencyMeasureController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrEmergencyMeasureNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "medida no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
