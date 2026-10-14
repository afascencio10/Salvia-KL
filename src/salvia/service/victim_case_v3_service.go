package service

import (
	"bitsflow/common/utils"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// VictimCaseV3Service — Registro V3 (DocsMD/Screens/Registro V3). Sin
// dinamic-form (decisión A1): el frontend llama directamente a estos
// métodos por módulo, no hay OnSectionUpdate genérico.
type VictimCaseV3Service interface {
	Search(ctx context.Context, docType, docNumber, phone string) ([]repository.VictimCaseV3Summary, error)
	CreateDraft(ctx context.Context, docType, docNumber, phone, residenceTown, agentID, channel string) (icode string, err error)
	SaveSection(ctx context.Context, icode, modulo string, answers map[string]string) (triadaCompleta bool, err error)
	SaveTamizaje(ctx context.Context, icode string, answers map[string]string, wasPartner bool, esOsigd bool, subtipoFeminicidio, subtipoHomicidioAbyecto bool) (TamizajeV3Result, error)
	Finalize(ctx context.Context, icode string) error
	GetCase(ctx context.Context, icode string) (*CaseStateV3, error)
	EnumCatalog(ctx context.Context) (map[string][]repository.EnumOptionV3, error)

	AddEntity(ctx context.Context, e repository.VictimCaseEntityInput) (string, error)
	RemoveEntity(ctx context.Context, entityICode string) error
	ListEntities(ctx context.Context, icode string) ([]repository.VictimCaseEntityInput, error)
	AddBarrier(ctx context.Context, b repository.VictimCaseBarrierInput) (string, error)
	RemoveBarrier(ctx context.Context, barrierICode string) error
	ListBarriers(ctx context.Context, icode string) ([]repository.VictimCaseBarrierInput, error)
}

// CaseStateV3 -- estado completo de un Caso para repoblar el Formulario al
// reabrir un Borrador (header + respuestas + entidades + barreras).
type CaseStateV3 struct {
	Summary  *repository.VictimCaseV3Summary     `json:"summary"`
	Answers  map[string]string                   `json:"answers"`
	Entities []repository.VictimCaseEntityInput  `json:"entities"`
	Barriers []repository.VictimCaseBarrierInput `json:"barriers"`
}

type TamizajeV3Result struct {
	RiskType  string `json:"riskType"`  // "feminicidio" | "homicidio_abyecto"
	RiskLevel string `json:"riskLevel"` // "bajo" | "moderado" | "alto" | "extremo"
	Score     int    `json:"score"`
	TieFlag   bool   `json:"tieFlag"` // ambos subtipos de tentativa marcados a la vez (G-V3-02)
}

type victimCaseV3Service struct {
	repo repository.VictimCaseV3Repository
}

func NewVictimCaseV3Service(repo repository.VictimCaseV3Repository) VictimCaseV3Service {
	return &victimCaseV3Service{repo: repo}
}

func (s *victimCaseV3Service) Search(ctx context.Context, docType, docNumber, phone string) ([]repository.VictimCaseV3Summary, error) {
	// D7: solo documento y/o teléfono, nunca por nombre.
	return s.repo.FindByDocumentOrPhone(ctx, docType, docNumber, phone)
}

func (s *victimCaseV3Service) CreateDraft(ctx context.Context, docType, docNumber, phone, residenceTown, agentID, channel string) (string, error) {
	login := randomLogin(6)
	plainPassword := randomPassword(4)
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	if err != nil {
		return "", fmt.Errorf("createDraft: generar hash: %w", err)
	}

	icode, err := s.repo.CreateDraft(ctx, repository.CreateDraftV3Input{
		DocType:            docType,
		DocNumber:          docNumber,
		Phone:              phone,
		ResidenceTownCode:  residenceTown,
		AgentID:            agentID,
		Channel:            channel,
		GeneralUserICode:   utils.GetUUID(),
		GeneralUserLogin:   login,
		GeneralUserPassSha: string(hashed),
	})
	if err != nil {
		return "", fmt.Errorf("createDraft: %w", err)
	}
	if err := s.repo.StoreCredentials(ctx, icode, login, plainPassword); err != nil {
		return "", fmt.Errorf("createDraft: guardar credenciales: %w", err)
	}
	log.Printf("[victimCaseV3] draft creado: caso=%s canal=%s agente=%s", icode, channel, agentID)
	return icode, nil
}

// triadaFieldsV3 -- obligatorias [1er contacto] de CI+DP+CA+HV+AG que
// determinan triadaCompleta (D4, unificado con Borrador→Iniciado). Solo
// el núcleo incondicional -- las variantes condicionales (ej. AG Q5 "depende
// económicamente", DP Q13 "pueblo indígena") no se evalúan aquí; se asume
// que si la pregunta está oculta por su propia condición, no se le exige al
// agente.
var triadaFieldsV3 = []string{ //nolint:gochecknoglobals
	"CI_1", "CI_2", "CI_3", "CI_4", "CI_5", "CI_6", "CI_7", "CI_9",
	"DP_1", "DP_2", "DP_5", "DP_8", "DP_9", "DP_11", "DP_12", "DP_14", "DP_17", "DP_27", "DP_28", "DP_29",
	"HV_1", "HV_2", "HV_4", "HV_6", "HV_8", "HV_9", "HV_11", "HV_13", "HV_14",
	"AG_1", "AG_2", "AG_3",
}

func (s *victimCaseV3Service) SaveSection(ctx context.Context, icode, modulo string, answers map[string]string) (bool, error) {
	if err := s.repo.SaveSection(ctx, icode, modulo, answers); err != nil {
		return false, err
	}

	// triadaCompleta / Borrador→Iniciado solo se reevalúa cuando el módulo
	// guardado es parte de la tríada (CI/DP/CA/HV/AG) -- TZ/GS no aportan.
	if modulo != "CI" && modulo != "DP" && modulo != "CA" && modulo != "HV" && modulo != "AG" {
		return false, nil
	}

	// Se evalúa contra el ACUMULADO histórico (todas las secciones ya
	// guardadas), no solo el delta de este guardado -- SaveSection es un
	// UPDATE parcial por módulo, así que un guardado de un solo módulo
	// nunca podría satisfacer por sí mismo campos de los otros 4.
	accumulated, err := s.repo.GetAnswers(ctx, icode)
	if err != nil {
		return false, err
	}
	complete := true
	for _, key := range triadaFieldsV3 {
		if v, ok := accumulated[key]; !ok || v == "" {
			complete = false
			break
		}
	}
	if complete {
		// 'in' = Iniciado (código nuevo, salvia/config/Enums.go) -- solo
		// transiciona si el Caso sigue en Borrador; no revierte un Caso que
		// ya esté Terminado/Cerrado por una reapertura parcial.
		if err := s.repo.UpdateStatusIfCurrent(ctx, icode, "bo", "in"); err != nil {
			return false, err
		}
	}
	return complete, nil
}

func (s *victimCaseV3Service) SaveTamizaje(ctx context.Context, icode string, answers map[string]string, wasPartner bool, esOsigd bool, subtipoFeminicidio, subtipoHomicidioAbyecto bool) (TamizajeV3Result, error) {
	// D3: tipo decidido por identidad/orientación (esOsigd), NO por la
	// pregunta VPP -- esa solo suma puntaje.
	riskType := "feminicidio"
	if esOsigd {
		riskType = "homicidio_abyecto"
	}

	tieFlag := false
	switch {
	case subtipoFeminicidio && subtipoHomicidioAbyecto:
		// G-V3-02 / B6: empate -- Feminicidio tiene prioridad.
		riskType = "feminicidio"
		tieFlag = true
	case subtipoFeminicidio:
		riskType = "feminicidio"
	case subtipoHomicidioAbyecto:
		riskType = "homicidio_abyecto"
	}

	score := 0
	for _, v := range answers {
		if v == "true" {
			score++
		}
	}

	level := computeRiskLevelV3(riskType, wasPartner, score)

	if err := s.repo.SaveTamizajeBranch(ctx, icode, wasPartner, answers, riskTypeCodeV3(riskType), riskLevelCodeV3(level), score, tieFlag); err != nil {
		return TamizajeV3Result{}, err
	}

	return TamizajeV3Result{RiskType: riskType, RiskLevel: level, Score: score, TieFlag: tieFlag}, nil
}

// computeRiskLevelV3 -- rangos de D4/Guion Funcional, marcados "pendientes
// de ajuste final" por el propio cliente (obs. #3) -- no recalibrar sin
// instrucción nueva.
func computeRiskLevelV3(riskType string, wasPartner bool, score int) string {
	type threshold struct {
		max   int
		level string
	}
	var thresholds []threshold
	switch {
	case riskType == "feminicidio" && wasPartner:
		thresholds = []threshold{{4, "bajo"}, {8, "moderado"}, {15, "alto"}, {1 << 30, "extremo"}}
	case riskType == "feminicidio" && !wasPartner:
		thresholds = []threshold{{2, "bajo"}, {5, "moderado"}, {8, "alto"}, {1 << 30, "extremo"}}
	default: // homicidio_abyecto, comparte escala pareja/no-pareja
		thresholds = []threshold{{5, "bajo"}, {11, "moderado"}, {16, "alto"}, {1 << 30, "extremo"}}
	}
	for _, t := range thresholds {
		if score <= t.max {
			return t.level
		}
	}
	return "extremo"
}

// riskTypeCodeV3/riskLevelCodeV3 -- traducen los valores descriptivos
// internos a los 2-char code sembrados en migrations/add_registro_v3.sql
// (categorías victim_case_form2_tzv3_risk_type/_risk_level).
func riskTypeCodeV3(riskType string) string {
	if riskType == "homicidio_abyecto" {
		return "ha"
	}
	return "fe"
}

func riskLevelCodeV3(level string) string {
	switch level {
	case "moderado":
		return "mo"
	case "alto":
		return "al"
	case "extremo":
		return "ex"
	default:
		return "ba"
	}
}

func (s *victimCaseV3Service) Finalize(ctx context.Context, icode string) error {
	return s.repo.UpdateStatus(ctx, icode, "te") // 'te' = Terminado -- requiere agregarse a VICTIM_CASE_STATUS, ver GAP abajo
}

func (s *victimCaseV3Service) AddEntity(ctx context.Context, e repository.VictimCaseEntityInput) (string, error) {
	return s.repo.AddEntity(ctx, e)
}
func (s *victimCaseV3Service) RemoveEntity(ctx context.Context, entityICode string) error {
	return s.repo.RemoveEntity(ctx, entityICode)
}
func (s *victimCaseV3Service) ListEntities(ctx context.Context, icode string) ([]repository.VictimCaseEntityInput, error) {
	return s.repo.ListEntities(ctx, icode)
}
func (s *victimCaseV3Service) AddBarrier(ctx context.Context, b repository.VictimCaseBarrierInput) (string, error) {
	return s.repo.AddBarrier(ctx, b)
}
func (s *victimCaseV3Service) RemoveBarrier(ctx context.Context, barrierICode string) error {
	return s.repo.RemoveBarrier(ctx, barrierICode)
}
func (s *victimCaseV3Service) ListBarriers(ctx context.Context, icode string) ([]repository.VictimCaseBarrierInput, error) {
	return s.repo.ListBarriers(ctx, icode)
}

// GetCase -- estado completo del Caso (header + respuestas + entidades +
// barreras) para repoblar el Formulario (ej. reabrir un Borrador).
func (s *victimCaseV3Service) GetCase(ctx context.Context, icode string) (*CaseStateV3, error) {
	summary, err := s.repo.GetCaseSummary(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	answers, err := s.repo.GetAnswers(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	entities, err := s.repo.ListEntities(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	barriers, err := s.repo.ListBarriers(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	return &CaseStateV3{Summary: summary, Answers: answers, Entities: entities, Barriers: barriers}, nil
}

// EnumCatalog -- ver ListEnumCategories. Cacheado en memoria del lado del
// frontend tras el primer fetch (la tabla victim_case_form2_enums es
// normativa, no cambia en caliente).
func (s *victimCaseV3Service) EnumCatalog(ctx context.Context) (map[string][]repository.EnumOptionV3, error) {
	return s.repo.ListEnumCategories(ctx)
}
