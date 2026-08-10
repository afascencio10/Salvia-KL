// Package service — case_detail_service.go
// Lógica de negocio para la pantalla de detalle de caso (rol sv).
package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	internaldb "bitsflow/internal/db"
	salvia_config "bitsflow/salvia/config"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var ErrCaseNotFound = errors.New("caso no encontrado")
var ErrInvalidICode = errors.New("icode inválido")

type CaseDetailService interface {
	GetDetail(ctx context.Context, caseICode string) (*repository.CaseDetailData, error)
	CreateFollowUp(ctx context.Context, caseICode, agentID, scheduledDate, notas, createdBy string) (*models.FollowUpV2, error)
	AddTimelineEvent(ctx context.Context, caseICode, eventType, description, actorID, actorName string) error
	ReassignFollowUp(ctx context.Context, followUpID, newAgentID string) error
	EditFollowUpDate(ctx context.Context, followUpID, scheduledDate, actorName string) error
	ReasignarCaso(ctx context.Context, caseICode, newOperadorICode string) error
	GetDB() *gorm.DB
}

type caseDetailService struct {
	repo repository.CaseDetailRepository
	db   *gorm.DB
}

func NewCaseDetailService(repo repository.CaseDetailRepository, db *gorm.DB) CaseDetailService {
	return &caseDetailService{repo: repo, db: db}
}

func (s *caseDetailService) GetDB() *gorm.DB {
	return s.db
}

func (s *caseDetailService) GetDetail(ctx context.Context, caseICode string) (*repository.CaseDetailData, error) {
	if caseICode == "" {
		return nil, ErrInvalidICode
	}
	var detail *repository.CaseDetailData
	var err error
	err = internaldb.WithRetry(func() error {
		detail, err = s.repo.GetByICode(ctx, caseICode)
		return err
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCaseNotFound
		}
		return nil, err
	}

	// Ajustar status de barreras: si tiene tareas pendientes y está en OPEN → "En Gestion"
	for i := range detail.Barriers {
		if detail.Barriers[i].Status == models.BarrierV2StatusOpen {
			var pendingCount int64
			s.db.WithContext(ctx).Raw(`
				SELECT COUNT(*) FROM salvia.case_task
				WHERE barrier_id = ? AND status = 'ToDo' AND deleted_at IS NULL
			`, detail.Barriers[i].ID).Scan(&pendingCount)
			if pendingCount > 0 {
				detail.Barriers[i].Status = models.BarrierV2StatusEnGestion
				// También actualizar en BD para que quede consistente
				s.db.Exec(`UPDATE salvia.barrier_v2 SET status = ? WHERE id = ?`, models.BarrierV2StatusEnGestion, detail.Barriers[i].ID)
			}
		}
	}

	// Traducir enums con Locale
	locale := salvia_config.Locale["sp"]
	for i, v := range detail.TipoViolencia {
		if t, ok := locale[v]; ok {
			detail.TipoViolencia[i] = t
		}
	}
	for i, v := range detail.SubtipoViolencia {
		if t, ok := locale[v]; ok {
			detail.SubtipoViolencia[i] = t
		}
	}
	for i, v := range detail.AmbitoViolencia {
		if t, ok := locale[v]; ok {
			detail.AmbitoViolencia[i] = t
		}
	}
	if t, ok := locale[detail.Genero]; ok {
		detail.Genero = t
	}
	if t, ok := locale[detail.Nacionalidad]; ok {
		detail.Nacionalidad = t
	}
	if t, ok := locale[detail.TipoAgresorResumen]; ok {
		detail.TipoAgresorResumen = t
	}
	for i, v := range detail.PlanAtencion {
		if t, ok := locale[v]; ok {
			detail.PlanAtencion[i] = t
		}
	}
	for i, v := range detail.AjusteRazonable {
		if t, ok := locale[v]; ok {
			detail.AjusteRazonable[i] = t
		}
	}

	return detail, nil
}

func (s *caseDetailService) CreateFollowUp(ctx context.Context, caseICode, agentID, scheduledDate, notas, createdBy string) (*models.FollowUpV2, error) {
	if caseICode == "" {
		return nil, ErrInvalidICode
	}

	// Parsear en zona horaria de Colombia para evitar desfase de fecha
	loc, _ := time.LoadLocation("America/Bogota")
	fecha, err := time.ParseInLocation("2006-01-02T15:04", scheduledDate, loc)
	if err != nil {
		fecha, err = time.ParseInLocation("2006-01-02", scheduledDate, loc)
		if err != nil {
			return nil, errors.New("formato de fecha inválido, use YYYY-MM-DD o YYYY-MM-DDTHH:MM")
		}
	}

	// Validar que la fecha no sea del pasado (en zona horaria de Colombia)
	hoy := time.Now().In(loc).Truncate(24 * time.Hour)
	if fecha.Before(hoy) {
		return nil, errors.New("no se permiten seguimientos con fecha anterior a hoy")
	}

	// Si es hoy y tiene hora, validar que la hora no haya pasado
	ahora := time.Now().In(loc)
	if len(scheduledDate) > 10 && fecha.Year() == ahora.Year() && fecha.Month() == ahora.Month() && fecha.Day() == ahora.Day() {
		if fecha.Before(ahora) {
			return nil, errors.New("la hora programada ya pasó, seleccione una hora futura")
		}
	}

	// Extraer la hora como string para el campo scheduled_time
	// Solo guardar hora si el usuario la envió (formato con T indica que tiene hora)
	hora := ""
	if len(scheduledDate) > 10 {
		hora = fecha.Format("15:04")
	}

	// Calcular sequence_number: contar los existentes + 1 (sin límite máximo)
	count, err := s.repo.CountFollowUpsByCaseID(ctx, caseICode)
	if err != nil {
		return nil, err
	}

	nextSeq := count + 1

	sinEvaluar := "SIN_EVALUAR"
	followUp := &models.FollowUpV2{
		CaseID:         caseICode,
		AgentID:        &agentID,
		Status:         models.FollowUpStatusPendiente,
		ScheduledDate:  fecha,
		ScheduledTime:  hora,
		Summary:        &notas,
		Team:           "SIN_EQUIPO",
		RiskStatus:     &sinEvaluar,
		SequenceNumber: nextSeq,
	}

	if err := s.repo.CreateFollowUpV2(ctx, followUp); err != nil {
		return nil, err
	}

	// Obtener nombre del agente asignado para el timeline
	var agentName string
	s.db.Raw(`SELECT gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names
		FROM security.general_user gu
		JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile
		WHERE gu.general_user_i_code = ?`, agentID).Scan(&agentName)
	if agentName == "" {
		agentName = agentID
	}

	// Registrar evento en el timeline con información completa
	descripcion := "Seguimiento #" + fmt.Sprintf("%d", nextSeq) + " creado — Asignado a: " + agentName + " — Fecha: " + fecha.Format("2006-01-02")
	if hora != "" {
		descripcion += " " + hora
	}
	actorTimeline := createdBy
	if actorTimeline == "" {
		actorTimeline = agentName
	}
	now := time.Now()
	s.repo.CreateTimelineEvent(ctx, &models.CaseTimelineEvent{
		CaseID:      caseICode,
		EventType:   models.TimelineEventSeguimiento,
		Category:    models.TimelineCategorySeguimientos,
		Type:        models.TimelineTypeSeguimientoProgramado,
		Icon:        models.TimelineIconSeguimiento,
		Color:       models.TimelineColorGreen,
		Date:        now,
		Description: descripcion,
		ActorName:   actorTimeline,
		EventUserID: agentID,
		FollowUpID:  followUp.ID,
		CreatedAt:   now,
	})

	return followUp, nil
}

func (s *caseDetailService) AddTimelineEvent(ctx context.Context, caseICode, eventType, description, actorID, actorName string) error {
	// Mapear event_type legacy a los campos nuevos (Category, Type, Icon, Color)
	category := models.TimelineCategoryGeneral
	tlType := eventType
	icon := models.TimelineIconNota
	color := models.TimelineColorGray

	switch eventType {
	case models.TimelineEventReasignacion:
		category = models.TimelineCategoryGeneral
		tlType = models.TimelineTypeReasignacionCaso
		icon = models.TimelineIconReasignacion
		color = models.TimelineColorBlue
	case models.TimelineEventReasignSeg:
		category = models.TimelineCategorySeguimientos
		tlType = models.TimelineTypeReasignacionSeg
		icon = models.TimelineIconReasignacion
		color = models.TimelineColorOrange
	case models.TimelineEventEstadoCambio:
		category = models.TimelineCategoryGeneral
		tlType = models.TimelineTypeCambioEstado
		icon = models.TimelineIconEstado
		color = models.TimelineColorOrange
	case models.TimelineEventSeguimiento:
		category = models.TimelineCategorySeguimientos
		tlType = models.TimelineTypeSeguimientoProgramado
		icon = models.TimelineIconSeguimiento
		color = models.TimelineColorGreen
	case models.TimelineEventBarrera:
		category = models.TimelineCategoryBarreras
		tlType = models.TimelineTypeBarreraIdentificada
		icon = models.TimelineIconBarrera
		color = models.TimelineColorOrange
	case models.TimelineEventNota:
		category = models.TimelineCategoryGeneral
		tlType = models.TimelineTypeNota
		icon = models.TimelineIconNota
		color = models.TimelineColorGray
	}

	now := time.Now()
	return s.repo.CreateTimelineEvent(ctx, &models.CaseTimelineEvent{
		CaseID:      caseICode,
		EventType:   eventType,
		Category:    category,
		Type:        tlType,
		Icon:        icon,
		Color:       color,
		Date:        now,
		Description: description,
		ActorID:     actorID,
		ActorName:   actorName,
		EventUserID: actorID,
		CreatedAt:   now,
	})
}

func (s *caseDetailService) ReassignFollowUp(ctx context.Context, followUpID, newAgentID string) error {
	return s.repo.UpdateFollowUpAgent(ctx, followUpID, newAgentID)
}

func (s *caseDetailService) EditFollowUpDate(ctx context.Context, followUpID, scheduledDate, actorName string) error {
	// Parsear en zona horaria de Colombia
	loc, _ := time.LoadLocation("America/Bogota")
	fecha, err := time.ParseInLocation("2006-01-02T15:04", scheduledDate, loc)
	if err != nil {
		fecha, err = time.ParseInLocation("2006-01-02", scheduledDate, loc)
		if err != nil {
			return errors.New("formato de fecha inválido")
		}
	}

	// Validar que no sea fecha pasada
	hoy := time.Now().In(loc).Truncate(24 * time.Hour)
	if fecha.Before(hoy) {
		return errors.New("no se permiten fechas anteriores a hoy")
	}

	// Extraer hora solo si se envió
	hora := ""
	if len(scheduledDate) > 10 {
		hora = fecha.Format("15:04")
	}

	// Actualizar en BD
	fields := map[string]interface{}{
		"scheduled_date": fecha,
	}
	if hora != "" {
		fields["scheduled_time"] = hora
	} else {
		fields["scheduled_time"] = ""
	}

	if err := s.db.Model(&models.FollowUpV2{}).Where("id = ?", followUpID).Updates(fields).Error; err != nil {
		return err
	}

	// Obtener el case_id del seguimiento para registrar en el timeline
	var fu models.FollowUpV2
	if err := s.db.Where("id = ?", followUpID).First(&fu).Error; err == nil {
		now := time.Now()
		descripcion := "Seguimiento editado — Nueva fecha: " + fecha.Format("2006-01-02")
		if hora != "" {
			descripcion += " " + hora
		}
		s.repo.CreateTimelineEvent(ctx, &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			EventType:   models.TimelineEventSeguimiento,
			Category:    models.TimelineCategorySeguimientos,
			Type:        models.TimelineTypeSeguimientoEditado,
			Icon:        models.TimelineIconPospuesto,
			Color:       models.TimelineColorBlue,
			Date:        now,
			Description: descripcion,
			ActorName:   actorName,
			FollowUpID:  followUpID,
			CreatedAt:   now,
		})
	}

	return nil
}

func (s *caseDetailService) ReasignarCaso(ctx context.Context, caseICode, newOperadorICode string) error {
	return internaldb.WithRetry(func() error {
		return s.reasignarCasoInternal(ctx, caseICode, newOperadorICode)
	})
}

func (s *caseDetailService) reasignarCasoInternal(ctx context.Context, caseICode, newOperadorICode string) error {
	db := s.db

	// 1. Verificar que el caso existe
	var victimCaseID int64
	if err := db.Raw("SELECT victim_case_id FROM salvia.victim_case WHERE victim_case_i_code = ?", caseICode).Scan(&victimCaseID).Error; err != nil {
		return errors.New("caso no encontrado")
	}
	if victimCaseID == 0 {
		return errors.New("caso no encontrado")
	}

	// 2. Obtener el team del nuevo operador
	var newTeam string
	db.Raw(`SELECT COALESCE(general_user_team, '') FROM security.general_user WHERE general_user_i_code = ?`, newOperadorICode).Scan(&newTeam)

	// 3. Actualizar agent_id y victim_case_team directamente en victim_case
	if newTeam != "" {
		db.Exec(`UPDATE salvia.victim_case SET agent_id = ?, victim_case_team = ? WHERE victim_case_i_code = ?`, newOperadorICode, newTeam, caseICode)
	} else {
		db.Exec(`UPDATE salvia.victim_case SET agent_id = ? WHERE victim_case_i_code = ?`, newOperadorICode, caseICode)
	}

	// 4. Obtener nombre del operador para el historial
	var fullName string
	db.Raw(`SELECT gup.general_user_profile_names || ' ' || gup.general_user_profile_last_names FROM security.general_user gu JOIN security.general_user_profile gup ON gup.general_user_profile_id = gu.general_user_general_user_profile WHERE gu.general_user_i_code = ?`, newOperadorICode).Scan(&fullName)

	// 5. Actualizar victim_case_owner_description (append al historial)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newEntry := "(" + timestamp + ") " + fullName + " [Operador]"
	db.Exec(`UPDATE salvia.victim_case SET victim_case_owner_description = COALESCE(victim_case_owner_description, '') || ' , ' || ? WHERE victim_case_i_code = ?`, newEntry, caseICode)

	// 6. Reasignar seguimientos PENDIENTES al nuevo operador
	db.Exec(`UPDATE salvia.follow_up_v2 SET agent_id = ? WHERE case_id = ? AND status = 'PENDIENTE'`, newOperadorICode, caseICode)

	// 7. Reasignar tareas pendientes (CaseTask) al nuevo operador
	db.Exec(`UPDATE salvia.case_task SET assigned_user_id = ? WHERE case_id = ? AND status = 'ToDo'`, newOperadorICode, caseICode)

	return nil
}
