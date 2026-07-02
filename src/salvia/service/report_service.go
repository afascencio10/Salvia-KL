package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/xuri/excelize/v2"
)

type FollowUpsReportRange struct {
	StartDate time.Time
	EndDate   time.Time
}

type ReportService interface {
	GenerateConsolidatedFollowUpsReport(ctx context.Context, r FollowUpsReportRange) (*excelize.File, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GenerateConsolidatedFollowUpsReport(ctx context.Context, r FollowUpsReportRange) (*excelize.File, error) {
	// 1. Consultar casos creados en el rango
	cases, err := s.repo.FindCasesInDateRange(ctx, r.StartDate, r.EndDate)
	if err != nil {
		return nil, err
	}

	var caseICodes []string
	for _, c := range cases {
		caseICodes = append(caseICodes, c.VictimCaseICode)
	}

	// 2. Si no hay casos creados en el rango, no generamos el Excel; retornamos error de negocio
	if len(caseICodes) == 0 {
		return nil, errors.New("no_cases_found")
	}

	// 3. Consultar los seguimientos asociados en lote
	followUps, err := s.repo.FindFollowUpsByCaseICodes(ctx, caseICodes)
	if err != nil {
		return nil, err
	}

	// 4. Consultar las respuestas del formulario asociadas
	var submissionIDs []string
	for _, fu := range followUps {
		if fu.FormSubmissionID != nil && *fu.FormSubmissionID != "" {
			submissionIDs = append(submissionIDs, *fu.FormSubmissionID)
		}
	}

	var answers []repository.AnswerDTO
	if len(submissionIDs) > 0 {
		answers, err = s.repo.FindAnswersByFormSubmissions(ctx, submissionIDs)
		if err != nil {
			return nil, err
		}
	}

	// 5. Consultar los eventos del timeline en lote
	events, err := s.repo.FindTimelineEventsByCaseICodes(ctx, caseICodes)
	if err != nil {
		return nil, err
	}

	// 6. Construir el Excel consolidado con los datos enriquecidos
	return BuildFollowUpsExcel(cases, followUps, answers, events)
}
