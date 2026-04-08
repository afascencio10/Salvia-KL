package controller

import (
	"bitsflow/internal/models"
	"bitsflow/salvia/service"
	"encoding/json"
	"errors"
	"net/http"
)

type QuestionController struct {
	svc service.QuestionService
}

func NewQuestionController(svc service.QuestionService) *QuestionController {
	return &QuestionController{svc: svc}
}

func (c *QuestionController) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	q, err := c.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "pregunta no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, q)
}

func (c *QuestionController) ListHandler(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 0)
	limit := queryInt(r, "limit", 20)
	result, err := c.svc.List(r.Context(), page, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (c *QuestionController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var q models.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	if q.FormID == "" || q.FormSectionID == "" || q.QuestionTypeID == "" || q.Description == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "campos requeridos: form_id, form_section_id, question_type_id, description"})
		return
	}
	if err := c.svc.Create(r.Context(), &q); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusCreated, q)
}

func (c *QuestionController) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	var q models.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "body inválido"})
		return
	}
	q.ID = id
	if err := c.svc.Update(r.Context(), &q); err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "pregunta no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	writeJSON(w, http.StatusOK, q)
}

func (c *QuestionController) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "parámetro 'id' requerido"})
		return
	}
	if err := c.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "pregunta no encontrada"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error interno"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
