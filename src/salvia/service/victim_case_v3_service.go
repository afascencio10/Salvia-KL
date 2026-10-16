package service

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

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
	Finalize(ctx context.Context, icode string) (*FinalizeResultV3, error)
	GetCase(ctx context.Context, icode string) (*CaseStateV3, error)
	EnumCatalog(ctx context.Context) (map[string][]repository.EnumOptionV3, error)

	// AddEntity/AddBarrier -- reconstruidos sobre el modelo real (entity_case/
	// entity_case_obligation/barrier_v2), el mismo que usan la pestaña
	// "Gestión institucional" y el form de Seguimiento. Ver
	// registro-v3-pendientes.md § B.
	AddEntity(ctx context.Context, e repository.VictimCaseEntityInput) (string, error)
	RemoveEntity(ctx context.Context, entityCaseID string) error
	ListEntities(ctx context.Context, icode string) ([]models.EntityCaseWithRelations, error)
	AddBarrier(ctx context.Context, b repository.VictimCaseBarrierInput) (string, error)
	RemoveBarrier(ctx context.Context, barrierID string) error
	ListBarriers(ctx context.Context, icode string) ([]models.BarrierV2, error)
}

// AddEntity devuelve ErrEntityCaseDuplicate (entity_case_service.go, mismo
// paquete) cuando el caso ya tiene esa sede activa -- el controller la
// mapea a 409, igual que case-entities.

// CaseStateV3 -- estado completo de un Caso para repoblar el Formulario al
// reabrir un Borrador (header + respuestas + entidades + barreras).
//
// AnsweredKeys acompaña a Answers porque Answers (vía GetAnswers) refleja
// lo que hay guardado en BD, dummy o real -- indistinguibles para varios
// Kind (ver hallazgo H3 en registro-v3-qa-plan.md). El frontend debe
// repoblar el formulario SOLO con las keys presentes en AnsweredKeys;
// cualquier otra key en Answers es un placeholder que nunca respondió el
// agente y no debe mostrarse como si lo hubiera hecho.
type CaseStateV3 struct {
	Summary      *repository.VictimCaseV3Summary `json:"summary"`
	Answers      map[string]string               `json:"answers"`
	AnsweredKeys []string                        `json:"answeredKeys"`
	// TriadaCompleta -- A2: antes solo se sabía tras un guardado en vivo
	// (SaveSection); un Caso ya completo se reabría con el sidebar
	// bloqueado otra vez porque el frontend nunca recibía este dato al
	// cargar. Mismo criterio que SaveSection (isTriadaCompleteV3).
	TriadaCompleta bool                             `json:"triadaCompleta"`
	Entities       []models.EntityCaseWithRelations `json:"entities"`
	Barriers       []models.BarrierV2               `json:"barriers"`
	// TzResultado -- BUG encontrado en QA (2026-10-15): antes nil siempre;
	// el Cierre mostraba "Nivel de riesgo: —" en cualquier Caso cuyo
	// Tamizaje no se hubiera guardado en la MISMA sesión de navegador que
	// hizo Finalizar (tzResultado vivía solo en memoria del frontend). Se
	// repuebla aquí desde lo que SaveTamizajeBranch ya persiste.
	TzResultado *TamizajeV3Result `json:"tzResultado"`
}

type TamizajeV3Result struct {
	RiskType  string `json:"riskType"`  // "feminicidio" | "homicidio_abyecto"
	RiskLevel string `json:"riskLevel"` // "bajo" | "moderado" | "alto" | "extremo"
	Score     int    `json:"score"`
	TieFlag   bool   `json:"tieFlag"` // ambos subtipos de tentativa marcados a la vez (G-V3-02)
}

type victimCaseV3Service struct {
	repo repository.VictimCaseV3Repository

	// entityCaseRepo/entityCaseObligationRepo/barrierV2Repo -- acceso directo
	// al modelo real de Entidades/Barreras (el mismo que usan la pestaña
	// "Gestión institucional" y el form de Seguimiento), inyectados desde
	// main.go reusando las MISMAS instancias que ya construye esa feature,
	// no copias propias de V3. Ver registro-v3-pendientes.md § B.
	entityCaseRepo           repository.EntityCaseRepository
	entityCaseObligationRepo repository.EntityCaseObligationRepository
	barrierV2Repo            repository.BarrierV2Repository
}

func NewVictimCaseV3Service(
	repo repository.VictimCaseV3Repository,
	entityCaseRepo repository.EntityCaseRepository,
	entityCaseObligationRepo repository.EntityCaseObligationRepository,
	barrierV2Repo repository.BarrierV2Repository,
) VictimCaseV3Service {
	return &victimCaseV3Service{
		repo:                     repo,
		entityCaseRepo:           entityCaseRepo,
		entityCaseObligationRepo: entityCaseObligationRepo,
		barrierV2Repo:            barrierV2Repo,
	}
}

// perfilObligatoriasV3 -- D1/D9: "Registro completo" se evalúa solo contra
// CI+DP+CA (el Perfil de víctima), nunca contra HV/AG/TZ/GS (eso es
// propiedad de cada Caso, no del Perfil que se duplica/comparte entre
// Casos del mismo documento -- ver D1 en registro-v3-data.md). Es
// triadaFieldsV3 filtrado a CI_/DP_ + las obligatorias de CA y los [máx 2º
// contacto] de DP que sí cuentan para Finalizar (C1).
var perfilObligatoriasV3 = func() []string { //nolint:gochecknoglobals
	out := make([]string, 0, len(triadaFieldsV3))
	for _, k := range triadaFieldsV3 {
		if strings.HasPrefix(k, "CI_") || strings.HasPrefix(k, "DP_") {
			out = append(out, k)
		}
	}
	out = append(out, dpSegundoContactoFinalizarV3...)
	out = append(out, caObligatoriasFinalizarV3...)
	return out
}()

func (s *victimCaseV3Service) Search(ctx context.Context, docType, docNumber, phone string) ([]repository.VictimCaseV3Summary, error) {
	// D7: solo documento y/o teléfono, nunca por nombre.
	results, err := s.repo.FindByDocumentOrPhone(ctx, docType, docNumber, phone)
	if err != nil || len(results) == 0 {
		return results, err
	}

	// D1/D9: "Registro completo" se calcula una sola vez, sobre el Caso MÁS
	// RECIENTE del documento -- se propaga el mismo valor a todas las filas
	// devueltas que comparten ese documento (puede haber varias, una por
	// Caso). Si la búsqueda fue solo por teléfono (sin documento), cada fila
	// puede tener un documento distinto -- se agrupa por (docType,docNumber).
	cache := map[string]bool{}
	for i := range results {
		key := results[i].DocType + "|" + results[i].DocNumber
		completo, ok := cache[key]
		if !ok {
			icode, err := s.repo.MostRecentCaseByDocument(ctx, results[i].DocType, results[i].DocNumber)
			if err != nil || icode == "" {
				completo = false
			} else {
				answered, err := s.repo.GetAnsweredKeys(ctx, icode)
				completo = err == nil && isSubsetAnsweredV3(answered, perfilObligatoriasV3)
			}
			cache[key] = completo
		}
		results[i].RegistroCompleto = completo
	}
	return results, nil
}

func isSubsetAnsweredV3(answered map[string]bool, keys []string) bool {
	for _, k := range keys {
		if !answered[k] {
			return false
		}
	}
	return true
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

// isTriadaCompleteV3 -- true si todas las obligatorias de triadaFieldsV3 ya
// fueron respondidas alguna vez (answered, de GetAnsweredKeys). Extraída de
// SaveSection para que GetCase (A2: triadaCompleta al reabrir un Caso) use
// exactamente el mismo criterio, en vez de depender de que el frontend haya
// recibido un `.then` de guardado en la misma sesión del navegador.
func isTriadaCompleteV3(answered map[string]bool) bool {
	for _, key := range triadaFieldsV3 {
		if !answered[key] {
			return false
		}
	}
	return true
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
	//
	// Usa GetAnsweredKeys, NO GetAnswers -- los valores dummy (D8) son
	// indistinguibles de una respuesta real una vez guardados (ej. dummy
	// "No" en bool_enum == respuesta real "No"), así que inferir
	// "respondido" desde el VALOR en BD marcaría la tríada completa
	// inmediatamente después de CreateDraft, sin que el agente respondiera
	// nada. GetAnsweredKeys rastrea qué keys llegaron alguna vez en un
	// guardado real (ver markKeysAnswered), independiente del valor.
	answered, err := s.repo.GetAnsweredKeys(ctx, icode)
	if err != nil {
		return false, err
	}
	complete := isTriadaCompleteV3(answered)
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

// riskTypeFromCodeV3/riskLevelFromCodeV3 -- inversas de riskTypeCodeV3/
// riskLevelCodeV3, para repoblar TzResultado en GetCase a partir de los
// codes que GetTamizajeResult relee de BD.
func riskTypeFromCodeV3(code string) string {
	if code == "ha" {
		return "homicidio_abyecto"
	}
	return "feminicidio"
}

func riskLevelFromCodeV3(code string) string {
	switch code {
	case "mo":
		return "moderado"
	case "al":
		return "alto"
	case "ex":
		return "extremo"
	default:
		return "bajo"
	}
}

// tamizajeInicialesFinalizarV3 -- de las 5 preguntas obligatorias de gating
// (D4/B4). El resto de Tamizaje (TZ_5/TZ_6/TZ_8 y las 24 de la rama pareja
// o no-pareja) NO entra en la validación de Finalizar: están condicionadas
// a que alguna de estas 5 sea "Sí" (`tzRevelarResto`), y la rama misma se
// decide en tiempo de ejecución (`wasPartner`) -- exigirlas siempre haría
// imposible finalizar en el caso común de que ninguna inicial sea "Sí".
var tamizajeInicialesFinalizarV3 = []string{ //nolint:gochecknoglobals
	"TZ_1", "TZ_2", "TZ_3", "TZ_4", "TZ_7",
}

// gsObligatoriasFinalizarV3 -- preguntas de GS que SIEMPRE están visibles
// (no dependen de otra respuesta). GS_16 se agregó en esta ronda (C1,
// 2026-10-15) -- existía en registro-v3-secciones.md como [1er contacto]
// pero no tenía ninguna UI todavía. Deliberadamente sigue excluyendo: GS_3
// (solo si GS_2 tiene selección), GS_6/7/9-12 (solo si GS_5="Sí", y varían
// según equipo elegido), GS_13 (solo si GS_5="Sí"), GS_15 (definida en el
// catálogo de campos pero sin UI todavía -- exigirla haría imposible
// finalizar cualquier Caso), GS_19 (solo si GS_18="Sí").
var gsObligatoriasFinalizarV3 = []string{ //nolint:gochecknoglobals
	"GS_1", "GS_2", "GS_4", "GS_5", "GS_14", "GS_16", "GS_17", "GS_18",
}

// caObligatoriasFinalizarV3 -- CA_1..CA_4 (máx. 2º contacto). C1 (2026-10-15):
// antes `registro-v3-campos.js` las tenía marcadas por error como [Deseable]
// (req: "D"), así que CA quedaba fuera de obligatoriasFinalizarV3 a pesar de
// que registro-v3-secciones.md documenta el cambio de criterio explícito
// "ya no es opcional" vs. V2-real. Se corrigió el catálogo (CA_1..4 ahora
// req: "C2") y se agrega esta lista acá.
var caObligatoriasFinalizarV3 = []string{ //nolint:gochecknoglobals
	"CA_1", "CA_2", "CA_3", "CA_4",
}

// dpSegundoContactoFinalizarV3 -- las [máx. 2º contacto] de DP que NO
// dependen de otra respuesta (DP_10 sexo asignado al nacer, DP_15 estado
// civil). El resto de C2 en DP es condicional (DP_3/DP_6/DP_26) y se deja
// fuera por la misma razón que el resto de condicionales de
// obligatoriasFinalizarV3 -- ver triadaFieldsV3.
var dpSegundoContactoFinalizarV3 = []string{ //nolint:gochecknoglobals
	"DP_10", "DP_15",
}

// obligatoriasFinalizarV3 -- C1 (2026-10-15): unión de TODAS las
// obligatorias [1er contacto]+[máx. 2º contacto] incondicionales de los 7
// módulos (antes solo cubría un subconjunto curado de [1er contacto] de
// CI/DP/HV/AG + TZ + GS, sin CA y sin los [máx. 2º contacto] de DP).
// triadaFieldsV3 ya sigue probado para Borrador→Iniciado, se reutiliza tal
// cual. Las condicionales (ej. AG Q5, DP Q13) se siguen excluyendo en todos
// los grupos, consistente con el criterio ya documentado en triadaFieldsV3:
// si la pregunta está oculta por su propia condición, no se le exige al
// agente.
var obligatoriasFinalizarV3 = append(append(append(append(append([]string{}, //nolint:gochecknoglobals
	triadaFieldsV3...), dpSegundoContactoFinalizarV3...), caObligatoriasFinalizarV3...),
	tamizajeInicialesFinalizarV3...), gsObligatoriasFinalizarV3...)

// ErrModulosIncompletos -- Finalize la devuelve cuando faltan obligatorias
// por responder (flow-E15: "Validar formulario... SI hay obligatorias sin
// responder → mostrar en qué módulo(s) falta diligenciamiento"). El
// controller la distingue de un error genérico con `errors.As`.
type ErrModulosIncompletos struct {
	Modulos []string
}

func (e *ErrModulosIncompletos) Error() string {
	return fmt.Sprintf("módulos incompletos: %s", strings.Join(e.Modulos, ", "))
}

// Finalize -- BUG encontrado en QA (2026-10-15): antes de este fix llamaba
// a UpdateStatus sin ninguna validación, podía terminar un Caso que nunca
// pasó por Iniciado. Ahora exige las obligatorias de obligatoriasFinalizarV3
// (ver flow-E15 § "Validar formulario") y solo transiciona si el estado
// actual es 'in' -- igual criterio de "nunca revertir/saltar estados" que
// ya usa SaveSection para Borrador→Iniciado (UpdateStatusIfCurrent).
// FinalizeResultV3 -- C2/C3 (2026-10-15): flow-E15 pide "generar la Denuncia
// fácil si corresponde" y "confirmar la remisión a procesos especializados
// si aplica". Ninguna de las dos tiene una integración real que replicar
// (a diferencia de Entidades/Barreras, no existe en todo el código ninguna
// mención de "denuncia fácil"/FGN, ni una tabla de remisión separada de la
// respuesta GS_6 misma) -- inventar un workflow nuevo sin que el cliente lo
// haya especificado sería peor que no hacer nada. Se implementa como lo que
// SÍ se puede confirmar honestamente: Finalize devuelve estos datos para que
// el agente vea una confirmación explícita de lo que va a pasar, en vez de
// silenciarlo. GS_6 queda crudo (CSV de icodes) porque el frontend ya sabe
// resolverlo a nombres de equipo (etiquetaSeleccionada/opcionesPara) -- no
// hace falta duplicar esa resolución de enums en Go para un mensaje.
type FinalizeResultV3 struct {
	DenunciaFacilSolicitada bool   `json:"denunciaFacilSolicitada"` // GS_1 == "true"
	EquiposRemisionCSV      string `json:"equiposRemisionCSV"`      // GS_6 crudo, "" si GS_5 != Sí
}

func (s *victimCaseV3Service) Finalize(ctx context.Context, icode string) (*FinalizeResultV3, error) {
	answered, err := s.repo.GetAnsweredKeys(ctx, icode)
	if err != nil {
		return nil, err
	}

	faltantesPorModulo := map[string]bool{}
	for _, key := range obligatoriasFinalizarV3 {
		if !answered[key] {
			modulo := strings.SplitN(key, "_", 2)[0]
			faltantesPorModulo[modulo] = true
		}
	}
	if len(faltantesPorModulo) > 0 {
		modulos := make([]string, 0, len(faltantesPorModulo))
		for m := range faltantesPorModulo {
			modulos = append(modulos, m)
		}
		sort.Strings(modulos)
		return nil, &ErrModulosIncompletos{Modulos: modulos}
	}

	answers, err := s.repo.GetAnswers(ctx, icode)
	if err != nil {
		return nil, err
	}
	result := &FinalizeResultV3{DenunciaFacilSolicitada: answers["GS_1"] == "true"}
	if answers["GS_5"] == "true" {
		result.EquiposRemisionCSV = answers["GS_6"]
	}

	// C4: no hay traspaso que ejecutar acá -- Registro V3 escribe directo
	// sobre salvia.victim_case/victim_case_form2, las MISMAS tablas que lee
	// Seguimiento (ver internal/repository/cases_list_repository.go, sin
	// tabla ni copia intermedia). 'te' (Terminado) es la señal de traspaso:
	// en cuanto se setea, el Caso ya es visible/consultable por Seguimiento
	// con el Perfil+tríada+Tamizaje+acciones que V3 ya escribió -- no existe
	// un segundo datastore que sincronizar (ver D1/D8/D9,
	// registro-v3-data.md). La transición Terminado→Cerrado es propiedad de
	// Seguimiento, no de esta pantalla (D5 obs. #1).
	if err := s.repo.UpdateStatusIfCurrent(ctx, icode, "in", "te"); err != nil {
		return nil, err
	}
	return result, nil
}

// AddEntity -- "+ Siguiente entidad" (GS § C), reconstruido sobre
// salvia.entity_case/entity_case_obligation (D9): EntityBranchID ya viene
// resuelto por la búsqueda Departamento/Ciudad/Municipio + Sector del
// frontend (idéntica a case-entities.js / Seguimiento), no por texto libre.
func (s *victimCaseV3Service) AddEntity(ctx context.Context, e repository.VictimCaseEntityInput) (string, error) {
	exists, err := s.entityCaseRepo.ExistsActive(ctx, e.VictimCaseICode, e.EntityBranchID)
	if err != nil {
		return "", fmt.Errorf("addEntity: verificar duplicado: %w", err)
	}
	if exists {
		return "", ErrEntityCaseDuplicate
	}

	var objetivo, solicitud, radicado *string
	if v := strings.TrimSpace(e.OtherInfo); v != "" {
		objetivo = &v
	}
	if v := strings.TrimSpace(e.ActivationChannel); v != "" {
		// El canal de activación no tiene columna propia en entity_case (a
		// diferencia del repeater real de Seguimiento, que además dispara
		// entity_letter/case_task vía procesarCanalActivacion) -- se deja
		// como nota en Solicitud para no perder el dato. Ver
		// registro-v3-pendientes.md § B (GAP canal de activación).
		s2 := "Canal de activación: " + v
		solicitud = &s2
	}
	if v := strings.TrimSpace(e.FileNumber); v != "" {
		radicado = &v
	}

	ec := &models.EntityCase{
		CaseID:         e.VictimCaseICode,
		EntityBranchID: e.EntityBranchID,
		Objetivo:       objetivo,
		Solicitud:      solicitud,
		Radicado:       radicado,
		CreatedByID:    e.Agent,
	}
	if err := s.entityCaseRepo.Create(ctx, ec); err != nil {
		return "", fmt.Errorf("addEntity: crear entity_case: %w", err)
	}

	if s.entityCaseObligationRepo != nil {
		s.registrarObligacionesEntidad(ctx, ec.ID, e.CompletedActions, models.EntityCaseObligationStatusCompletada, e.Agent)
		s.registrarObligacionesEntidad(ctx, ec.ID, e.PendingActions, models.EntityCaseObligationStatusPendiente, e.Agent)
	}
	return ec.ID, nil
}

// registrarObligacionesEntidad crea un entity_case_obligation por cada id
// del catálogo marcado -- mismo patrón que
// salvia/service/victim_case_entidades.go (ProcessVictimCaseEntidadEntries)
// para Registro de Caso V2. "otra" se ignora -- el checklist de V3 no
// captura el texto libre de "Otra" todavía (el frontend ya la excluye de
// las opciones, igual que case-entities.js).
func (s *victimCaseV3Service) registrarObligacionesEntidad(ctx context.Context, entityCaseID string, ids []string, status, actorID string) {
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" || id == "otra" {
			continue
		}
		existe, err := s.entityCaseObligationRepo.ExistsForObligation(ctx, entityCaseID, id)
		if err != nil {
			log.Printf("[registrarObligacionesEntidad] WARN: no se pudo verificar duplicado (entity_case=%s, obligation=%s): %v", entityCaseID, id, err)
			continue
		}
		if existe {
			continue
		}
		obligID := id
		eco := &models.EntityCaseObligation{
			EntityCaseID:       entityCaseID,
			EntityObligationID: &obligID,
			Status:             status,
			CreatedByID:        actorID,
		}
		if err := s.entityCaseObligationRepo.Create(ctx, eco); err != nil {
			log.Printf("[registrarObligacionesEntidad] WARN: no se pudo crear entity_case_obligation (entity_case=%s, obligation=%s): %v", entityCaseID, id, err)
		}
	}
}

func (s *victimCaseV3Service) RemoveEntity(ctx context.Context, entityCaseID string) error {
	return s.entityCaseRepo.Delete(ctx, entityCaseID)
}

func (s *victimCaseV3Service) ListEntities(ctx context.Context, icode string) ([]models.EntityCaseWithRelations, error) {
	return s.entityCaseRepo.FindByCaseIDWithRelations(ctx, icode)
}

// gestionBarreraV3CodeToSlug -- V3 muestra el catálogo seeded en
// victim_case_form2_enums (category victim_case_form2_barrier_management,
// codes "01".."06") porque ya tiene traducción en Locale.go, pero el valor
// real que entienden Seguimiento/Psicosocial (y que algún día dispararía
// case_task/entity_letter si V3 llegara a replicar esa automatización, ver
// AddBarrier) es un slug distinto -- gestionLabels/gestionSinTarea en
// salvia/service/psicosocial_barreras.go. Traducido acá, en el borde, para
// que barrier_v2.management_actions quede en el mismo vocabulario que
// cualquier barrera creada desde esas otras dos pantallas.
var gestionBarreraV3CodeToSlug = map[string]string{
	"01": "orientacion_llamada",
	"02": "gestion_llamada",
	"03": "activacion_ruta_interinstitucional",
	"04": "articulacion_institucional",
	"05": "escalamiento_organismo_control",
	"06": "alerta_barreras",
}

// AddBarrier -- "+ Siguiente barrera" (GS § D), reconstruido sobre
// salvia.barrier_v2 (D9) -- el mismo modelo real que usan Seguimiento y la
// sesión psicosocial (ver salvia/service/psicosocial_barreras.go). No
// replica la automatización de case_task/entity_letter que dispara
// "Gestión de la barrera" en esos dos flujos -- ver
// registro-v3-pendientes.md § B4 (fuera de alcance de esta ronda: no existe
// Seguimiento todavía sobre un Caso que recién se está registrando).
func (s *victimCaseV3Service) AddBarrier(ctx context.Context, b repository.VictimCaseBarrierInput) (string, error) {
	var entityBranchID *int64
	if b.EntityBranchID != 0 {
		id := b.EntityBranchID
		entityBranchID = &id
	}
	gestionSlugs := make([]string, 0, len(b.Management))
	for _, code := range b.Management {
		if slug, ok := gestionBarreraV3CodeToSlug[code]; ok {
			gestionSlugs = append(gestionSlugs, slug)
		} else {
			gestionSlugs = append(gestionSlugs, code)
		}
	}
	barrier := &models.BarrierV2{
		CaseID:                  b.VictimCaseICode,
		CreatedByID:             b.Agent,
		Status:                  models.BarrierV2StatusOpen,
		Sector:                  b.SectorCode,
		SpecificBarriers:        strings.Join(b.Specific, ","),
		OtherBarrierDesc:        b.OtherDetail,
		InstitutionName:         b.EntityName,
		DepartmentID:            b.Department,
		CityID:                  b.City,
		TownID:                  b.Town,
		StructuralInstitutional: strings.Join(b.Institutional, ","),
		StructuralEconomic:      strings.Join(b.Economic, ","),
		StructuralTerritorial:   strings.Join(b.Territorial, ","),
		StructuralDifferential:  strings.Join(b.Differential, ","),
		BarrierDate:             b.Date,
		OfficialDependency:      b.Official,
		Description:             b.Description,
		ManagementActions:       strings.Join(gestionSlugs, ","),
		Competences:             strings.Join(b.Competences, ","),
		EntityBranchID:          entityBranchID,
	}
	if err := s.barrierV2Repo.Create(ctx, barrier); err != nil {
		return "", fmt.Errorf("addBarrier: crear barrier_v2: %w", err)
	}
	return barrier.ID, nil
}

func (s *victimCaseV3Service) RemoveBarrier(ctx context.Context, barrierID string) error {
	return s.barrierV2Repo.Delete(ctx, barrierID)
}

func (s *victimCaseV3Service) ListBarriers(ctx context.Context, icode string) ([]models.BarrierV2, error) {
	return s.barrierV2Repo.FindByCaseID(ctx, icode)
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
	answeredSet, err := s.repo.GetAnsweredKeys(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	answeredKeys := make([]string, 0, len(answeredSet))
	for k := range answeredSet {
		answeredKeys = append(answeredKeys, k)
	}
	entities, err := s.ListEntities(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	barriers, err := s.ListBarriers(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	riskTypeCode, riskLevelCode, riskScore, err := s.repo.GetTamizajeResult(ctx, icode)
	if err != nil {
		return nil, fmt.Errorf("getCase: %w", err)
	}
	var tzResultado *TamizajeV3Result
	if riskLevelCode != "" {
		tzResultado = &TamizajeV3Result{
			RiskType:  riskTypeFromCodeV3(riskTypeCode),
			RiskLevel: riskLevelFromCodeV3(riskLevelCode),
			Score:     riskScore,
		}
	}
	return &CaseStateV3{
		Summary: summary, Answers: answers, AnsweredKeys: answeredKeys,
		TriadaCompleta: isTriadaCompleteV3(answeredSet),
		Entities:       entities, Barriers: barriers,
		TzResultado: tzResultado,
	}, nil
}

// EnumCatalog -- ver ListEnumCategories. Cacheado en memoria del lado del
// frontend tras el primer fetch (la tabla victim_case_form2_enums es
// normativa, no cambia en caliente).
func (s *victimCaseV3Service) EnumCatalog(ctx context.Context) (map[string][]repository.EnumOptionV3, error) {
	return s.repo.ListEnumCategories(ctx)
}
