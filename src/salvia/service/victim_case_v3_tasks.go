package service

// Tareas, oficios y remisiones que crea el vaciado del borrador de Registro V3
// (WP-4, REG-38*; ver victim_case_v3_draft.go).
//
// Replican -sin editar esos archivos, que son de Seguimiento- lo que hace el
// formulario de Seguimiento v4 al registrarse el mismo dato:
//
//   - canales de activación de una entidad: form_service_v4.go v4ProcesarCanales
//   - gestión de una barrera nueva:         v4GestionBarreraNueva / v4AlertaComite
//   - remisiones a otros equipos:           form_service.go processEquiposYMedidas
//
// Todo se crea con mejor esfuerzo (como en Seguimiento): un fallo se registra en
// el log y no revierte las filas ya materializadas. Candidato a un helper
// compartido cuando Seguimiento y Registro V3 se unifiquen (ver «Para
// integración» en cambios-WP4.md).

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"log"
	"strings"
	"time"
)

func warnV3(format string, args ...any) {
	log.Printf("[victimCaseV3] WARN: "+format, args...)
}

func (s *victimCaseV3Service) crearTareaV3(ctx context.Context, task *models.CaseTask) {
	if s.draft == nil || s.draft.deps.CaseTaskRepo == nil {
		return
	}
	if err := s.draft.deps.CaseTaskRepo.Create(ctx, task); err != nil {
		warnV3("no se pudo crear la tarea %s del caso %s: %v", task.Type, task.CaseID, err)
	}
}

// crearOficioV3 devuelve el id del oficio, o nil si no se pudo crear.
func (s *victimCaseV3Service) crearOficioV3(ctx context.Context, letter *models.EntityLetter) *string {
	if s.draft == nil || s.draft.deps.EntityLetterRepo == nil {
		return nil
	}
	if err := s.draft.deps.EntityLetterRepo.Create(ctx, letter); err != nil {
		warnV3("no se pudo crear el oficio del caso %s: %v", letter.CaseID, err)
		return nil
	}
	id := letter.ID
	return &id
}

// crearTareasMaterializacionV3 crea, para lo que acaba de materializarse:
// los canales de activación de cada entidad y la gestión de cada barrera.
func (s *victimCaseV3Service) crearTareasMaterializacionV3(ctx context.Context, icode string, entities []flushedEntityV3, barriers []flushedBarrierV3) {
	for _, fe := range entities {
		e := fe.Draft
		// Igual que Seguimiento: los canales solo operan si la víctima está de
		// acuerdo Y se hizo la debida diligencia. La entidad a la que ya acudió
		// solo ofrece el bloque de activación si le faltan acciones por completar
		// (D17); sin pendientes no hay nada que activar.
		if fe.EntityCaseID == "" || !e.Agreement || !e.DueDiligence {
			continue
		}
		if e.Kind != "activacion" && len(e.Pending) == 0 {
			continue
		}
		s.procesarCanalesV3(ctx, icode, e.By, fe.EntityCaseID, e.EntityBranchID, e.ActivationChannels)
	}
	for _, fb := range barriers {
		s.procesarGestionBarreraV3(ctx, icode, fb.Draft, fb.BarrierID)
	}
}

func (s *victimCaseV3Service) procesarCanalesV3(ctx context.Context, caseID, actorID, entityCaseID string, branch int64, channels []string) {
	if entityCaseID == "" || branch <= 0 {
		return
	}
	for _, canal := range channels {
		switch canal {
		case "notificacion":
			branchID := branch
			agent := actorID
			letterID := s.crearOficioV3(ctx, &models.EntityLetter{
				CaseID: caseID, State: models.EntityLetterStatePorProyectar,
				AgentID: &agent, EntityBranchID: &branchID,
			})
			ecID := entityCaseID
			s.crearTareaV3(ctx, &models.CaseTask{
				Category: "Entidades", Type: "proyectar_oficio",
				Description:    "Proyectar oficio de activación de ruta",
				AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
				CaseID: caseID, EntityCaseID: &ecID, EntityLetterID: letterID,
			})
		case "misalvia":
			ecID := entityCaseID
			s.crearTareaV3(ctx, &models.CaseTask{
				Category: "Entidades", Type: "solicitud_entidad",
				Description:    "Hacer solicitudes a la entidad sobre obligaciones pendientes",
				AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
				CaseID: caseID, EntityCaseID: &ecID,
			})
		case "articulacion":
			s.gestionLlamadaV3(ctx, caseID, actorID, entityCaseID, "Articulación institucional telefónica o por mensajería")
		case "denuncia":
			s.gestionLlamadaV3(ctx, caseID, actorID, entityCaseID, "Denuncia fácil")
		default:
			warnV3("canal de activación sin efecto: %q", canal)
		}
	}
}

func (s *victimCaseV3Service) gestionLlamadaV3(ctx context.Context, caseID, actorID, entityCaseID, description string) {
	ecID := entityCaseID
	s.crearTareaV3(ctx, &models.CaseTask{
		Category: "Entidades", Type: "gestion_llamada",
		Description: description, AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
		CaseID: caseID, EntityCaseID: &ecID,
	})
}

func (s *victimCaseV3Service) procesarGestionBarreraV3(ctx context.Context, caseID string, b DraftBarrierV3, barrierID string) {
	var branch *int64
	if b.EntityBranchID != 0 {
		id := b.EntityBranchID
		branch = &id
	}
	for _, g := range b.Management {
		switch g {
		case "activacion":
			agent := b.By
			letterID := s.crearOficioV3(ctx, &models.EntityLetter{
				BarrierID: barrierID, CaseID: caseID,
				State: models.EntityLetterStatePorProyectar, AgentID: &agent, EntityBranchID: branch,
			})
			id := barrierID
			s.crearTareaV3(ctx, &models.CaseTask{
				Category: "Barreras", Type: "proyectar_oficio",
				Description:    "Activación de ruta interinstitucional (Generar oficio)",
				AssignedUserID: b.By, Status: models.CaseTaskStatusToDo,
				CaseID: caseID, BarrierID: &id, EntityLetterID: letterID,
			})
		case "alerta":
			id := barrierID
			s.crearTareaV3(ctx, &models.CaseTask{
				Category: "Barreras", Type: "comite_caso",
				Description: "Alerta por barreras", AssignedUserID: "",
				Status: models.CaseTaskStatusToDo, CaseID: caseID, BarrierID: &id,
			})
		}
	}
}

// ── Remisiones a otros equipos (D16a) ───────────────────────────────────────

const (
	catEquiposV3        = "victim_case_form2_referral_team"
	catMedidasV3        = "victim_case_form2_emergency_measure"
	catCriteriosPsicoV3 = "victim_case_form2_psycho_criteria"
	catCriteriosHomV3   = "victim_case_form2_men_criteria"
	catServiciosDiscV3  = "victim_case_form2_disability_service"
	catCriteriosEstV3   = "victim_case_form2_stabiliz_criteria"
)

// Código de V3 -> slug que entienden Seguimiento y los equipos.
var (
	medidaV3ToTipo = map[string]string{ //nolint:gochecknoglobals
		"al": "alojamiento", "tr": "transporte", "am": "alimentacion",
		"ve": "vestuario", "ps": "apoyo_psico", "ot": "otras_me",
	}
	criterioPsicoV3ToSlug = map[string]string{ //nolint:gochecknoglobals
		"ob": "criterio_obligatorio", "cs": "conducta_suicida", "in": "interseccionalidad",
		"sr": "sin_ruta", "ct": "condiciones_territoriales", "sa": "sin_acceso_psico", "nv": "naturalizacion_vbg",
	}
	criterioPsicoV3Puntos = map[string]int{ //nolint:gochecknoglobals
		"cs": 3, "in": 2, "sr": 1, "ct": 1, "sa": 1, "nv": 1,
	}
	servicioDiscV3ToSlug = map[string]string{ //nolint:gochecknoglobals
		"ls": "apoyo_lsc", "ed": "enfoque_discapacidad",
	}
)

// codigosEnumV3 traduce una respuesta multi-selección (CSV de icodes) a los
// códigos cortos de su categoría.
func codigosEnumV3(cats map[string][]repository.EnumOptionV3, category, csv string) []string {
	if strings.TrimSpace(csv) == "" {
		return nil
	}
	byICode := map[string]string{}
	for _, o := range cats[category] {
		byICode[o.ICode] = o.Code
	}
	out := []string{}
	for _, ic := range splitCSVV3(csv) {
		if code, ok := byICode[ic]; ok {
			out = append(out, code)
		}
	}
	return out
}

func contieneV3(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

// crearRemisionesV3 crea los registros de remisión a otros equipos con las
// respuestas GS_5..GS_12 guardadas hasta ese momento (D16a). Mismos criterios de
// elegibilidad que processEquiposYMedidas de Seguimiento; si no se cumplen, la
// remisión se omite. Se ejecuta una sola vez por Caso (materializedAt).
func (s *victimCaseV3Service) crearRemisionesV3(ctx context.Context, icode, actorID string) {
	if s.draft == nil {
		return
	}
	answers, err := s.repo.GetAnswers(ctx, icode)
	if err != nil {
		warnV3("remisiones: no se pudieron leer las respuestas de %s: %v", icode, err)
		return
	}
	if answers["GS_5"] != "true" {
		return
	}
	cats, err := s.repo.ListEnumCategories(ctx)
	if err != nil {
		warnV3("remisiones: no se pudo leer el catálogo de %s: %v", icode, err)
		return
	}
	equipos := codigosEnumV3(cats, catEquiposV3, answers["GS_6"])
	if len(equipos) == 0 {
		return
	}
	deps := s.draft.deps
	var submittedBy *string
	if v := strings.TrimSpace(actorID); v != "" {
		submittedBy = &v
	}
	submittedByTeam := s.equipoDelAgenteV3(ctx, actorID)
	conMedidas := contieneV3(equipos, "me")

	for _, equipo := range equipos {
		switch equipo {
		case "ap": // Atención psicosocial
			// Igual que Seguimiento: incompatible con medidas de emergencia en la misma remisión.
			if conMedidas {
				log.Printf("[victimCaseV3] remisión psicosocial omitida en %s: incompatible con medidas de emergencia", icode)
				continue
			}
			criterios := codigosEnumV3(cats, catCriteriosPsicoV3, answers["GS_9"])
			puntos := 0
			for _, c := range criterios {
				puntos += criterioPsicoV3Puntos[c]
			}
			if !contieneV3(criterios, "ob") || puntos < 3 {
				log.Printf("[victimCaseV3] remisión psicosocial omitida en %s: no cumple criterios (obligatorio=%v, puntos=%d)", icode, contieneV3(criterios, "ob"), puntos)
				continue
			}
			if deps.PsychosocialSupportRepo == nil {
				continue
			}
			slugs := make([]string, 0, len(criterios))
			for _, c := range criterios {
				if slug, ok := criterioPsicoV3ToSlug[c]; ok {
					slugs = append(slugs, slug)
				}
			}
			crit := strings.Join(slugs, ",")
			ps := &models.PsychosocialSupport{
				CaseID: icode, Type: "derivacion", Status: models.PsychosocialSupportStatusAbierto,
				SubmittedBy: submittedBy, SubmittedByTeam: submittedByTeam, ReferralCriteria: &crit,
			}
			if err := deps.PsychosocialSupportRepo.Create(ctx, ps); err != nil {
				warnV3("no se pudo crear la remisión psicosocial de %s: %v", icode, err)
				continue
			}
			psID := ps.ID
			s.crearTareaV3(ctx, &models.CaseTask{
				Category: "Psicosocial", Type: "validar_remision",
				Description:    "Validar que remision a psicosocial es valida",
				AssignedUserID: "", Status: models.CaseTaskStatusToDo,
				CaseID: icode, PsychosocialSupportID: &psID,
			})
		case "ah": // Atención hombres
			if !contieneV3(codigosEnumV3(cats, catCriteriosHomV3, answers["GS_10"]), "ob") {
				log.Printf("[victimCaseV3] remisión a atención hombres omitida en %s: sin criterio", icode)
				continue
			}
			if deps.MenTeamRemisionRepo == nil {
				continue
			}
			if err := deps.MenTeamRemisionRepo.Create(ctx, &models.MenTeamRemision{
				CaseID: icode, Type: "derivacion", Status: "ACTIVE",
			}); err != nil {
				warnV3("no se pudo crear la remisión a atención hombres de %s: %v", icode, err)
			}
		case "di": // Discapacidad: un registro por servicio
			if deps.DiscapacidadRemisionRepo == nil {
				continue
			}
			for _, svc := range codigosEnumV3(cats, catServiciosDiscV3, answers["GS_11"]) {
				slug, ok := servicioDiscV3ToSlug[svc]
				if !ok {
					continue
				}
				if err := deps.DiscapacidadRemisionRepo.Create(ctx, &models.DiscapacidadRemision{
					CaseID: icode, Service: slug, Status: "ACTIVE",
				}); err != nil {
					warnV3("no se pudo crear la remisión a discapacidad (%s) de %s: %v", slug, icode, err)
				}
			}
		case "es": // Estabilización económica: al menos un criterio
			if len(codigosEnumV3(cats, catCriteriosEstV3, answers["GS_12"])) == 0 {
				log.Printf("[victimCaseV3] remisión a estabilización omitida en %s: sin criterios", icode)
				continue
			}
			if deps.EconomicStabilizationRepo == nil {
				continue
			}
			if err := deps.EconomicStabilizationRepo.Create(ctx, &models.EconomicStabilization{
				CaseID: icode, Type: "derivacion", Status: "ACTIVE",
			}); err != nil {
				warnV3("no se pudo crear la remisión a estabilización de %s: %v", icode, err)
			}
		case "me": // Medidas de emergencia: una cabecera y una medida por tipo
			tipos := []string{}
			for _, m := range codigosEnumV3(cats, catMedidasV3, answers["GS_7"]) {
				if tipo, ok := medidaV3ToTipo[m]; ok {
					tipos = append(tipos, tipo)
				}
			}
			if len(tipos) == 0 || deps.EmergencyMeasureReferralRepo == nil {
				continue
			}
			now := time.Now()
			day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			referral := &models.EmergencyMeasureReferral{
				CaseID: icode, Status: models.EmergencyReferralStatusPendienteDeRevision,
				Urgency: models.EmergencyReferralUrgencyNormal, SubmittedBy: submittedBy,
				SubmittedByTeam: submittedByTeam, ReferredAt: &day,
			}
			if err := deps.EmergencyMeasureReferralRepo.CreateWithMeasures(ctx, referral, tipos, false); err != nil {
				warnV3("no se pudo crear la remisión a medidas de emergencia de %s: %v", icode, err)
			}
		case "sd": // Salvia Dignidad: sin registro propio (igual que Seguimiento)
		}
	}
}

func (s *victimCaseV3Service) equipoDelAgenteV3(ctx context.Context, actorID string) *string {
	repo := s.draft.deps.AgentLightRepo
	if repo == nil || strings.TrimSpace(actorID) == "" {
		return nil
	}
	agent, err := repo.FindByICode(ctx, actorID)
	if err != nil {
		warnV3("no se pudo leer el equipo de %s: %v", actorID, err)
		return nil
	}
	if agent == nil {
		return nil
	}
	if team := strings.TrimSpace(agent.Team); team != "" {
		return &team
	}
	return nil
}
