package service

import (
	"bitsflow/internal/repository"
	"context"
	"errors"
	"testing"
	"time"
)

type mockContactsReportRepository struct {
	repository.ReportRepository
	findContactsFunc func(ctx context.Context, start, end time.Time) ([]repository.ContactReportDTO, error)
}

func (m *mockContactsReportRepository) FindContactsInDateRange(ctx context.Context, start, end time.Time) ([]repository.ContactReportDTO, error) {
	if m.findContactsFunc != nil {
		return m.findContactsFunc(ctx, start, end)
	}
	return nil, nil
}

func testRange() FollowUpsReportRange {
	return FollowUpsReportRange{
		StartDate: time.Now().AddDate(0, -1, 0),
		EndDate:   time.Now(),
	}
}

func TestReportService_GenerateConsolidatedContactsReport_Empty(t *testing.T) {
	mockRepo := &mockContactsReportRepository{
		findContactsFunc: func(ctx context.Context, start, end time.Time) ([]repository.ContactReportDTO, error) {
			return []repository.ContactReportDTO{}, nil
		},
	}

	svc := NewReportService(mockRepo)
	_, err := svc.GenerateConsolidatedContactsReport(context.Background(), testRange())
	if err == nil || err.Error() != "no_contacts_found" {
		t.Fatalf("se esperaba error 'no_contacts_found' al no encontrar reportes, se obtuvo: %v", err)
	}
}

func TestReportService_GenerateConsolidatedContactsReport_Error(t *testing.T) {
	mockRepo := &mockContactsReportRepository{
		findContactsFunc: func(ctx context.Context, start, end time.Time) ([]repository.ContactReportDTO, error) {
			return nil, errors.New("db down")
		},
	}

	svc := NewReportService(mockRepo)
	_, err := svc.GenerateConsolidatedContactsReport(context.Background(), testRange())
	if err == nil || err.Error() != "db down" {
		t.Fatalf("se esperaba propagación del error de repositorio, se obtuvo: %v", err)
	}
}

func TestReportService_GenerateConsolidatedContactsReport_WithData(t *testing.T) {
	creation := time.Date(2026, 6, 10, 14, 30, 0, 0, time.UTC)
	mockRepo := &mockContactsReportRepository{
		findContactsFunc: func(ctx context.Context, start, end time.Time) ([]repository.ContactReportDTO, error) {
			return []repository.ContactReportDTO{
				{
					VictimContactId:        1,
					VictimContactICode:     "CON-001",
					CreationDate:           creation,
					UpdateDate:             creation,
					StatusDescription:      "Recontacto",
					Names:                  "María",
					LastNames:              "Pérez",
					Latitude:               4.60971,
					Longitude:              -74.08175,
					FormType:               "f1",
					Form1Nick:              "Mari",
					Form1DocType:           "cc",
					Form1DocNumber:         "1032456789",
					Form1BirthDate:         "01/05/1990",
					Form1TownCode:          "11001000",
					Form1TownName:          "Bogotá D.C.",
					Form1Address:           "Calle 1 # 2-3",
					Form1Phone:             "3001234567",
					Form1GenderIdentity:    "fe",
					Form1SexualOrientation: "he",
					Form1Origin:            "u",
					Form1Occupation:        "st",
					Form1FactsDescription:  "Descripción de los hechos",
				},
				{
					VictimContactId:          2,
					VictimContactICode:       "CON-002",
					CreationDate:             creation,
					UpdateDate:               creation,
					StatusDescription:        "Recontacto",
					FormType:                 "f2",
					Form2ReporterNames:       "Pedro Gómez",
					Form2ReporterPhone:       "3109876543",
					Form2VictimColPhone:      "3011112222",
					Form2FactsDescription:    "Reporte de un tercero",
					Form2BestContactTime: "14:00",
					Form2WillReceiveCall: "yes_no_y",
					Form2HasCareRole:     "yes_no_n",
					Form2VictimAware:     "yes_no_y",
					Form2ReportType:      "victim_case_form2_report_type_v",
					Form2AdjustmentsGBV:  "victim_case_form2_adjustments_gbv_nr, yes_no_y",
				},
			}, nil
		},
	}

	svc := NewReportService(mockRepo)
	file, err := svc.GenerateConsolidatedContactsReport(context.Background(), testRange())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) != 1 || sheets[0] != "Reportes" {
		t.Fatalf("se esperaba una única hoja 'Reportes', se obtuvo: %v", sheets)
	}

	rows, err := file.GetRows("Reportes")
	if err != nil {
		t.Fatalf("error leyendo filas: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("se esperaban 3 filas (1 encabezado + 2 reportes), se obtuvo: %d", len(rows))
	}

	// Las 34 cabeceras en el orden del spec
	if len(rows[0]) != 34 {
		t.Fatalf("se esperaban 34 cabeceras, se obtuvo: %d", len(rows[0]))
	}
	if rows[0][0] != "contacto_id" || rows[0][9] != "tipo_formulario_reporte" || rows[0][33] != "f2_ajustes_gbv" {
		t.Fatalf("cabeceras fuera de orden: %v", rows[0])
	}

	getCell := func(row []string, idx int) string {
		if idx < len(row) {
			return row[idx]
		}
		return ""
	}

	// Fila Form1: tipo diferenciado, códigos traducidos, f2_* vacías
	f1Row := rows[1]
	if getCell(f1Row, 9) != "Reporte propio (Formulario 1)" {
		t.Errorf("tipo_formulario_reporte F1 incorrecto: %q", getCell(f1Row, 9))
	}
	if getCell(f1Row, 11) != "Cédula de Ciudadanía" {
		t.Errorf("f1_tipo_documento sin traducir: %q", getCell(f1Row, 11))
	}
	if getCell(f1Row, 18) != "Mujer" {
		t.Errorf("f1_identidad_genero sin traducir: %q", getCell(f1Row, 18))
	}
	if getCell(f1Row, 20) != "Urbano" {
		t.Errorf("f1_procedencia_zona sin traducir: %q", getCell(f1Row, 20))
	}
	if getCell(f1Row, 21) != "Estudiante" {
		t.Errorf("f1_ocupacion sin traducir: %q", getCell(f1Row, 21))
	}
	if getCell(f1Row, 12) != "1032456789" {
		t.Errorf("f1_numero_documento debe ser texto sin alterar: %q", getCell(f1Row, 12))
	}
	if getCell(f1Row, 15) != "Bogotá D.C." {
		t.Errorf("f1_municipio incorrecto: %q", getCell(f1Row, 15))
	}
	if getCell(f1Row, 24) != "" || getCell(f1Row, 25) != "" {
		t.Errorf("columnas f2_* deben estar vacías en un reporte F1")
	}

	// Fila Form2: tipo diferenciado, f1_* vacías, teléfonos como texto
	f2Row := rows[2]
	if getCell(f2Row, 9) != "Reporte de tercero (Formulario 2)" {
		t.Errorf("tipo_formulario_reporte F2 incorrecto: %q", getCell(f2Row, 9))
	}
	if getCell(f2Row, 10) != "" || getCell(f2Row, 11) != "" {
		t.Errorf("columnas f1_* deben estar vacías en un reporte F2")
	}
	if getCell(f2Row, 25) != "3109876543" {
		t.Errorf("f2_telefono_reportante debe ser texto sin notación científica: %q", getCell(f2Row, 25))
	}
	if getCell(f2Row, 29) != "Sí" {
		t.Errorf("f2_recibira_llamada sin traducir desde clave yes_no_y: %q", getCell(f2Row, 29))
	}
	if getCell(f2Row, 30) != "No" {
		t.Errorf("f2_rol_cuidado sin traducir desde clave yes_no_n: %q", getCell(f2Row, 30))
	}
	if getCell(f2Row, 32) != "Usted es la víctima de violencia basada en género" {
		t.Errorf("f2_tipo_reporte sin traducir desde clave: %q", getCell(f2Row, 32))
	}
	if getCell(f2Row, 33) != "No requiere, Sí" {
		t.Errorf("f2_ajustes_gbv sin traducir claves de la lista: %q", getCell(f2Row, 33))
	}

	// Fechas amigables dd/mm/yyyy hh:mm
	if getCell(f1Row, 2) != "10/06/2026 14:30" {
		t.Errorf("contacto_fecha_creacion con formato incorrecto: %q", getCell(f1Row, 2))
	}
}
