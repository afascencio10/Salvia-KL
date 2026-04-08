package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

// FormController maneja los endpoints HTTP para Form.
type FormController struct {
	svc service.FormService
}

func NewFormController(svc service.FormService) *FormController {
	return &FormController{svc: svc}
}

// GET /forms?id={uuid}
func (c *FormController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	f, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// GET /forms/list?page={int}&limit={int}
func (c *FormController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// POST /forms
func (c *FormController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var f models.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if f.Name == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campo 'name' requerido"})
		return
	}
	if err := c.svc.Create(r.Context(), &f); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, f)
}

// PUT /forms?id={uuid}
func (c *FormController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var f models.Form
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	f.ID = id
	if err := c.svc.Update(r.Context(), &f); err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// DELETE /forms?id={uuid}
func (c *FormController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrFormNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "formulario no encontrado"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
