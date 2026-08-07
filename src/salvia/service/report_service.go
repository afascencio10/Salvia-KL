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
	GenerateConsolidatedContactsReport(ctx context.Context, r FollowUpsReportRange) (*excelize.File, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GenerateConsolidatedFollowUpsReport(ctx context.Context, r FollowUpsReportRange) (*excelize.File, error) {
	// 1. Consultar casos creados en el rango — por chunks de 2 meses para evitar timeout del pooler
	var cases []repository.CaseReportDTO
	chunkStart := r.StartDate
	for chunkStart.Before(r.EndDate) {
		chunkEnd := chunkStart.AddDate(0, 2, 0)
		if chunkEnd.After(r.EndDate) {
			chunkEnd = r.EndDate
		}
		chunk, err := s.repo.FindCasesInDateRange(ctx, chunkStart, chunkEnd)
		if err != nil {
			return nil, err
		}
		cases = append(cases, chunk...)
		chunkStart = chunkEnd.Add(time.Second)
	}

	var caseICodes []string
	for _, c := range cases {
		caseICodes = append(caseICodes, c.VictimCaseICode)
	}

	// 2. Si no hay casos creados en el rango, no generamos el Excel; retornamos error de negocio
	if len(caseICodes) == 0 {
		return nil, errors.New("no_cases_found")
	}

	// 3. Consultar los seguimientos asociados — por chunks de 500 caseICodes
	var followUps []repository.FollowUpReportDTO
	for i := 0; i < len(caseICodes); i += 500 {
		end := i + 500
		if end > len(caseICodes) {
			end = len(caseICodes)
		}
		chunk, err := s.repo.FindFollowUpsByCaseICodes(ctx, caseICodes[i:end])
		if err != nil {
			return nil, err
		}
		followUps = append(followUps, chunk...)
	}

	// 4. Consultar las respuestas del formulario asociadas — por chunks de 500
	var submissionIDs []string
	for _, fu := range followUps {
		if fu.FormSubmissionID != nil && *fu.FormSubmissionID != "" {
			submissionIDs = append(submissionIDs, *fu.FormSubmissionID)
		}
	}

	var answers []repository.AnswerDTO
	for i := 0; i < len(submissionIDs); i += 500 {
		end := i + 500
		if end > len(submissionIDs) {
			end = len(submissionIDs)
		}
		chunk, err := s.repo.FindAnswersByFormSubmissions(ctx, submissionIDs[i:end])
		if err != nil {
			return nil, err
		}
		answers = append(answers, chunk...)
	}

	// 5. Consultar los eventos del timeline — por chunks de 500
	var events []repository.TimelineReportDTO
	for i := 0; i < len(caseICodes); i += 500 {
		end := i + 500
		if end > len(caseICodes) {
			end = len(caseICodes)
		}
		chunk, err := s.repo.FindTimelineEventsByCaseICodes(ctx, caseICodes[i:end])
		if err != nil {
			return nil, err
		}
		events = append(events, chunk...)
	}

	// 6. Construir el Excel consolidado con los datos enriquecidos
	return BuildFollowUpsExcel(cases, followUps, answers, events)
}

func (s *reportService) GenerateConsolidatedContactsReport(ctx context.Context, r FollowUpsReportRange) (*excelize.File, error) {
	contacts, err := s.repo.FindContactsInDateRange(ctx, r.StartDate, r.EndDate)
	if err != nil {
		return nil, err
	}

	// Si no hay reportes creados en el rango, no se genera archivo (decisión de spec):
	// el facade mapea este error de negocio a 400 VALIDATION_FAILED.
	if len(contacts) == 0 {
		return nil, errors.New("no_contacts_found")
	}

	return BuildContactsExcel(contacts)
}
