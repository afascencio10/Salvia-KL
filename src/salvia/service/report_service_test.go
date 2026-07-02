package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
	"time"
)

type mockReportRepository struct {
	repository.ReportRepository
	findCasesFunc          func(ctx context.Context, start, end time.Time) ([]repository.CaseReportDTO, error)
	findFollowUpsFunc      func(ctx context.Context, caseICodes []string) ([]repository.FollowUpReportDTO, error)
	findAnswersFunc        func(ctx context.Context, submissionIDs []string) ([]repository.AnswerDTO, error)
	findTimelineEventsFunc func(ctx context.Context, caseICodes []string) ([]repository.TimelineReportDTO, error)
}

func (m *mockReportRepository) FindCasesInDateRange(ctx context.Context, start, end time.Time) ([]repository.CaseReportDTO, error) {
	if m.findCasesFunc != nil {
		return m.findCasesFunc(ctx, start, end)
	}
	return nil, nil
}

func (m *mockReportRepository) FindFollowUpsByCaseICodes(ctx context.Context, caseICodes []string) ([]repository.FollowUpReportDTO, error) {
	if m.findFollowUpsFunc != nil {
		return m.findFollowUpsFunc(ctx, caseICodes)
	}
	return nil, nil
}

func (m *mockReportRepository) FindAnswersByFormSubmissions(ctx context.Context, submissionIDs []string) ([]repository.AnswerDTO, error) {
	if m.findAnswersFunc != nil {
		return m.findAnswersFunc(ctx, submissionIDs)
	}
	return nil, nil
}

func (m *mockReportRepository) FindTimelineEventsByCaseICodes(ctx context.Context, caseICodes []string) ([]repository.TimelineReportDTO, error) {
	if m.findTimelineEventsFunc != nil {
		return m.findTimelineEventsFunc(ctx, caseICodes)
	}
	return nil, nil
}

func TestReportService_GenerateConsolidatedFollowUpsReport_Empty(t *testing.T) {
	mockRepo := &mockReportRepository{
		findCasesFunc: func(ctx context.Context, start, end time.Time) ([]repository.CaseReportDTO, error) {
			return []repository.CaseReportDTO{}, nil
		},
	}

	svc := NewReportService(mockRepo)
	r := FollowUpsReportRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}

	_, err := svc.GenerateConsolidatedFollowUpsReport(context.Background(), r)
	if err == nil || err.Error() != "no_cases_found" {
		t.Fatalf("se esperaba error 'no_cases_found' al no encontrar casos, se obtuvo: %v", err)
	}
}

func TestReportService_GenerateConsolidatedFollowUpsReport_WithData(t *testing.T) {
	age := 30
	mockRepo := &mockReportRepository{
		findCasesFunc: func(ctx context.Context, start, end time.Time) ([]repository.CaseReportDTO, error) {
			return []repository.CaseReportDTO{
				{
					VictimCaseId:             1,
					VictimCaseICode:          "CAS-001",
					VictimCaseStatus:         "ra",
					VictimCaseVictimName:     "María",
					VictimCaseVictimLastName: "Pérez",
					VictimCaseVictimAge:      &age,
					VictimCaseVictimGender:   "F",
					VictimCaseTownCode:       "11001000",
					VictimCaseTeam:           "Riesgo bajo",
					VictimCaseCreatedAt:      time.Now(),
					VictimPhone:              "3159999999",
					HasForm1:                 "Sí",
					HasForm2:                 "No",
					TownName:                 "Bogotá",
					DeptName:                 "Cundinamarca",
					AggressorName:            "Juan Pérez",
					RelationshipAggressor:    "Ex-pareja",
					RiskLevel:                "Alto",
				},
			}, nil
		},
		findFollowUpsFunc: func(ctx context.Context, caseICodes []string) ([]repository.FollowUpReportDTO, error) {
			subID := "submission-001"
			summary := "Primer contacto exitoso"
			risk := "bajo"
			return []repository.FollowUpReportDTO{
				{
					FollowUpV2: models.FollowUpV2{
						ID:               "FU-001",
						CaseID:           "CAS-001",
						FormSubmissionID: &subID,
						SequenceNumber:   1,
						ScheduledDate:    time.Now(),
						Status:           "REALIZADO",
						Team:             "Riesgo bajo",
						RiskStatus:       &risk,
						Summary:          &summary,
					},
					VictimDocNumber: "12345678",
				},
			}, nil
		},
		findAnswersFunc: func(ctx context.Context, submissionIDs []string) ([]repository.AnswerDTO, error) {
			now := time.Now()
			return []repository.AnswerDTO{
				{
					CaseICode:        "CAS-001",
					FollowUpID:       "FU-001",
					FormSubmissionID: "submission-001",
					QuestionDesc:     "¿Sufre amenazas?",
					AnswerValue:      "Sí",
					VictimDocNumber:  "12345678",
					FollowUpDate:     &now,
				},
			}, nil
		},
		findTimelineEventsFunc: func(ctx context.Context, caseICodes []string) ([]repository.TimelineReportDTO, error) {
			return []repository.TimelineReportDTO{
				{
					CaseTimelineEvent: models.CaseTimelineEvent{
						ID:          "EV-001",
						CaseID:      "CAS-001",
						Date:        time.Now(),
						Category:    "General",
						Type:        "Creación de Caso",
						Description: "Caso creado por primer contacto",
						ActorName:   "Supervisor Juan",
					},
					VictimDocNumber: "12345678",
				},
			}, nil
		},
	}

	svc := NewReportService(mockRepo)
	r := FollowUpsReportRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}

	excelFile, err := svc.GenerateConsolidatedFollowUpsReport(context.Background(), r)
	if err != nil {
		t.Fatalf("se esperaba éxito, error: %v", err)
	}
	if excelFile == nil {
		t.Fatalf("se esperaba archivo Excel válido")
	}
	defer excelFile.Close()

	// Validar que se puedan leer las celdas de las hojas enriquecidas
	val, _ := excelFile.GetCellValue("Casos", "A2")
	if val != "CAS-001" {
		t.Errorf("se esperaba CAS-001 en Casos A2, se obtuvo: %s", val)
	}

	// Género amigable
	val, _ = excelFile.GetCellValue("Casos", "G2")
	if val != "Mujer" {
		t.Errorf("se esperaba Mujer en Casos G2, se obtuvo: %s", val)
	}

	val, _ = excelFile.GetCellValue("Casos", "H2")
	if val != "3159999999" {
		t.Errorf("se esperaba teléfono en Casos H2, se obtuvo: %s", val)
	}

	val, _ = excelFile.GetCellValue("Casos", "I2")
	if val != "Bogotá" {
		t.Errorf("se esperaba Bogotá en Casos I2, se obtuvo: %s", val)
	}

	// Validar hoja "Registro"
	val, _ = excelFile.GetCellValue("Registro", "AJ2")
	if val != "Alto" {
		t.Errorf("se esperaba Alto en Registro AJ2, se obtuvo: %s", val)
	}
	val, _ = excelFile.GetCellValue("Registro", "AK2")
	if val != "Juan Pérez" {
		t.Errorf("se esperaba Juan Pérez en Registro AK2, se obtuvo: %s", val)
	}

	val, _ = excelFile.GetCellValue("Seguimientos", "B2")
	if val != "12345678" {
		t.Errorf("se esperaba documento víctima en Seguimientos B2, se obtuvo: %s", val)
	}

	val, _ = excelFile.GetCellValue("Respuestas de formulario", "F2")
	if val != "¿Sufre amenazas?" {
		t.Errorf("se esperaba la pregunta en F2, se obtuvo: %s", val)
	}

	val, _ = excelFile.GetCellValue("Timeline", "B2")
	if val != "12345678" {
		t.Errorf("se esperaba documento víctima en Timeline B2, se obtuvo: %s", val)
	}
}

func TestReportService_GenerateConsolidatedFollowUpsReport_Error(t *testing.T) {
	mockRepo := &mockReportRepository{
		findCasesFunc: func(ctx context.Context, start, end time.Time) ([]repository.CaseReportDTO, error) {
			return nil, errors.New("error de BD simulado")
		},
	}

	svc := NewReportService(mockRepo)
	r := FollowUpsReportRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}

	_, err := svc.GenerateConsolidatedFollowUpsReport(context.Background(), r)
	if err == nil {
		t.Fatalf("se esperaba un error de base de datos")
	}
}
