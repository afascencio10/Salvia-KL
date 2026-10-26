package service

import (
	"bitsflow/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// seguimientoV4ID reproduce id() de scripts/gen_seguimiento_v4_seed.js:
// UUID v4 a partir de sha256("seguimiento-v4:" + clave).
func seguimientoV4ID(key string) string {
	sum := sha256.Sum256([]byte("seguimiento-v4:" + key))
	var b [16]byte
	copy(b[:], sum[:16])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b[:])
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
}

func v4Key(key string) string { return seguimientoV4ID(key) }

type v4Entry struct {
	entry   models.RepeaterEntry
	answers map[string]string
}

type v4Loaded struct {
	direct  map[string]string
	byGroup map[string][]v4Entry
}

type v4BarrierQ struct {
	sector, tiposSalud, instSalud, otraSalud          string
	tiposJusticia, instJusticia, otraJusticia         string
	tiposProteccion, instProteccion, otraProteccion   string
	talento, eco, terr, diferencial                   string
	fecha, funcionario, depto, municipio, descripcion string
	gestion                                           string
}

func v4BarrierQuestions(prefix string) v4BarrierQ {
	k := func(field string) string { return v4Key("q:" + prefix + "." + field) }
	return v4BarrierQ{
		sector: k("sector"), tiposSalud: k("tiposSalud"), instSalud: k("instSalud"), otraSalud: k("otraTipoSalud"),
		tiposJusticia: k("tiposJusticia"), instJusticia: k("instJusticia"), otraJusticia: k("otraTipoJusticia"),
		tiposProteccion: k("tiposProteccion"), instProteccion: k("instProteccion"), otraProteccion: k("otraTipoProteccion"),
		talento: k("talento"), eco: k("eco"), terr: k("terr"), diferencial: k("diferencial"),
		fecha: k("fecha"), funcionario: k("funcionario"), depto: k("depto"), municipio: k("municipio"),
		descripcion: k("descripcion"), gestion: k("gestion"),
	}
}

var (
	v4IdentBarrier = v4BarrierQuestions("ident.barrera")
	v4SegBarrier   = v4BarrierQuestions("seg.barrera")
)

// v4Reasignacion es la decisión ya tomada. applyV4Reasignacion la ejecuta
// después de marcar el seguimiento REALIZADO.
type v4Reasignacion struct {
	apply bool
	level int
	team  string
}

func v4Warn(format string, args ...any) {
	log.Printf("[processFollowUpV4] "+format, args...)
}

func (s *formService) loadV4Submission(ctx context.Context, submissionID string) (*v4Loaded, error) {
	if s.answerRepo == nil || s.repeaterEntryRepo == nil {
		return nil, fmt.Errorf("processFollowUpV4: repositorios de respuestas no configurados")
	}
	directRows, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("processFollowUpV4: leer respuestas: %w", err)
	}
	loaded := &v4Loaded{
		direct:  make(map[string]string, len(directRows)),
		byGroup: map[string][]v4Entry{},
	}
	for _, a := range directRows {
		loaded.direct[a.QuestionID] = a.Value
	}

	groups := []string{
		v4Key("repeater:identEntidades"),
		v4Key("repeater:identBarreras"),
		v4Key("repeater:segEntidades"),
		v4Key("repeater:segBarreras"),
		v4Key("repeater:segBarrerasActivas"),
	}
	entries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, groups)
	if err != nil {
		return nil, fmt.Errorf("processFollowUpV4: leer entradas: %w", err)
	}
	for _, entry := range entries {
		rows, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
		if err != nil {
			return nil, fmt.Errorf("processFollowUpV4: leer entrada %s: %w", entry.ID, err)
		}
		answers := make(map[string]string, len(rows))
		for _, a := range rows {
			answers[a.QuestionID] = a.Value
		}
		loaded.byGroup[entry.RepeaterGroupID] = append(loaded.byGroup[entry.RepeaterGroupID], v4Entry{
			entry: entry, answers: answers,
		})
	}
	return loaded, nil
}

func v4Roots(entries []v4Entry) []v4Entry {
	var out []v4Entry
	for _, e := range entries {
		if e.entry.ParentEntryID == nil || *e.entry.ParentEntryID == "" {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].entry.Iteration < out[j].entry.Iteration })
	return out
}

func v4Kids(entries []v4Entry, parentID string) []v4Entry {
	var out []v4Entry
	for _, e := range entries {
		if e.entry.ParentEntryID != nil && *e.entry.ParentEntryID == parentID {
			out = append(out, e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].entry.Iteration < out[j].entry.Iteration })
	return out
}

func v4Has(csv, want string) bool {
	for _, v := range splitCSVValues(csv) {
		if v == want {
			return true
		}
	}
	return false
}

// v4CountedValues descarta Ninguna/Ninguno. En el V.4 el valor de la opción es "ninguna".
func v4CountedValues(csv string) []string {
	var out []string
	for _, v := range splitCSVValues(csv) {
		if v == "ninguna" || v == "ninguno" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func v4ParseBranch(raw string) (int64, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "otra" {
		return 0, false
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, false
	}
	return id, true
}

func v4BothYes(answers map[string]string, acuerdoKey, debidaKey string) bool {
	return answers[acuerdoKey] == "true" && answers[debidaKey] == "true"
}

// v4ReasignacionDestino decide el nivel de destino.
// Un factor extremo con el caso en nivel 1 o 2 gana sobre el desplegable.
// riesgo y general no mueven el caso si ya está en ese equipo.
// feminicidio y ninguno no mueven el caso.
func v4ReasignacionDestino(currentLevel int, extremosCSV, equipo string) v4Reasignacion {
	if len(v4CountedValues(extremosCSV)) > 0 && currentLevel >= 1 && currentLevel <= 2 {
		return v4Reasignacion{apply: true, level: 4, team: "Riesgo alto"}
	}
	switch strings.TrimSpace(equipo) {
	case "riesgo":
		if currentLevel >= 3 {
			return v4Reasignacion{}
		}
		if currentLevel >= 1 {
			return v4Reasignacion{apply: true, level: 3, team: "Riesgo alto"}
		}
	case "general":
		if currentLevel >= 1 && currentLevel <= 2 {
			return v4Reasignacion{}
		}
		if currentLevel >= 3 {
			return v4Reasignacion{apply: true, level: 2, team: "Riesgo bajo"}
		}
	}
	return v4Reasignacion{}
}

func v4PsicoAplica(csv string) bool {
	scores := map[string]int{
		"poblacion": 0, "suicida": 3, "temor": 1, "inter": 2,
		"territorio": 1, "sin_acceso": 1, "naturalizacion": 1,
	}
	hasPoblacion := false
	points := 0
	for _, v := range splitCSVValues(csv) {
		if v == "poblacion" {
			hasPoblacion = true
		}
		points += scores[v]
	}
	return hasPoblacion && points >= 3
}

func (s *formService) processFollowUpV4(ctx context.Context, fu *models.FollowUpV2, submissionID, actorID string) (int, int, v4Reasignacion, error) {
	loaded, err := s.loadV4Submission(ctx, submissionID)
	if err != nil {
		return 0, 0, v4Reasignacion{}, err
	}
	if loaded.direct == nil {
		loaded.direct = map[string]string{}
	}

	s.v4Hechos(ctx, fu, actorID, loaded.direct)

	dest := v4Reasignacion{}
	if s.caseRepo != nil {
		level, err := s.caseRepo.FindRiskLevelByICode(ctx, fu.CaseID)
		if err != nil {
			v4Warn("no se pudo leer el nivel de riesgo del caso %s: %v", fu.CaseID, err)
		} else {
			dest = v4ReasignacionDestino(level, loaded.direct[v4Key("q:risk.extremos")], loaded.direct[v4Key("q:plan.equipoReasigna")])
		}
	}

	barriers := 0
	remisiones := 0
	if loaded.direct[v4Key("q:ident.conocimiento")] == "true" {
		barriers += s.v4ProcessIdent(ctx, fu, actorID, loaded)
	}
	b, err := s.v4ProcessSeg(ctx, fu, actorID, loaded)
	if err != nil {
		v4Warn("seguimiento a entidades: %v", err)
	}
	barriers += b
	remisiones += s.v4ProcessPlan(ctx, fu, actorID, loaded.direct)
	s.v4ProcessConclusion(ctx, fu, actorID, loaded.direct)
	return barriers, remisiones, dest, nil
}

func (s *formService) v4Hechos(ctx context.Context, fu *models.FollowUpV2, actorID string, direct map[string]string) {
	if direct[v4Key("q:risk.hechos")] != "true" || s.caseTimelineRepo == nil {
		return
	}
	fecha := time.Now()
	if raw := strings.TrimSpace(direct[v4Key("q:risk.hechosFecha")]); raw != "" {
		if t, err := time.Parse("2006-01-02", raw); err == nil {
			fecha = t
		}
	}
	event := &models.CaseTimelineEvent{
		CaseID:      fu.CaseID,
		FollowUpID:  fu.ID,
		Category:    models.TimelineCategoryGeneral,
		Type:        models.TimelineTypeHechosCaso,
		Icon:        models.TimelineIconHechosCaso,
		Color:       "#f87171",
		Description: strings.TrimSpace(direct[v4Key("q:risk.hechosTexto")]),
		EventUserID: actorID,
		Date:        fecha,
		CreatedAt:   time.Now(),
	}
	if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
		v4Warn("no se pudo crear el evento de hechos: %v", err)
	}
}

func (s *formService) v4ProcessIdent(ctx context.Context, fu *models.FollowUpV2, actorID string, loaded *v4Loaded) int {
	roots := v4Roots(loaded.byGroup[v4Key("repeater:identEntidades")])
	kids := loaded.byGroup[v4Key("repeater:identBarreras")]
	acuerdo := v4Key("q:ident.acuerdo")
	debida := v4Key("q:ident.debida")
	count := 0
	for _, entry := range roots {
		branch, hasBranch := v4ParseBranch(entry.answers[v4Key("q:ident.entidad")])
		var branchPtr *int64
		entityCaseID := ""
		if hasBranch {
			b := branch
			branchPtr = &b
			var objetivo, radicado *string
			if v := strings.TrimSpace(entry.answers[v4Key("q:ident.otraRuta")]); v != "" {
				objetivo = &v
			}
			if v := strings.TrimSpace(entry.answers[v4Key("q:ident.radicado")]); v != "" {
				radicado = &v
			}
			entityCaseID = s.v4EnsureEntityCase(ctx, fu.CaseID, branch, objetivo, radicado, actorID)
			if entityCaseID != "" {
				s.v4AddCatalogObligations(ctx, entityCaseID, fu.ID, actorID, entry.answers[v4Key("q:ident.completadas")], models.EntityCaseObligationStatusCompletada)
				s.v4AddCatalogObligations(ctx, entityCaseID, fu.ID, actorID, entry.answers[v4Key("q:ident.pendientes")], models.EntityCaseObligationStatusPendiente)
				if v4BothYes(entry.answers, acuerdo, debida) {
					s.v4ProcesarCanales(ctx, fu.CaseID, &fu.ID, actorID, entityCaseID, branch, entry.answers[v4Key("q:ident.canales")])
				}
			}
		}
		if entry.answers[v4Key("q:ident.nuevaBarrera")] == "true" {
			count += s.v4NewBarriers(ctx, fu, actorID, branchPtr, v4Kids(kids, entry.entry.ID), v4IdentBarrier)
		}
	}
	return count
}

func (s *formService) v4ProcessSeg(ctx context.Context, fu *models.FollowUpV2, actorID string, loaded *v4Loaded) (int, error) {
	roots := v4Roots(loaded.byGroup[v4Key("repeater:segEntidades")])
	if len(roots) == 0 {
		return 0, nil
	}
	if s.entityCaseSvc == nil {
		v4Warn("entityCaseSvc es nil — no se procesa el seguimiento a entidades")
		return 0, nil
	}
	entidades, err := s.entityCaseSvc.ListByCase(ctx, fu.CaseID)
	if err != nil {
		return 0, err
	}
	active := s.v4OrderedActiveBarriers(ctx, fu)
	newKids := loaded.byGroup[v4Key("repeater:segBarreras")]
	activeKids := loaded.byGroup[v4Key("repeater:segBarrerasActivas")]
	acuerdo := v4Key("q:seg.acuerdo")
	debida := v4Key("q:seg.debida")
	count := 0
	for i, entry := range roots {
		if i >= len(entidades) {
			v4Warn("entrada de seguimiento a entidades %d sin entidad del caso", i)
			continue
		}
		ec := entidades[i]
		branch := ec.EntityBranchID
		s.v4SeguimientoEntidad(ctx, fu, actorID, ec.RelID, branch, entry.answers, acuerdo, debida)
		if entry.answers[v4Key("q:seg.nuevaBarrera")] == "true" {
			count += s.v4NewBarriers(ctx, fu, actorID, &branch, v4Kids(newKids, entry.entry.ID), v4SegBarrier)
		}
		s.v4ActiveBarriers(ctx, fu, actorID, v4Kids(activeKids, entry.entry.ID), v4BarriersForBranch(active, branch))
	}
	return count, nil
}

func (s *formService) v4SeguimientoEntidad(ctx context.Context, fu *models.FollowUpV2, actorID, entityCaseID string, branch int64, answers map[string]string, acuerdoKey, debidaKey string) {
	texto := strings.TrimSpace(answers[v4Key("q:seg.actualizacion")])
	activa := v4BothYes(answers, acuerdoKey, debidaKey)
	canales := ""
	if activa {
		canales = strings.TrimSpace(answers[v4Key("q:seg.canales")])
	}
	if s.entityCaseFollowUpRepo != nil {
		row := &models.EntityCaseFollowUp{
			EntityCaseID:            entityCaseID,
			FollowUpID:              fu.ID,
			RutaActualizada:         texto != "",
			RequiereNuevaActivacion: activa,
			CreatedByID:             actorID,
		}
		if texto != "" {
			row.MotivoActualizacion = &texto
		}
		if canales != "" {
			row.CanalActivacion = &canales
		}
		if err := s.entityCaseFollowUpRepo.Create(ctx, row); err != nil {
			v4Warn("no se pudo crear el seguimiento de la entidad %s: %v", entityCaseID, err)
		}
	}
	resumen := "Seguimiento registrado sin cambios en la ruta"
	if texto != "" {
		resumen = "Ruta actualizada: " + texto
	}
	if s.entityCaseSvc != nil {
		if err := s.entityCaseSvc.UpdateLastAction(ctx, entityCaseID, resumen); err != nil {
			v4Warn("no se pudo actualizar last_action de %s: %v", entityCaseID, err)
		}
	}
	if s.entityCaseObligationRepo != nil {
		for _, id := range splitCSVValues(answers[v4Key("q:seg.completadas")]) {
			if id == "" || id == "otra" {
				continue
			}
			if err := s.entityCaseObligationRepo.UpdateFields(ctx, id, map[string]interface{}{
				"status": models.EntityCaseObligationStatusCompletada,
			}); err != nil {
				v4Warn("no se pudo completar la obligación %s: %v", id, err)
			}
		}
	}
	if v4Has(answers[v4Key("q:seg.pendientes")], "otra") {
		label := strings.TrimSpace(answers[v4Key("q:seg.otraPendiente")])
		if label == "" {
			label = "Otra"
		}
		s.v4AddCustomObligation(ctx, entityCaseID, fu.ID, actorID, label, models.EntityCaseObligationStatusPendiente)
	}
	if activa {
		s.v4ProcesarCanales(ctx, fu.CaseID, &fu.ID, actorID, entityCaseID, branch, canales)
	}
}

func (s *formService) v4EnsureEntityCase(ctx context.Context, caseID string, branch int64, objetivo, radicado *string, actorID string) string {
	if s.entityCaseSvc == nil {
		v4Warn("entityCaseSvc es nil — no se crea la entidad %d", branch)
		return ""
	}
	created, err := s.entityCaseSvc.Create(ctx, CreateEntityCaseInput{
		CaseID: caseID, EntityBranchID: branch,
		Objetivo: objetivo, Radicado: radicado, CreatedByID: actorID,
	})
	if err == nil {
		return created.ID
	}
	if !errors.Is(err, ErrEntityCaseDuplicate) {
		v4Warn("no se pudo crear entity_case (sede %d): %v", branch, err)
		return ""
	}
	existentes, findErr := s.entityCaseSvc.ListByCase(ctx, caseID)
	if findErr != nil {
		v4Warn("no se pudo resolver la entidad ya asociada (sede %d): %v", branch, findErr)
		return ""
	}
	for _, e := range existentes {
		if e.EntityBranchID == branch {
			return e.RelID
		}
	}
	v4Warn("entidad duplicada no encontrada (sede %d)", branch)
	return ""
}

func (s *formService) v4AddCatalogObligations(ctx context.Context, entityCaseID, followUpID, actorID, raw, status string) {
	if s.entityCaseObligationRepo == nil || entityCaseID == "" {
		return
	}
	for _, val := range splitCSVValues(raw) {
		eco := &models.EntityCaseObligation{
			EntityCaseID: entityCaseID,
			Status:       status,
			FollowUpID:   &followUpID,
			CreatedByID:  actorID,
		}
		if val == "otra" {
			label := "Otra"
			eco.CustomLabel = &label
		} else {
			id := val
			eco.EntityObligationID = &id
		}
		if err := s.createEntityCaseObligationIfNotExists(ctx, eco); err != nil {
			v4Warn("no se pudo crear la obligación %s de %s: %v", status, entityCaseID, err)
		}
	}
}

func (s *formService) v4AddCustomObligation(ctx context.Context, entityCaseID, followUpID, actorID, label, status string) {
	if s.entityCaseObligationRepo == nil || entityCaseID == "" {
		return
	}
	eco := &models.EntityCaseObligation{
		EntityCaseID: entityCaseID,
		Status:       status,
		FollowUpID:   &followUpID,
		CreatedByID:  actorID,
		CustomLabel:  &label,
	}
	if err := s.entityCaseObligationRepo.Create(ctx, eco); err != nil {
		v4Warn("no se pudo crear la obligación libre de %s: %v", entityCaseID, err)
	}
}

func (s *formService) v4ProcesarCanales(ctx context.Context, caseID string, followUpID *string, actorID, entityCaseID string, branch int64, canalRaw string) {
	if entityCaseID == "" || branch <= 0 {
		return
	}
	for _, canal := range splitCSVValues(canalRaw) {
		switch canal {
		case "ninguna", "":
			continue
		case "notificacion":
			var letterID *string
			if s.entityLetterRepo != nil {
				letter := &models.EntityLetter{
					CaseID:         caseID,
					State:          models.EntityLetterStatePorProyectar,
					AgentID:        &actorID,
					EntityBranchID: &branch,
				}
				if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
					v4Warn("no se pudo crear el oficio de %s: %v", entityCaseID, err)
				} else {
					letterID = &letter.ID
				}
			}
			s.v4Task(ctx, &models.CaseTask{
				Category: "Entidades", Type: "proyectar_oficio",
				Description:    "Proyectar oficio de activación de ruta",
				AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
				CaseID: caseID, FollowUpID: followUpID,
				EntityCaseID: &entityCaseID, EntityLetterID: letterID,
			})
		case "misalvia":
			s.v4Task(ctx, &models.CaseTask{
				Category: "Entidades", Type: "solicitud_entidad",
				Description:    "Hacer solicitudes a la entidad sobre obligaciones pendientes",
				AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
				CaseID: caseID, FollowUpID: followUpID, EntityCaseID: &entityCaseID,
			})
		case "articulacion":
			s.v4GestionLlamada(ctx, caseID, followUpID, actorID, entityCaseID, "Articulación institucional telefónica o por mensajería")
		case "denuncia":
			s.v4GestionLlamada(ctx, caseID, followUpID, actorID, entityCaseID, "Denuncia fácil")
		default:
			v4Warn("canal sin efecto: %s", canal)
		}
	}
}

func (s *formService) v4GestionLlamada(ctx context.Context, caseID string, followUpID *string, actorID, entityCaseID, description string) {
	s.v4Task(ctx, &models.CaseTask{
		Category: "Entidades", Type: "gestion_llamada",
		Description: description, AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
		CaseID: caseID, FollowUpID: followUpID, EntityCaseID: &entityCaseID,
	})
}

func (s *formService) v4Task(ctx context.Context, task *models.CaseTask) {
	if s.caseTaskRepo == nil {
		return
	}
	if err := s.caseTaskRepo.Create(ctx, task); err != nil {
		v4Warn("no se pudo crear la tarea %s: %v", task.Type, err)
	}
}

func (s *formService) v4NewBarriers(ctx context.Context, fu *models.FollowUpV2, actorID string, branch *int64, kids []v4Entry, q v4BarrierQ) int {
	if s.barrierV2Repo == nil {
		if len(kids) > 0 {
			v4Warn("barrierV2Repo es nil — no se crean barreras nuevas")
		}
		return 0
	}
	count := 0
	for _, kid := range kids {
		m := kid.answers
		sector := strings.TrimSpace(m[q.sector])
		b := &models.BarrierV2{
			CaseID: fu.CaseID, FollowUpID: &fu.ID, CreatedByID: actorID,
			Status: models.BarrierV2StatusOpen, Sector: sector,
			StructuralInstitutional: m[q.talento],
			StructuralEconomic:      m[q.eco],
			StructuralTerritorial:   m[q.terr],
			StructuralDifferential:  m[q.diferencial],
			BarrierDate:             m[q.fecha],
			OfficialDependency:      m[q.funcionario],
			Description:             m[q.descripcion],
			ManagementActions:       m[q.gestion],
			DepartmentID:            strings.TrimSpace(m[q.depto]),
			CityID:                  strings.TrimSpace(m[q.municipio]),
			EntityBranchID:          branch,
		}
		switch sector {
		case "salud":
			b.SpecificBarriers, b.SpecificInstitutions, b.OtherBarrierDesc = m[q.tiposSalud], m[q.instSalud], m[q.otraSalud]
		case "justicia":
			b.SpecificBarriers, b.SpecificInstitutions, b.OtherBarrierDesc = m[q.tiposJusticia], m[q.instJusticia], m[q.otraJusticia]
		case "proteccion":
			b.SpecificBarriers, b.SpecificInstitutions, b.OtherBarrierDesc = m[q.tiposProteccion], m[q.instProteccion], m[q.otraProteccion]
		}
		if err := s.barrierV2Repo.Create(ctx, b); err != nil {
			v4Warn("no se pudo crear la barrera: %v", err)
			continue
		}
		count++
		s.v4GestionBarreraNueva(ctx, fu, actorID, b.ID, branch, m[q.gestion])
	}
	return count
}

func (s *formService) v4GestionBarreraNueva(ctx context.Context, fu *models.FollowUpV2, actorID, barrierID string, branch *int64, gestionCSV string) {
	for _, g := range splitCSVValues(gestionCSV) {
		switch g {
		case "activacion":
			var letterID *string
			if s.entityLetterRepo != nil {
				letter := &models.EntityLetter{
					BarrierID: barrierID, CaseID: fu.CaseID,
					State: models.EntityLetterStatePorProyectar, AgentID: &actorID,
					EntityBranchID: branch,
				}
				if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
					v4Warn("no se pudo crear el oficio de la barrera %s: %v", barrierID, err)
				} else {
					letterID = &letter.ID
				}
			}
			id := barrierID
			s.v4Task(ctx, &models.CaseTask{
				Category: "Barreras", Type: "proyectar_oficio",
				Description:    "Activación de ruta interinstitucional (Generar oficio)",
				AssignedUserID: actorID, Status: models.CaseTaskStatusToDo,
				CaseID: fu.CaseID, FollowUpID: &fu.ID,
				BarrierID: &id, EntityLetterID: letterID,
			})
		case "alerta":
			s.v4AlertaComite(ctx, fu, barrierID)
		}
	}
}

func (s *formService) v4AlertaComite(ctx context.Context, fu *models.FollowUpV2, barrierID string) {
	id := barrierID
	s.v4Task(ctx, &models.CaseTask{
		Category: "Barreras", Type: "comite_caso",
		Description: "Alerta por barreras", AssignedUserID: "",
		Status: models.CaseTaskStatusToDo, CaseID: fu.CaseID, FollowUpID: &fu.ID,
		BarrierID: &id,
	})
}

func (s *formService) v4OrderedActiveBarriers(ctx context.Context, fu *models.FollowUpV2) []models.BarrierV2 {
	if s.barrierV2Repo == nil {
		return nil
	}
	if fu.ActiveBarrierIDs != nil && strings.TrimSpace(*fu.ActiveBarrierIDs) != "" {
		ids := splitCSVValues(*fu.ActiveBarrierIDs)
		found, err := s.barrierV2Repo.FindByIDs(ctx, ids)
		if err != nil {
			v4Warn("no se pudieron leer las barreras activas: %v", err)
			return nil
		}
		byID := make(map[string]models.BarrierV2, len(found))
		for _, b := range found {
			byID[b.ID] = b
		}
		ordered := make([]models.BarrierV2, 0, len(ids))
		for _, id := range ids {
			if b, ok := byID[id]; ok {
				ordered = append(ordered, b)
			}
		}
		return ordered
	}
	found, err := s.barrierV2Repo.FindActiveByCaseID(ctx, fu.CaseID)
	if err != nil {
		v4Warn("no se pudieron leer las barreras activas del caso: %v", err)
		return nil
	}
	return found
}

func v4BarriersForBranch(all []models.BarrierV2, branch int64) []models.BarrierV2 {
	var out []models.BarrierV2
	for _, b := range all {
		if b.EntityBranchID != nil && *b.EntityBranchID == branch {
			out = append(out, b)
		}
	}
	return out
}

func (s *formService) v4ActiveBarriers(ctx context.Context, fu *models.FollowUpV2, actorID string, kids []v4Entry, barriers []models.BarrierV2) {
	qPersiste := v4Key("q:seg.barreraActiva.persisten")
	qCierre := v4Key("q:seg.barreraActiva.cierre")
	qMotivo := v4Key("q:seg.barreraActiva.motivoCierre")
	qGestion := v4Key("q:seg.barreraActiva.gestion")
	for i, kid := range kids {
		if i >= len(barriers) {
			v4Warn("seguimiento de barrera sin barrera activa en la posición %d", i)
			continue
		}
		barrierID := barriers[i].ID
		m := kid.answers
		closes := m[qCierre] == "true"
		if closes && s.barrierV2Repo != nil {
			if err := s.barrierV2Repo.UpdateStatus(ctx, barrierID, models.BarrierV2StatusManaged); err != nil {
				v4Warn("no se pudo cerrar la barrera %s: %v", barrierID, err)
			}
		}
		gestion := strings.TrimSpace(m[qGestion])
		if gestion == "alerta" {
			s.v4AlertaComite(ctx, fu, barrierID)
		}
		if s.caseTimelineRepo != nil {
			persiste := "No persiste"
			if m[qPersiste] == "true" {
				persiste = "Persiste"
			}
			desc := persiste
			if gestion != "" {
				desc += " · Gestión: " + gestion
			}
			event := &models.CaseTimelineEvent{
				CaseID: fu.CaseID, FollowUpID: fu.ID, Category: "Barreras",
				Type: "Seguimiento a Barrera", Icon: "shield-halved", Color: "#6366f1",
				Description: desc, EventUserID: actorID, Date: time.Now(),
			}
			if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
				v4Warn("no se pudo crear el evento de la barrera %s: %v", barrierID, err)
			}
		}
		if s.barrierFollowUpRepo != nil {
			row := &models.BarrierFollowUp{
				BarrierID: barrierID, FollowUpID: &fu.ID, CreatedByID: actorID,
				Persists: m[qPersiste] == "true", ManagementActions: gestion,
				ClosesBarrier: closes, ClosureReason: m[qMotivo],
			}
			if err := s.barrierFollowUpRepo.Create(ctx, row); err != nil {
				v4Warn("no se pudo crear el seguimiento de la barrera %s: %v", barrierID, err)
			}
		}
	}
}

func (s *formService) v4ProcessPlan(ctx context.Context, fu *models.FollowUpV2, actorID string, direct map[string]string) int {
	if direct[v4Key("q:plan.remision")] != "true" {
		return 0
	}
	equipos := direct[v4Key("q:plan.equipos")]
	tieneMedidas := v4Has(equipos, "medidas")
	count := 0
	if tieneMedidas {
		tipos := splitCSVValues(direct[v4Key("q:plan.medidasTipos")])
		if len(tipos) > 0 {
			if err := s.createEmergencyMeasureReferral(ctx, fu.CaseID, &fu.ID, actorID, tipos); err != nil {
				v4Warn("no se pudo crear la remisión de medidas: %v", err)
			} else {
				count += len(tipos)
			}
		}
	}
	for _, equipo := range splitCSVValues(equipos) {
		switch equipo {
		case "psico":
			if tieneMedidas || !v4PsicoAplica(direct[v4Key("q:plan.psico")]) {
				continue
			}
			if s.v4CrearPsico(ctx, fu.CaseID, &fu.ID, actorID, "derivacion", direct[v4Key("q:plan.psico")], "Validar que remision a psicosocial es valida") {
				count++
			}
		case "osignh":
			if tieneMedidas {
				continue
			}
			if s.v4CrearPsico(ctx, fu.CaseID, &fu.ID, actorID, "osignh", direct[v4Key("q:plan.osignh")], "Validar que remision a psicosocial OSIGNH es valida") {
				count++
			}
		case "hombres":
			if s.menTeamRemisionRepo == nil {
				v4Warn("menTeamRemisionRepo es nil — no se crea la remisión de hombres")
				continue
			}
			row := &models.MenTeamRemision{CaseID: fu.CaseID, FollowUpID: &fu.ID, Type: "derivacion", Status: "ACTIVE"}
			if err := s.menTeamRemisionRepo.Create(ctx, row); err != nil {
				v4Warn("no se pudo crear la remisión de hombres: %v", err)
				continue
			}
			count++
		case "discapacidad":
			if s.discapacidadRemisionRepo == nil {
				v4Warn("discapacidadRemisionRepo es nil — no se crea la remisión")
				continue
			}
			for _, svc := range splitCSVValues(direct[v4Key("q:plan.discapacidad")]) {
				row := &models.DiscapacidadRemision{CaseID: fu.CaseID, FollowUpID: &fu.ID, Service: svc, Status: "ACTIVE"}
				if err := s.discapacidadRemisionRepo.Create(ctx, row); err != nil {
					v4Warn("no se pudo crear la remisión de discapacidad %s: %v", svc, err)
					continue
				}
				count++
			}
		case "estabilizacion":
			if len(splitCSVValues(direct[v4Key("q:plan.estabilizacion")])) == 0 {
				continue
			}
			if s.economicStabilizationRepo == nil {
				v4Warn("economicStabilizationRepo es nil — no se crea la remisión")
				continue
			}
			row := &models.EconomicStabilization{CaseID: fu.CaseID, FollowUpID: &fu.ID, Type: "derivacion", Status: "ACTIVE"}
			if err := s.economicStabilizationRepo.Create(ctx, row); err != nil {
				v4Warn("no se pudo crear la remisión de estabilización: %v", err)
				continue
			}
			count++
		}
	}
	return count
}

func (s *formService) v4CrearPsico(ctx context.Context, caseID string, followUpID *string, actorID, tipo, criterios, taskDesc string) bool {
	if s.psychosocialSupportRepo == nil {
		v4Warn("psychosocialSupportRepo es nil — no se crea la remisión %s", tipo)
		return false
	}
	var submittedBy, submittedByTeam *string
	if actorID != "" {
		actor := actorID
		submittedBy = &actor
	}
	if s.agentLightRepo != nil && actorID != "" {
		agent, err := s.agentLightRepo.FindByICode(ctx, actorID)
		if err == nil && agent != nil && strings.TrimSpace(agent.Team) != "" {
			team := strings.TrimSpace(agent.Team)
			submittedByTeam = &team
		}
	}
	var criteria *string
	if v := strings.TrimSpace(criterios); v != "" {
		criteria = &v
	}
	ps := &models.PsychosocialSupport{
		CaseID: caseID, FollowUpID: followUpID, Type: tipo,
		Status:      models.PsychosocialSupportStatusAbierto,
		SubmittedBy: submittedBy, SubmittedByTeam: submittedByTeam,
		ReferralCriteria: criteria,
	}
	if err := s.psychosocialSupportRepo.Create(ctx, ps); err != nil {
		v4Warn("no se pudo crear la remisión %s: %v", tipo, err)
		return false
	}
	psID := ps.ID
	s.v4Task(ctx, &models.CaseTask{
		Category: "Psicosocial", Type: "validar_remision", Description: taskDesc,
		AssignedUserID: "", Status: models.CaseTaskStatusToDo,
		CaseID: caseID, FollowUpID: followUpID, PsychosocialSupportID: &psID,
	})
	return true
}

func (s *formService) v4ProcessConclusion(ctx context.Context, fu *models.FollowUpV2, actorID string, direct map[string]string) {
	if direct[v4Key("q:con.concluye")] == "true" {
		if s.casoCierreService == nil {
			v4Warn("casoCierreService es nil — el caso %s no se cierra", fu.CaseID)
		} else if err := s.casoCierreService.CerrarCaso(ctx, CerrarCasoInput{
			CaseICode:               fu.CaseID,
			Motivo:                  direct[v4Key("q:con.motivo")],
			OtroMotivo:              direct[v4Key("q:con.otro")],
			Descripcion:             direct[v4Key("q:con.descripcion")],
			AccionesInstitucionales: direct[v4Key("q:con.accion")] == "true",
			ActorID:                 actorID,
		}); err != nil {
			v4Warn("no se pudo cerrar el caso %s: %v", fu.CaseID, err)
		}
	}
	if direct[v4Key("q:con.accion")] != "true" {
		return
	}
	branch, ok := v4ParseBranch(direct[v4Key("q:con.entidad")])
	if !ok {
		return
	}
	entityCaseID := s.v4EnsureEntityCase(ctx, fu.CaseID, branch, nil, nil, actorID)
	s.v4ProcesarCanales(ctx, fu.CaseID, &fu.ID, actorID, entityCaseID, branch, direct[v4Key("q:con.canales")])
}

func (s *formService) applyV4Reasignacion(ctx context.Context, fu *models.FollowUpV2, actorID string, dest v4Reasignacion) {
	if !dest.apply || s.caseRepo == nil {
		return
	}
	if err := s.caseRepo.UpdateRiskLevelByICode(ctx, fu.CaseID, dest.level); err != nil {
		v4Warn("no se pudo actualizar el nivel de riesgo: %v", err)
		return
	}
	if s.followUpV2Svc == nil {
		v4Warn("followUpV2Svc es nil — no se reasigna el calendario")
		return
	}
	if err := s.followUpV2Svc.ReasignarCalendario(ctx, fu.CaseID, dest.level); err != nil {
		v4Warn("ReasignarCalendario falló: %v", err)
		return
	}
	newAgentID := ""
	if s.followUpRepo != nil {
		if pendientes, err := s.followUpRepo.FindPendingByCaseID(ctx, fu.CaseID); err == nil && len(pendientes) > 0 && pendientes[0].AgentID != nil {
			newAgentID = *pendientes[0].AgentID
		}
	}
	if err := s.caseRepo.UpdateTeamAndAgent(ctx, fu.CaseID, dest.team, newAgentID); err != nil {
		v4Warn("no se pudo actualizar equipo y agente: %v", err)
	}
	if s.caseTimelineRepo == nil {
		return
	}
	now := time.Now()
	event := &models.CaseTimelineEvent{
		CaseID: fu.CaseID, FollowUpID: fu.ID, Category: models.TimelineCategoryGeneral,
		Type: models.TimelineTypeReasignacionCaso, Icon: "arrows-rotate", Color: "#f59e0b",
		Date: now, EventUserID: actorID, CreatedAt: now,
	}
	if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
		v4Warn("no se pudo crear el evento de reasignación: %v", err)
	}
}
