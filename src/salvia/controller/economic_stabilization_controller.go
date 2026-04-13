package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type EconomicStabilizationController struct {
	svc service.EconomicStabilizationService
}

func NewEconomicStabilizationController(svc service.EconomicStabilizationService) *EconomicStabilizationController {
	return &EconomicStabilizationController{svc: svc}
}

func (c *EconomicStabilizationController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	es, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrEconomicStabilizationNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "estabilización no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, es)
}

func (c *EconomicStabilizationController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *EconomicStabilizationController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var es models.EconomicStabilization
	if err := json.NewDecoder(r.Body).Decode(&es); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if es.CaseID == "" || es.FollowUpID == "" || es.Type == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos requeridos: case_id, follow_up_id, type"})
		return
	}
	if err := c.svc.Create(r.Context(), &es); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, es)
}

func (c *EconomicStabilizationController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var es models.EconomicStabilization
	if err := json.NewDecoder(r.Body).Decode(&es); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	es.ID = id
	if err := c.svc.Update(r.Context(), &es); err != nil {
		if errors.Is(err, service.ErrEconomicStabilizationNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "estabilización no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, es)
}

func (c *EconomicStabilizationController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrEconomicStabilizationNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "estabilización no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
