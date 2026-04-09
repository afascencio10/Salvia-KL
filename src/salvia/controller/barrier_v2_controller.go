package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type BarrierV2Controller struct {
	svc service.BarrierV2Service
}

func NewBarrierV2Controller(svc service.BarrierV2Service) *BarrierV2Controller {
	return &BarrierV2Controller{svc: svc}
}

func (c *BarrierV2Controller) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	b, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrBarrierV2NotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "barrera no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (c *BarrierV2Controller) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *BarrierV2Controller) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var b models.BarrierV2
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if b.CaseID == "" || b.FollowUpID == "" || b.Sector == "" || b.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos requeridos: case_id, follow_up_id, sector, description"})
		return
	}
	if err := c.svc.Create(r.Context(), &b); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (c *BarrierV2Controller) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var b models.BarrierV2
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	b.ID = id
	if err := c.svc.Update(r.Context(), &b); err != nil {
		if errors.Is(err, service.ErrBarrierV2NotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "barrera no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (c *BarrierV2Controller) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrBarrierV2NotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "barrera no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
