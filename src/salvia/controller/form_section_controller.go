package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type FormSectionController struct {
	svc service.FormSectionService
}

func NewFormSectionController(svc service.FormSectionService) *FormSectionController {
	return &FormSectionController{svc: svc}
}

func (c *FormSectionController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	fs, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "sección no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, fs)
}

func (c *FormSectionController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *FormSectionController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var fs models.FormSection
	if err := json.NewDecoder(r.Body).Decode(&fs); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if fs.Name == "" || fs.FormID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos 'name' y 'form_id' requeridos"})
		return
	}
	if err := c.svc.Create(r.Context(), &fs); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, fs)
}

func (c *FormSectionController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var fs models.FormSection
	if err := json.NewDecoder(r.Body).Decode(&fs); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	fs.ID = id
	if err := c.svc.Update(r.Context(), &fs); err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "sección no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, fs)
}

func (c *FormSectionController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrFormSectionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "sección no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
