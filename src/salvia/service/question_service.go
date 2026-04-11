package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrQuestionNotFound = errors.New("question: registro no encontrado")

type QuestionService interface {
	GetByID(ctx context.Context, id string) (*models.Question, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.Question], error)
	ListByFormID(ctx context.Context, formID string) ([]models.Question, error)
	Create(ctx context.Context, q *models.Question) error
	Update(ctx context.Context, q *models.Question) error
	Delete(ctx context.Context, id string) error
}

type questionService struct {
	repo repository.QuestionRepository
}

func NewQuestionService(repo repository.QuestionRepository) QuestionService {
	return &questionService{repo: repo}
}

func (s *questionService) GetByID(ctx context.Context, id string) (*models.Question, error) {
	q, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}
	return q, nil
}

func (s *questionService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Question], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *questionService) ListByFormID(ctx context.Context, formID string) ([]models.Question, error) {
	return s.repo.FindByFormID(ctx, formID)
}

func (s *questionService) Create(ctx context.Context, q *models.Question) error {
	return s.repo.Create(ctx, q)
}

func (s *questionService) Update(ctx context.Context, q *models.Question) error {
	return s.repo.Update(ctx, q)
}

func (s *questionService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrQuestionNotFound
	}
	return err
}
