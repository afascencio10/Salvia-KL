package service

import (
	"bitsflow/internal/constants"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrFormNotFound = errors.New("form: registro no encontrado")

// ─── Form ─────────────────────────────────────────────────────────────────────

type CreateFormInput struct {
	Name        string
	Description string
	Status      string
}

type UpdateFormInput struct {
	Name        *string
	Description *string
	Status      *string
}

// ─── getFormStructure response types ─────────────────────────────────────────

type QuestionStructure struct {
	models.Question
	Options       []models.Option               `json:"options"`
	Conditions    []models.VisibilityCondition  `json:"conditions"`
	Modifications []models.RenderModification   `json:"modifications"`
}

type RepeaterStructure struct {
	models.RepeaterGroup
	Questions     []QuestionStructure           `json:"questions"`
	Conditions    []models.VisibilityCondition  `json:"conditions"`
	Modifications []models.RenderModification   `json:"modifications"`
}

type SectionStructure struct {
	models.FormSection
	Questions     []QuestionStructure           `json:"questions"`
	Repeaters     []RepeaterStructure           `json:"repeaters"`
	Conditions    []models.VisibilityCondition  `json:"conditions"`
	Modifications []models.RenderModification   `json:"modifications"`
	IsAnswered    bool                          `json:"isAnswered"`
	IsVisible     bool                          `json:"isVisible"`
}

type FormStructure struct {
	models.Form
	Sections []SectionStructure `json:"sections"`
}

// ─── getFormSubmission response types ────────────────────────────────────────

type RepeaterEntryStructure struct {
	models.RepeaterEntry
	Answers []models.Answer `json:"answers"`
}

type SubmissionStructure struct {
	models.FormSubmission
	DirectAnswers   []models.Answer          `json:"directAnswers"`
	RepeaterEntries []RepeaterEntryStructure `json:"repeaterEntries"`
}

// ─── getFormSubmissionStructured response types ───────────────────────────────

type QuestionStructured struct {
	QuestionStructure
	SubmissionVisible bool             `json:"submissionVisible"`
	FailedConditions  []FailedCondition `json:"failedConditions,omitempty"`
	Answer            *models.Answer   `json:"answer,omitempty"`
}

type RepeaterEntryQuestion struct {
	QuestionStructure
	SubmissionVisible bool             `json:"submissionVisible"`
	FailedConditions  []FailedCondition `json:"failedConditions,omitempty"`
	Answer            *models.Answer   `json:"answer,omitempty"`
}

type RepeaterEntryStructured struct {
	models.RepeaterEntry
	Questions []RepeaterEntryQuestion `json:"questions"`
}

type RepeaterStructured struct {
	models.RepeaterGroup
	Questions         []QuestionStructure          `json:"questions"`
	Conditions        []models.VisibilityCondition `json:"conditions"`
	SubmissionVisible bool                         `json:"submissionVisible"`
	FailedConditions  []FailedCondition             `json:"failedConditions,omitempty"`
	Entries           []RepeaterEntryStructured    `json:"entries"`
}

type SectionStructured struct {
	models.FormSection
	Questions         []QuestionStructured         `json:"questions"`
	Repeaters         []RepeaterStructured         `json:"repeaters"`
	Conditions        []models.VisibilityCondition `json:"conditions"`
	SubmissionVisible bool                         `json:"submissionVisible"`
	FailedConditions  []FailedCondition             `json:"failedConditions,omitempty"`
}

type FormSubmissionStructured struct {
	models.Form
	SubmissionID *string             `json:"submissionId,omitempty"`
	Sections     []SectionStructured `json:"sections"`
}

type FormService interface {
	GetByID(ctx context.Context, id string) (*models.Form, error)
	List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error)
	Create(ctx context.Context, input CreateFormInput) (*models.Form, error)
	Update(ctx context.Context, id string, input UpdateFormInput) (*models.Form, error)
	Delete(ctx context.Context, id string) error
	GetFormStructure(ctx context.Context, formID string) (*FormStructure, error)
	GetFormSubmission(ctx context.Context, submissionID string) (*SubmissionStructure, error)
	GetFormSubmissionStructured(ctx context.Context, formID, submissionID string) (*FormSubmissionStructured, error)
	LoadForm(ctx context.Context, formID, submissionID string, formState map[string]interface{}) (*LoadFormResult, error)
	SaveSection(ctx context.Context, input SaveSectionInput) (*LoadFormResult, error)
	OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error
	// OnSectionUpdate se dispara después de CADA guardado de sección (completo
	// o no) -- a diferencia de OnEndFormSubmission, que solo se dispara al
	// completar el formulario. Ver comentario en su implementación.
	OnSectionUpdate(ctx context.Context, formID, submissionID, actorID string) error
	TestFunction(ctx context.Context, fn, id, submissionID string) (interface{}, error)
}

type FormServiceDeps struct {
	FormRepo                   repository.FormRepository
	FormSectionRepo            repository.FormSectionRepository
	QuestionRepo               repository.QuestionRepository
	RepeaterGroupRepo          repository.RepeaterGroupRepository
	OptionRepo                 repository.OptionRepository
	VisibilityCondRepo         repository.VisibilityConditionRepository
	RenderModificationRepo     repository.RenderModificationRepository
	FormSubmissionRepo         repository.FormSubmissionRepository
	RepeaterEntryRepo          repository.RepeaterEntryRepository
	AnswerRepo                 repository.AnswerRepository
	FollowUpRepo               repository.FollowUpRepository
	EmergencyMeasureRepo       repository.EmergencyMeasureRepository
	PsychosocialSupportRepo    repository.PsychosocialSupportRepository
	EconomicStabilizationRepo  repository.EconomicStabilizationRepository
	BarrierV2Repo              repository.BarrierV2Repository
	CaseTimelineEventRepo      repository.CaseTimelineEventRepository
	AgentLightRepo             repository.AgentLightRepository
	MenTeamRemisionRepo        repository.MenTeamRemisionRepository
	DiscapacidadRemisionRepo   repository.DiscapacidadRemisionRepository
	CasoCierreService          CasoCierreService
	CaseRepo                   repository.VictimCaseLightRepository
	FollowUpV2Svc              FollowUpV2Service
	CaseTaskRepo               repository.CaseTaskRepository
	EntityLetterRepo           repository.EntityLetterRepository
	BarrierFollowUpRepo        repository.BarrierFollowUpRepository
	TeamContactRepo            repository.TeamContactRepository
	DuplaRepo                  repository.DuplaRepository
	VictimCaseFormSvc          VictimCaseFormService
}

type formService struct {
	repo                   repository.FormRepository
	formSectionRepo        repository.FormSectionRepository
	questionRepo           repository.QuestionRepository
	repeaterGroupRepo      repository.RepeaterGroupRepository
	optionRepo             repository.OptionRepository
	visibilityCondRepo         repository.VisibilityConditionRepository
	renderModificationRepo     repository.RenderModificationRepository
	submissionRepo         repository.FormSubmissionRepository
	repeaterEntryRepo      repository.RepeaterEntryRepository
	answerRepo             repository.AnswerRepository
	followUpRepo           repository.FollowUpRepository
	emergencyMeasureRepo   repository.EmergencyMeasureRepository
	psychosocialSupportRepo    repository.PsychosocialSupportRepository
	economicStabilizationRepo  repository.EconomicStabilizationRepository
	barrierV2Repo              repository.BarrierV2Repository
	caseTimelineRepo           repository.CaseTimelineEventRepository
	agentLightRepo             repository.AgentLightRepository
	menTeamRemisionRepo        repository.MenTeamRemisionRepository
	discapacidadRemisionRepo   repository.DiscapacidadRemisionRepository
	casoCierreService          CasoCierreService
	caseRepo                   repository.VictimCaseLightRepository
	followUpV2Svc              FollowUpV2Service
	caseTaskRepo               repository.CaseTaskRepository
	entityLetterRepo           repository.EntityLetterRepository
	barrierFollowUpRepo        repository.BarrierFollowUpRepository
	teamContactRepo            repository.TeamContactRepository
	duplaRepo                  repository.DuplaRepository
	victimCaseFormSvc          VictimCaseFormService
}

func NewFormService(deps FormServiceDeps) FormService {
	return &formService{
		repo:                      deps.FormRepo,
		formSectionRepo:           deps.FormSectionRepo,
		questionRepo:              deps.QuestionRepo,
		repeaterGroupRepo:         deps.RepeaterGroupRepo,
		optionRepo:                deps.OptionRepo,
		visibilityCondRepo:        deps.VisibilityCondRepo,
		renderModificationRepo:    deps.RenderModificationRepo,
		submissionRepo:            deps.FormSubmissionRepo,
		repeaterEntryRepo:         deps.RepeaterEntryRepo,
		answerRepo:                deps.AnswerRepo,
		followUpRepo:              deps.FollowUpRepo,
		emergencyMeasureRepo:      deps.EmergencyMeasureRepo,
		psychosocialSupportRepo:   deps.PsychosocialSupportRepo,
		economicStabilizationRepo: deps.EconomicStabilizationRepo,
		barrierV2Repo:             deps.BarrierV2Repo,
		caseTimelineRepo:          deps.CaseTimelineEventRepo,
		agentLightRepo:            deps.AgentLightRepo,
		menTeamRemisionRepo:       deps.MenTeamRemisionRepo,
		discapacidadRemisionRepo:  deps.DiscapacidadRemisionRepo,
		casoCierreService:         deps.CasoCierreService,
		caseRepo:                  deps.CaseRepo,
		followUpV2Svc:             deps.FollowUpV2Svc,
		caseTaskRepo:              deps.CaseTaskRepo,
		entityLetterRepo:          deps.EntityLetterRepo,
		barrierFollowUpRepo:       deps.BarrierFollowUpRepo,
		teamContactRepo:           deps.TeamContactRepo,
		duplaRepo:                 deps.DuplaRepo,
		victimCaseFormSvc:         deps.VictimCaseFormSvc,
	}
}

func (s *formService) GetByID(ctx context.Context, id string) (*models.Form, error) {
	form, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}
	return form, nil
}

func (s *formService) List(ctx context.Context, page, limit int) (repository.PageResult[models.Form], error) {
	return s.repo.FindWithPagination(ctx, page, limit)
}

func (s *formService) Create(ctx context.Context, input CreateFormInput) (*models.Form, error) {
	status := input.Status
	if status == "" {
		status = "active"
	}
	form := &models.Form{
		Name:        input.Name,
		Description: input.Description,
		Status:      status,
	}
	if err := s.repo.Create(ctx, form); err != nil {
		return nil, err
	}
	return form, nil
}

func (s *formService) Update(ctx context.Context, id string, input UpdateFormInput) (*models.Form, error) {
	fields := map[string]interface{}{}
	if input.Name != nil        { fields["name"] = *input.Name }
	if input.Description != nil { fields["description"] = *input.Description }
	if input.Status != nil      { fields["status"] = *input.Status }

	if err := s.repo.UpdateFields(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, id)
}

func (s *formService) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFormNotFound
	}
	return err
}

// GetFormStructure returns the full form tree (sections → questions/repeaters → options/conditions)
// using 6 bulk queries to avoid N+1 problems.
func (s *formService) GetFormStructure(ctx context.Context, formID string) (*FormStructure, error) {
	// Query 1: form
	form, err := s.repo.FindByID(ctx, formID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormNotFound
		}
		return nil, err
	}

	// Query 2: sections
	sections, err := s.formSectionRepo.FindByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	sectionIDs := make([]string, len(sections))
	for i, s := range sections {
		sectionIDs[i] = s.ID
	}

	// Query 3: all questions for the form (includes both section-level and repeater-level)
	questions, err := s.questionRepo.FindByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	// Query 4: all repeater groups for these sections
	repeaters, err := s.repeaterGroupRepo.FindBySectionIDs(ctx, sectionIDs)
	if err != nil {
		return nil, err
	}

	// Collect all IDs that need options and visibility conditions
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	repeaterIDs := make([]string, len(repeaters))
	for i, r := range repeaters {
		repeaterIDs[i] = r.ID
	}

	// Build full target ID list for bulk VC fetch: sections + repeaters + questions
	allTargetIDs := make([]string, 0, len(sectionIDs)+len(repeaterIDs)+len(questionIDs))
	allTargetIDs = append(allTargetIDs, sectionIDs...)
	allTargetIDs = append(allTargetIDs, repeaterIDs...)
	allTargetIDs = append(allTargetIDs, questionIDs...)

	// Query 5: all options
	options, err := s.optionRepo.FindByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, err
	}

	// Query 6: all visibility conditions
	vcs, err := s.visibilityCondRepo.FindByTargetIDs(ctx, allTargetIDs)
	if err != nil {
		return nil, err
	}

	// Query 7: all render modifications (non-fatal — tabla puede no existir aún)
	rms, err := s.renderModificationRepo.FindByTargetIDs(ctx, allTargetIDs)
	if err != nil {
		log.Printf("[WARN] render_modification query failed (tabla inexistente?): %v", err)
		rms = []models.RenderModification{}
	}

	// ── Build lookup maps ──────────────────────────────────────────────────────

	// options by questionID
	optionsByQuestion := map[string][]models.Option{}
	for _, o := range options {
		optionsByQuestion[o.QuestionID] = append(optionsByQuestion[o.QuestionID], o)
	}

	// VCs by composite key "TargetType:TargetID" — validates both fields
	vcByKey := map[string][]models.VisibilityCondition{}
	for _, vc := range vcs {
		key := strings.ToUpper(vc.TargetType) + ":" + vc.TargetID
		vcByKey[key] = append(vcByKey[key], vc)
	}

	// RMs by targetID (target_type no es necesario para indexar — el targetID ya es único por entidad)
	rmByTargetID := map[string][]models.RenderModification{}
	for _, rm := range rms {
		rmByTargetID[rm.TargetID] = append(rmByTargetID[rm.TargetID], rm)
	}

	// questions by sectionID (section-level) and by repeaterGroupID
	sectionQuestions  := map[string][]QuestionStructure{}
	repeaterQuestions := map[string][]QuestionStructure{}
	for _, q := range questions {
		qs := QuestionStructure{
			Question:      q,
			Options:       optionsByQuestion[q.ID],
			Conditions:    vcByKey["QUESTION:"+q.ID],
			Modifications: rmByTargetID[q.ID],
		}
		if qs.Options == nil        { qs.Options = []models.Option{} }
		if qs.Conditions == nil     { qs.Conditions = []models.VisibilityCondition{} }
		if qs.Modifications == nil  { qs.Modifications = []models.RenderModification{} }

		if q.RepeaterGroupID != nil {
			repeaterQuestions[*q.RepeaterGroupID] = append(repeaterQuestions[*q.RepeaterGroupID], qs)
		} else {
			sectionQuestions[q.FormSectionID] = append(sectionQuestions[q.FormSectionID], qs)
		}
	}

	// repeaters by sectionID
	repeatersBySection := map[string][]RepeaterStructure{}
	for _, rg := range repeaters {
		rs := RepeaterStructure{
			RepeaterGroup: rg,
			Questions:     repeaterQuestions[rg.ID],
			Conditions:    vcByKey["REPEATER_GROUP:"+rg.ID],
			Modifications: rmByTargetID[rg.ID],
		}
		if rs.Questions == nil      { rs.Questions = []QuestionStructure{} }
		if rs.Conditions == nil     { rs.Conditions = []models.VisibilityCondition{} }
		if rs.Modifications == nil  { rs.Modifications = []models.RenderModification{} }
		repeatersBySection[rg.FormSectionID] = append(repeatersBySection[rg.FormSectionID], rs)
	}

	// ── Assemble result ────────────────────────────────────────────────────────
	sectionStructures := make([]SectionStructure, len(sections))
	for i, sec := range sections {
		qs := sectionQuestions[sec.ID]
		rs := repeatersBySection[sec.ID]
		cs := vcByKey["SECTION:"+sec.ID]
		ms := rmByTargetID[sec.ID]
		if qs == nil { qs = []QuestionStructure{} }
		if rs == nil { rs = []RepeaterStructure{} }
		if cs == nil { cs = []models.VisibilityCondition{} }
		if ms == nil { ms = []models.RenderModification{} }
		sectionStructures[i] = SectionStructure{
			FormSection:   sec,
			Questions:     qs,
			Repeaters:     rs,
			Conditions:    cs,
			Modifications: ms,
		}
	}

	return &FormStructure{
		Form:     *form,
		Sections: sectionStructures,
	}, nil
}

// GetFormSubmission returns a submission with its direct answers and repeater entries (each with their answers).
// Uses 3 queries: submission + directAnswers + repeaterEntries, then 1 bulk query for all entry answers.
func (s *formService) GetFormSubmission(ctx context.Context, submissionID string) (*SubmissionStructure, error) {
	// PASO 1 — Consultar Submission
	submission, err := s.submissionRepo.FindByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFormSubmissionNotFound
		}
		return nil, err
	}

	// PASO 2 — Consultar Answers directas (sin repeaterEntryId)
	directAnswers, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if directAnswers == nil {
		directAnswers = []models.Answer{}
	}

	// PASO 3 — Consultar RepeaterEntries ordenadas
	entries, err := s.repeaterEntryRepo.FindBySubmissionIDOrdered(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	// Bulk fetch de todas las answers de los entries (evita N+1)
	entryIDs := make([]string, len(entries))
	for i, e := range entries {
		entryIDs[i] = e.ID
	}
	allEntryAnswers, err := s.answerRepo.FindByRepeaterEntryIDs(ctx, entryIDs)
	if err != nil {
		return nil, err
	}

	// Agrupar answers por entryID
	answersByEntry := map[string][]models.Answer{}
	for _, a := range allEntryAnswers {
		if a.RepeaterEntryID != nil {
			answersByEntry[*a.RepeaterEntryID] = append(answersByEntry[*a.RepeaterEntryID], a)
		}
	}

	// Armar estructuras de entries
	entryStructures := make([]RepeaterEntryStructure, len(entries))
	for i, e := range entries {
		answers := answersByEntry[e.ID]
		if answers == nil {
			answers = []models.Answer{}
		}
		entryStructures[i] = RepeaterEntryStructure{
			RepeaterEntry: e,
			Answers:       answers,
		}
	}

	return &SubmissionStructure{
		FormSubmission:  *submission,
		DirectAnswers:   directAnswers,
		RepeaterEntries: entryStructures,
	}, nil
}

// GetFormSubmissionStructured combina getFormStructure + getFormSubmission y devuelve
// el árbol del form con las respuestas del submission asociadas a cada pregunta.
func (s *formService) GetFormSubmissionStructured(ctx context.Context, formID, submissionID string) (*FormSubmissionStructured, error) {
	fs, err := s.GetFormStructure(ctx, formID)
	if err != nil {
		return nil, err
	}

	// Índices de respuestas (se populan solo si hay submissionID)
	directByQuestion := map[string]models.Answer{}
	entriesByGroup := map[string][]RepeaterEntryStructure{}

	var sub *SubmissionStructure
	var subID *string
	if submissionID != "" {
		var err error
		sub, err = s.GetFormSubmission(ctx, submissionID)
		if err != nil {
			return nil, err
		}
		subID = &sub.ID

		for _, a := range sub.DirectAnswers {
			directByQuestion[a.QuestionID] = a
		}
		for _, e := range sub.RepeaterEntries {
			entriesByGroup[e.RepeaterGroupID] = append(entriesByGroup[e.RepeaterGroupID], e)
		}
	}

	// Construir secciones estructuradas
	sections := make([]SectionStructured, len(fs.Sections))
	for i, sec := range fs.Sections {
		secVis := checkVisibility(fs, sub, nil, sectionItem(sec), nil)

		// Preguntas directas con su answer y visibilidad
		questions := make([]QuestionStructured, len(sec.Questions))
		for j, q := range sec.Questions {
			qVis := checkVisibility(fs, sub, nil, directQuestionItem(q, sec.Order), nil)
			qs := QuestionStructured{
				QuestionStructure: q,
				SubmissionVisible: qVis.Visible,
				FailedConditions:  qVis.FailedConditions,
			}
			if a, ok := directByQuestion[q.ID]; ok {
				qs.Answer = &a
			}
			questions[j] = qs
		}

		// Repeaters con visibilidad, entries y preguntas por entry
		repeaters := make([]RepeaterStructured, len(sec.Repeaters))
		for k, r := range sec.Repeaters {
			rVis := checkVisibility(fs, sub, nil, repeaterItem(r, sec.Order), nil)

			rawEntries := entriesByGroup[r.ID]
			entries := make([]RepeaterEntryStructured, len(rawEntries))
			for l, e := range rawEntries {
				entryQuestions := make([]RepeaterEntryQuestion, len(r.Questions))
				for m, q := range r.Questions {
					qVis := checkVisibility(fs, sub, e.Answers, repeaterQuestionItem(q, sec.Order, r.Order, r.ID), nil)
					eq := RepeaterEntryQuestion{
						QuestionStructure: q,
						SubmissionVisible: qVis.Visible,
						FailedConditions:  qVis.FailedConditions,
					}
					for idx := range e.Answers {
						if e.Answers[idx].QuestionID == q.ID {
							eq.Answer = &e.Answers[idx]
							break
						}
					}
					entryQuestions[m] = eq
				}
				entries[l] = RepeaterEntryStructured{
					RepeaterEntry: e.RepeaterEntry,
					Questions:     entryQuestions,
				}
			}

			repeaters[k] = RepeaterStructured{
				RepeaterGroup:     r.RepeaterGroup,
				Questions:         r.Questions,
				Conditions:        r.Conditions,
				SubmissionVisible: rVis.Visible,
				FailedConditions:  rVis.FailedConditions,
				Entries:           entries,
			}
		}

		sections[i] = SectionStructured{
			FormSection:       sec.FormSection,
			Questions:         questions,
			Repeaters:         repeaters,
			Conditions:        sec.Conditions,
			SubmissionVisible: secVis.Visible,
			FailedConditions:  secVis.FailedConditions,
		}
	}

	return &FormSubmissionStructured{
		Form:         fs.Form,
		SubmissionID: subID,
		Sections:     sections,
	}, nil
}

func (s *formService) TestFunction(ctx context.Context, fn, id, submissionID string) (interface{}, error) {
	switch fn {

	// Lista submissions de un form para obtener IDs reales
	// ?fn=listSubmissions&id=<formId>
	case "listSubmissions":
		return s.submissionRepo.FindByFormID(ctx, id)

	// Valida que cada pregunta tenga solo 1 Answer (sin duplicados por questionId)
	// ?fn=validateAnswers&submissionId=<submissionId>
	case "validateAnswers":
		answers, err := s.answerRepo.FindBySubmissionID(ctx, submissionID)
		if err != nil {
			return nil, err
		}
		countByQuestion := map[string]int{}
		for _, a := range answers {
			countByQuestion[a.QuestionID]++
		}
		duplicates := map[string]int{}
		for qID, count := range countByQuestion {
			if count > 1 {
				duplicates[qID] = count
			}
		}
		return map[string]interface{}{
			"totalAnswers":    len(answers),
			"uniqueQuestions": len(countByQuestion),
			"duplicates":      duplicates,
			"ok":              len(duplicates) == 0,
		}, nil

	// Devuelve FormStructure con answers y entries asociadas
	// ?fn=getFormSubmissionStructured&id=<formId>&submissionId=<submissionId>
	case "getFormSubmissionStructured":
		return s.GetFormSubmissionStructured(ctx, id, submissionID)

	// Evalúa visibilidad de todos los items del form
	// ?fn=checkVisibility&id=<formId>&submissionId=<submissionId>
	case "checkVisibility":
		fs, err := s.GetFormStructure(ctx, id)
		if err != nil {
			return nil, err
		}
		var submission *SubmissionStructure
		if submissionID != "" {
			submission, err = s.GetFormSubmission(ctx, submissionID)
			if err != nil {
				return nil, err
			}
		}
		// Indexar entries por repeaterGroupID para acceso rápido
		entriesByGroup := map[string][]RepeaterEntryStructure{}
		if submission != nil {
			for _, e := range submission.RepeaterEntries {
				entriesByGroup[e.RepeaterGroupID] = append(entriesByGroup[e.RepeaterGroupID], e)
			}
		}

		result := map[string]interface{}{}
		for _, sec := range fs.Sections {
			result["section:"+sec.ID] = checkVisibility(fs, submission, nil, sectionItem(sec), nil)
			for _, q := range sec.Questions {
				result["question:"+q.ID] = checkVisibility(fs, submission, nil, directQuestionItem(q, sec.Order), nil)
			}
			for _, r := range sec.Repeaters {
				result["repeater:"+r.ID] = checkVisibility(fs, submission, nil, repeaterItem(r, sec.Order), nil)
				for _, q := range r.Questions {
					entries := entriesByGroup[r.ID]
					if len(entries) == 0 {
						// Sin entries: evaluar sin entryAnswers
						result["repeaterQ:"+q.ID+":noEntry"] = checkVisibility(fs, submission, nil, repeaterQuestionItem(q, sec.Order, r.Order, r.ID), nil)
					} else {
						// Evaluar para cada entry
						for _, entry := range entries {
							key := "repeaterQ:" + q.ID + ":entry" + entry.ID
							result[key] = checkVisibility(fs, submission, entry.Answers, repeaterQuestionItem(q, sec.Order, r.Order, r.ID), nil)
						}
					}
				}
			}
		}
		return result, nil

	// Pruebas de validateAnswer
	// ?fn=validateAnswer
	case "validateAnswer":
		type testCase struct {
			Desc   string              `json:"desc"`
			Input  string              `json:"input"`
			Want   bool                `json:"want"`
			Result ValidateAnswerResult `json:"result"`
			Pass   bool                `json:"pass"`
		}

		makeQ := func(typeID string, required bool) QuestionStructure {
			return QuestionStructure{
				Question: models.Question{QuestionTypeID: typeID, Required: required},
			}
		}
		makeA := func(val string) models.Answer {
			return models.Answer{Value: val}
		}

		cases := []testCase{
			// required=false → siempre válido sin importar el valor
			{"required=false ignora valor vacío",   "",            true,  ValidateAnswerResult{}, false},
			{"required=false ignora valor inválido","abc",         true,  ValidateAnswerResult{}, false},

			// text
			{"text requerido vacío",                "",            false, ValidateAnswerResult{}, false},
			{"text requerido solo espacios",        "   ",         false, ValidateAnswerResult{}, false},
			{"text requerido con valor",            "hola",        true,  ValidateAnswerResult{}, false},

			// number
			{"number requerido vacío",              "",            false, ValidateAnswerResult{}, false},
			{"number requerido texto",              "abc",         false, ValidateAnswerResult{}, false},
			{"number requerido decimal",            "3.14",        true,  ValidateAnswerResult{}, false},
			{"number requerido negativo",           "-5",          true,  ValidateAnswerResult{}, false},

			// date
			{"date requerido vacío",                "",            false, ValidateAnswerResult{}, false},
			{"date requerido mes inválido",         "2026-13-01",  false, ValidateAnswerResult{}, false},
			{"date requerido formato incorrecto",   "11/04/2026",  false, ValidateAnswerResult{}, false},
			{"date requerido formato correcto",     "2026-04-11",  true,  ValidateAnswerResult{}, false},

			// datetime
			{"datetime requerido vacío",            "",                    false, ValidateAnswerResult{}, false},
			{"datetime requerido solo fecha",       "2026-04-11",          false, ValidateAnswerResult{}, false},
			{"datetime requerido datetime-local",   "2026-04-11T17:30",    true,  ValidateAnswerResult{}, false},
			{"datetime requerido RFC3339",          "2026-04-11T17:30:00Z",true,  ValidateAnswerResult{}, false},

			// single / dropdown
			{"single requerido vacío",              "",       false, ValidateAnswerResult{}, false},
			{"single requerido con valor",          "Alto",   true,  ValidateAnswerResult{}, false},
			{"dropdown requerido vacío",            "",       false, ValidateAnswerResult{}, false},
			{"dropdown requerido con valor",        "Ana M.", true,  ValidateAnswerResult{}, false},

			// boolean
			{"boolean requerido vacío",             "",       false, ValidateAnswerResult{}, false},
			{"boolean requerido valor inválido",    "maybe",  false, ValidateAnswerResult{}, false},
			{"boolean requerido false es válido",   "false",  true,  ValidateAnswerResult{}, false},
			{"boolean requerido true es válido",    "true",   true,  ValidateAnswerResult{}, false},

			// multiple
			{"multiple requerido vacío",            "",               false, ValidateAnswerResult{}, false},
			{"multiple requerido con valores csv",  "Salud,Justicia", true,  ValidateAnswerResult{}, false},
		}

		// Tipos correspondientes a cada caso (alineados con el slice anterior)
		types := []struct{ typeID string; required bool }{
			{"text",     false}, {"number",   false},
			{"text",     true},  {"text",     true},  {"text",     true},
			{"number",   true},  {"number",   true},  {"number",   true},  {"number",   true},
			{"date",     true},  {"date",     true},  {"date",     true},  {"date",     true},
			{"datetime", true},  {"datetime", true},  {"datetime", true},  {"datetime", true},
			{"single",   true},  {"single",   true},  {"dropdown", true},  {"dropdown", true},
			{"boolean",  true},  {"boolean",  true},  {"boolean",  true},  {"boolean",  true},
			{"multiple", true},  {"multiple", true},
		}

		passed, failed := 0, 0
		for i := range cases {
			q := makeQ(types[i].typeID, types[i].required)
			a := makeA(cases[i].Input)
			result := validateAnswer(a, q)
			cases[i].Result = result
			cases[i].Pass = result.Valid == cases[i].Want
			if cases[i].Pass {
				passed++
			} else {
				failed++
			}
		}

		return map[string]interface{}{
			"passed": passed,
			"failed": failed,
			"total":  len(cases),
			"ok":     failed == 0,
			"cases":  cases,
		}, nil

	// Pruebas de isAnsweredQuestion
	// ?fn=isAnsweredQuestion&id=<formId>&submissionId=<submissionId>
	case "isAnsweredQuestion":
		fs, err := s.GetFormStructure(ctx, id)
		if err != nil {
			return nil, err
		}
		var sub *SubmissionStructure
		if submissionID != "" {
			sub, err = s.GetFormSubmission(ctx, submissionID)
			if err != nil {
				return nil, err
			}
		}

		// Índice de respuestas directas por questionID
		answerByQuestion := map[string]*models.Answer{}
		if sub != nil {
			for i := range sub.DirectAnswers {
				a := sub.DirectAnswers[i]
				answerByQuestion[a.QuestionID] = &a
			}
		}

		// Helper: encuentra QuestionStructure y sectionOrder por ID
		type qWithOrder struct {
			q     QuestionStructure
			order int
		}
		findQ := func(questionID string) (qWithOrder, bool) {
			for _, sec := range fs.Sections {
				for _, q := range sec.Questions {
					if q.ID == questionID {
						return qWithOrder{q, sec.Order}, true
					}
				}
			}
			return qWithOrder{}, false
		}

		type testCase struct {
			Desc       string          `json:"desc"`
			QuestionID string          `json:"questionId"`
			HasAnswer  bool            `json:"hasAnswer"`
			AnswerVal  string          `json:"answerValue"`
			Want       IsAnsweredResult `json:"want"`
			Got        IsAnsweredResult `json:"got"`
			Pass       bool            `json:"pass"`
		}

		makeAnswer := func(qID, val string) *models.Answer {
			a := models.Answer{QuestionID: qID, Value: val}
			return &a
		}

		// IDs conocidos del form de prueba
		const (
			qNivelRiesgo       = "a1937f54-9282-48bb-bf6b-2a8bd9933c85" // single, required
			qNuevosHechos      = "da90b1c9-45a1-4e2a-83dd-a2ae77cf5227" // boolean, required
			qDescHechos        = "f8453544-2a8c-461d-9986-302ee4719492" // text, NOT required — visible solo si nuevos_hechos=true
			qCodigoInterno     = "81d7b17f-3438-4d07-b72f-b45f858d56c4" // text, NOT required
			qInfoPsicosocial   = "6891fbec-da1d-46e6-9c3b-3bb824d19631" // text, required
		)

		// Para los casos de visibilidad condicional construimos un sub sintético
		// con nuevos_hechos="false" → descripcion_hechos NO visible
		subConFalse := &SubmissionStructure{
			DirectAnswers: []models.Answer{
				{QuestionID: qNuevosHechos, Value: "false"},
			},
		}
		// con nuevos_hechos="true" → descripcion_hechos SÍ visible
		subConTrue := &SubmissionStructure{
			DirectAnswers: []models.Answer{
				{QuestionID: qNuevosHechos, Value: "true"},
			},
		}

		cases := []struct {
			desc     string
			qID      string
			answer   *models.Answer
			sub      *SubmissionStructure
			wantVis  bool
			wantAns  bool
		}{
			// visible, required, sin respuesta → not answered
			{"required sin respuesta (nil)",          qNivelRiesgo,     nil,                          sub,         true,  false},
			// visible, required, answer vacío → not answered
			{"required answer vacío",                 qNivelRiesgo,     makeAnswer(qNivelRiesgo, ""), sub,         true,  false},
			// visible, required, answer válido → answered
			{"required answer válido",                qNivelRiesgo,     makeAnswer(qNivelRiesgo, "Alto"), sub,     true,  true},
			// visible, NOT required, sin respuesta → answered (no es requerida)
			{"not required sin respuesta",            qCodigoInterno,   nil,                          sub,         true,  true},
			// NOT visible (condición no cumplida: nuevos_hechos=false) → answered
			{"not visible por condición false",       qDescHechos,      nil,                          subConFalse, false, true},
			// SÍ visible (condición cumplida: nuevos_hechos=true), not required, sin answer → answered
			{"visible por condición true, not req",   qDescHechos,      nil,                          subConTrue,  true,  true},
			// sin submission (sub=nil): condición no evaluable → visible por defecto, required, sin answer
			{"sub nil, required sin respuesta",       qInfoPsicosocial, nil,                          nil,         true,  false},
		}

		type result struct {
			Desc   string          `json:"desc"`
			Got    IsAnsweredResult `json:"got"`
			Pass   bool            `json:"pass"`
		}
		var results []result
		passed, failed := 0, 0

		for _, c := range cases {
			qwo, ok := findQ(c.qID)
			if !ok {
				results = append(results, result{c.desc + " [PREGUNTA NO ENCONTRADA]", IsAnsweredResult{}, false})
				failed++
				continue
			}
			got := isAnsweredQuestion(qwo.q, qwo.order, c.answer, fs, c.sub, nil)
			pass := got.IsVisible == c.wantVis && got.IsAnswered == c.wantAns
			if pass { passed++ } else { failed++ }
			results = append(results, result{c.desc, got, pass})
		}

		return map[string]interface{}{
			"passed":  passed,
			"failed":  failed,
			"total":   len(cases),
			"ok":      failed == 0,
			"cases":   results,
		}, nil

	// Pruebas de isAnsweredRepeater
	// ?fn=isAnsweredRepeater&id=<formId>
	case "isAnsweredRepeater":
		fs, err := s.GetFormStructure(ctx, id)
		if err != nil {
			return nil, err
		}

		// IDs conocidos
		const (
			subSinRespuestas = "90b396bb-54c6-48aa-a41b-1a59ad73d518"
			subConRepeater   = "ed09276b-47b5-4ea5-9a43-49afb45f9d7b"
			repeaterID       = "5fd3ecdc-2e5f-4b31-97ef-8a994580586a"
			qSector          = "f19378b6-55c5-4fdf-b765-7ebcc3978741" // dropdown, required
		)

		// Cargar submissions reales
		subVacio, err := s.GetFormSubmission(ctx, subSinRespuestas)
		if err != nil {
			return nil, err
		}
		subConEntries, err := s.GetFormSubmission(ctx, subConRepeater)
		if err != nil {
			return nil, err
		}

		// Encontrar el repeater y su sección en el form
		var repeater    RepeaterStructure
		var sectionOrder int
		for _, sec := range fs.Sections {
			for _, r := range sec.Repeaters {
				if r.ID == repeaterID {
					repeater     = r
					sectionOrder = sec.Order
				}
			}
		}

		// Filtrar entries de cada submission para este repeater
		entriesVacio := []RepeaterEntryStructure{}
		for _, e := range subVacio.RepeaterEntries {
			if e.RepeaterGroupID == repeaterID {
				entriesVacio = append(entriesVacio, e)
			}
		}
		entriesCompletas := []RepeaterEntryStructure{}
		for _, e := range subConEntries.RepeaterEntries {
			if e.RepeaterGroupID == repeaterID {
				entriesCompletas = append(entriesCompletas, e)
			}
		}

		// Repeater sintético con MinRepetitions=1 (para probar PASO 2)
		repeaterConMin := repeater
		repeaterConMin.MinRepetitions = 1

		// Entry sintética con sector vacío (required, dropdown)
		entryConSectorVacio := RepeaterEntryStructure{
			RepeaterEntry: models.RepeaterEntry{
				ID:              "synthetic-entry-001",
				RepeaterGroupID: repeaterID,
				Iteration:       1,
			},
			Answers: []models.Answer{
				{QuestionID: qSector, Value: ""},
			},
		}

		// Sub sintético: nuevos_hechos=false → para probar "repeater no visible"
		// (el repeater de barreras no tiene condición, así que usamos nil sub
		//  y verificamos que sin condiciones siempre es visible)

		type testCase struct {
			Desc       string                   `json:"desc"`
			WantVis    bool                     `json:"wantVisible"`
			WantAns    bool                     `json:"wantAnswered"`
			Got        IsAnsweredRepeaterResult `json:"got"`
			Pass       bool                     `json:"pass"`
		}

		runCase := func(desc string, r RepeaterStructure, entries []RepeaterEntryStructure, sub *SubmissionStructure, wantVis, wantAns bool) testCase {
			got  := isAnsweredRepeater(r, sectionOrder, entries, fs, sub, nil)
			pass := got.IsVisible == wantVis && got.IsAnswered == wantAns
			return testCase{desc, wantVis, wantAns, got, pass}
		}

		cases := []testCase{
			// minRep=0, sin entries → answered (loop vacío)
			runCase("minRep=0, sin entries → ok",
				repeater, entriesVacio, subVacio, true, true),

			// minRep=1, sin entries → not answered
			runCase("minRep=1, sin entries → no ok",
				repeaterConMin, entriesVacio, subVacio, true, false),

			// minRep=0, entries reales todas con sector respondido → answered
			runCase("entries completas (sub real) → ok",
				repeater, entriesCompletas, subConEntries, true, true),

			// minRep=0, entry sintética con sector vacío → not answered
			runCase("entry con required vacío → no ok",
				repeater, []RepeaterEntryStructure{entryConSectorVacio}, subVacio, true, false),

			// minRep=2, solo 1 entry → not answered (no cumple mínimo)
			runCase("minRep=2, 1 entry → no ok",
				func() RepeaterStructure { r := repeater; r.MinRepetitions = 2; return r }(),
				[]RepeaterEntryStructure{entryConSectorVacio}, subVacio, true, false),
		}

		passed, failed := 0, 0
		for _, c := range cases {
			if c.Pass { passed++ } else { failed++ }
		}

		return map[string]interface{}{
			"passed": passed,
			"failed": failed,
			"total":  len(cases),
			"ok":     failed == 0,
			"cases":  cases,
		}, nil

	// Pruebas de isAnsweredSection
	// ?fn=isAnsweredSection&id=<formId>
	case "isAnsweredSection":
		fs, err := s.GetFormStructure(ctx, id)
		if err != nil {
			return nil, err
		}

		const (
			subSinRespuestasID = "90b396bb-54c6-48aa-a41b-1a59ad73d518"
			subSec1ID          = "1c5b8876-e47a-4c77-9d74-700ff5034ef1"
			subSec2ID          = "10a14015-20fc-4e19-b6dc-a6d2259895da"
			subSec4ID          = "ed09276b-47b5-4ea5-9a43-49afb45f9d7b"
		)

		loadSub := func(subID string) *SubmissionStructure {
			sub, e := s.GetFormSubmission(ctx, subID)
			if e != nil {
				return nil
			}
			return sub
		}

		subVacio  := loadSub(subSinRespuestasID)
		subSec1   := loadSub(subSec1ID)
		subSec2   := loadSub(subSec2ID)
		subSec4   := loadSub(subSec4ID)

		// Índice de secciones por order
		secByOrder := map[int]SectionStructure{}
		for _, sec := range fs.Sections {
			secByOrder[sec.Order] = sec
		}

		type testCase struct {
			Desc    string          `json:"desc"`
			WantVis bool            `json:"wantVisible"`
			WantAns bool            `json:"wantAnswered"`
			Got     IsAnsweredResult `json:"got"`
			Pass    bool            `json:"pass"`
		}

		run := func(desc string, secOrder int, sub *SubmissionStructure, wantVis, wantAns bool) testCase {
			got  := isAnsweredSection(secByOrder[secOrder], fs, sub, nil)
			pass := got.IsVisible == wantVis && got.IsAnswered == wantAns
			return testCase{desc, wantVis, wantAns, got, pass}
		}

		cases := []testCase{
			// sub nil con sección que tiene required → no respondida
			run("sub nil, sección 1 (tiene required)",         1, nil,      true, false),
			// sin ninguna respuesta, sección 1 tiene required → no respondida
			run("sin respuestas, sección 1",                   1, subVacio, true, false),
			// sub con sección 1 completa → respondida
			run("sub sec1, sección 1 → ok",                    1, subSec1,  true, true),
			// mismo sub pero evaluando sección 2 → no respondida (no llega a sec2)
			run("sub sec1, sección 2 → no ok",                 2, subSec1,  true, false),
			// sub hasta sec2, sección 2 → respondida
			run("sub sec2, sección 2 → ok",                    2, subSec2,  true, true),
			// sub hasta sec2, sección 3 → no respondida
			run("sub sec2, sección 3 → no ok",                 3, subSec2,  true, false),
			// sub hasta sec4 con repeater (minRep=0) → sección 4 respondida
			run("sub sec4, sección 4 (repeater minRep=0) → ok",4, subSec4,  true, true),
		}

		passed, failed := 0, 0
		for _, c := range cases {
			if c.Pass {
				passed++
			} else {
				failed++
			}
		}

		return map[string]interface{}{
			"passed": passed,
			"failed": failed,
			"total":  len(cases),
			"ok":     failed == 0,
			"cases":  cases,
		}, nil

	// Pruebas de loadForm
	// ?fn=loadForm&id=<formId>                         → sin submission
	// ?fn=loadForm&id=<formId>&submissionId=<subId>    → con submission
	case "loadForm":
		const (
			subSinRespuestasID = "90b396bb-54c6-48aa-a41b-1a59ad73d518"
			subSec1ID          = "1c5b8876-e47a-4c77-9d74-700ff5034ef1"
			subSec2ID          = "10a14015-20fc-4e19-b6dc-a6d2259895da"
			subSec4ID          = "ed09276b-47b5-4ea5-9a43-49afb45f9d7b"
		)

		type testCase struct {
			Desc              string `json:"desc"`
			WantCurrentSecOrd int    `json:"wantCurrentSectionOrder"`
			GotCurrentSecOrd  int    `json:"gotCurrentSectionOrder"`
			GotCurrentSecName string `json:"gotCurrentSectionName"`
			Pass              bool   `json:"pass"`
		}

		run := func(desc string, subID string, wantOrder int) testCase {
			res, e := s.LoadForm(ctx, id, subID, nil)
			if e != nil {
				return testCase{desc + " [ERROR: " + e.Error() + "]", wantOrder, -1, "", false}
			}
			gotOrder := 0
			gotName  := ""
			if res.CurrentSection != nil {
				gotOrder = res.CurrentSection.Order
				gotName  = res.CurrentSection.Name
			}
			return testCase{desc, wantOrder, gotOrder, gotName, gotOrder == wantOrder}
		}

		cases := []testCase{
			run("sin submission → sec 1",                    "",              1),
			run("sub sin respuestas → sec 1",                subSinRespuestasID, 1),
			run("sub sec1 completa → sec 2",                 subSec1ID,       2),
			run("sub sec1+2 completa → sec 3",               subSec2ID,       3),
			run("sub sec1-4 completa → sec 5",               subSec4ID,       5),
		}

		passed, failed := 0, 0
		for _, c := range cases {
			if c.Pass { passed++ } else { failed++ }
		}

		return map[string]interface{}{
			"passed": passed,
			"failed": failed,
			"total":  len(cases),
			"ok":     failed == 0,
			"cases":  cases,
		}, nil

	// Default: lista submissions del form de prueba
	default:
		return s.submissionRepo.FindByFormID(ctx, "2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff")
	}
}

// ─── visibilityItem ───────────────────────────────────────────────────────────

type visibilityItem struct {
	ID              string
	Name            string
	Conditions      []models.VisibilityCondition
	SectionOrder    int
	OrderInSection  int     // 0 para secciones → solo secciones anteriores pueden disparar
	RepeaterGroupID *string // non-nil si el item está dentro de un repeater
	OrderInRepeater int     // solo si RepeaterGroupID != nil
}

func sectionItem(sec SectionStructure) visibilityItem {
	return visibilityItem{
		ID:           sec.ID,
		Name:         sec.Name,
		Conditions:   sec.Conditions,
		SectionOrder: sec.Order,
	}
}

func repeaterItem(r RepeaterStructure, sectionOrder int) visibilityItem {
	return visibilityItem{
		ID:             r.ID,
		Name:           r.Name,
		Conditions:     r.Conditions,
		SectionOrder:   sectionOrder,
		OrderInSection: r.Order,
	}
}

func directQuestionItem(q QuestionStructure, sectionOrder int) visibilityItem {
	return visibilityItem{
		ID:             q.ID,
		Name:           q.Description,
		Conditions:     q.Conditions,
		SectionOrder:   sectionOrder,
		OrderInSection: q.Order,
	}
}

func repeaterQuestionItem(q QuestionStructure, sectionOrder, repeaterOrder int, repeaterGroupID string) visibilityItem {
	rgID := repeaterGroupID
	return visibilityItem{
		ID:              q.ID,
		Name:            q.Description,
		Conditions:      q.Conditions,
		SectionOrder:    sectionOrder,
		OrderInSection:  repeaterOrder,
		RepeaterGroupID: &rgID,
		OrderInRepeater: q.Order,
	}
}

// ─── questionRef ──────────────────────────────────────────────────────────────

type questionRef struct {
	Description     string
	SectionOrder    int
	RepeaterGroupID *string
	OrderInSection  int // order en sección (si directa) o order del repeater (si en repeater)
	OrderInRepeater int // solo si RepeaterGroupID != nil
}

func findQuestion(fs *FormStructure, questionID string) *questionRef {
	for _, sec := range fs.Sections {
		for _, q := range sec.Questions {
			if q.ID == questionID {
				return &questionRef{
					Description:    q.Description,
					SectionOrder:   sec.Order,
					OrderInSection: q.Order,
				}
			}
		}
		for _, r := range sec.Repeaters {
			for _, q := range r.Questions {
				if q.ID == questionID {
					rgID := r.ID
					return &questionRef{
						Description:     q.Description,
						SectionOrder:    sec.Order,
						RepeaterGroupID: &rgID,
						OrderInSection:  r.Order,
						OrderInRepeater: q.Order,
					}
				}
			}
		}
	}
	return nil
}

// ─── isBeforeInFlow ───────────────────────────────────────────────────────────
// Determina si una trigger directa viene ANTES de item en el flujo del form.
// Regla unificada (funciona para Section, RepeaterGroup y direct Question):
//   trigger.SectionOrder < item.SectionOrder
//   OR (misma sección AND trigger.OrderInSection < item.OrderInSection)
func isBeforeInFlow(trigger *questionRef, item visibilityItem) bool {
	return trigger.SectionOrder < item.SectionOrder ||
		(trigger.SectionOrder == item.SectionOrder &&
			trigger.OrderInSection < item.OrderInSection)
}

// ─── resolveStatePath ─────────────────────────────────────────────────────────

// resolveStatePath resuelve un path dotted (ej. "currentCase.status", "items.0.name")
// en un objeto anidado de interface{}. Soporta maps y arrays (índice numérico).
// Retorna string vacío si el path no existe, el objeto es nil, o el índice está fuera de rango.
func resolveStatePath(state map[string]interface{}, path string) string {
	if state == nil || path == "" {
		return ""
	}
	var current interface{} = state
	for _, key := range strings.Split(path, ".") {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		case []interface{}:
			idx, err := strconv.Atoi(key)
			if err != nil || idx < 0 || idx >= len(v) {
				return ""
			}
			current = v[idx]
		default:
			return ""
		}
		if current == nil {
			return ""
		}
	}
	return fmt.Sprintf("%v", current)
}

// ─── checkVisibility ──────────────────────────────────────────────────────────

type FailedCondition struct {
	models.VisibilityCondition
	TriggerQuestionName string `json:"triggerQuestionName"`
}

type VisibilityResult struct {
	Name             string            `json:"name"`
	Visible          bool              `json:"visible"`
	FailedConditions []FailedCondition `json:"failedConditions,omitempty"`
}

// checkVisibility evalúa si un item es visible según sus condiciones (todas AND).
//
//   fs           → árbol completo del form (getFormStructure)
//   submission   → submission con directAnswers y repeaterEntries (getFormSubmission), puede ser nil
//   entryAnswers → respuestas de la entry actual (solo cuando itemToCheck está en un repeaterGroup)
//   item         → elemento a evaluar
//   formState    → estado externo del padre; nil = sin condiciones de estado
func checkVisibility(
	fs          *FormStructure,
	submission  *SubmissionStructure,
	entryAnswers []models.Answer,
	item        visibilityItem,
	formState   map[string]interface{},
) VisibilityResult {
	if len(item.Conditions) == 0 {
		return VisibilityResult{Name: item.Name, Visible: true}
	}

	var directAnswers []models.Answer
	if submission != nil {
		directAnswers = submission.DirectAnswers
	}

	var applicable []models.VisibilityCondition
	for _, cond := range item.Conditions {
		// Condición de formState — aplica siempre, sin restricción de orden en el flujo
		if cond.TriggerStatePath != nil && *cond.TriggerStatePath != "" {
			applicable = append(applicable, cond)
			continue
		}

		// Condición de pregunta — lógica existente
		trigger := findQuestion(fs, cond.TriggerQuestionID)
		if trigger == nil {
			continue
		}

		if trigger.RepeaterGroupID != nil {
			if item.RepeaterGroupID == nil || *item.RepeaterGroupID != *trigger.RepeaterGroupID {
				continue
			}
			if entryAnswers == nil {
				continue
			}
			if trigger.OrderInRepeater >= item.OrderInRepeater {
				continue
			}
			applicable = append(applicable, cond)
		} else {
			if !isBeforeInFlow(trigger, item) {
				continue
			}
			applicable = append(applicable, cond)
		}
	}

	if len(applicable) == 0 {
		return VisibilityResult{Name: item.Name, Visible: true}
	}

	var failed []FailedCondition
	for _, cond := range applicable {
		trigVal := ""
		if cond.TriggerValue != nil {
			trigVal = *cond.TriggerValue
		}

		// Condición de formState
		if cond.TriggerStatePath != nil && *cond.TriggerStatePath != "" {
			stateVal := resolveStatePath(formState, *cond.TriggerStatePath)
			condFailed := false
			switch strings.ToLower(cond.Operator) {
			case "equals":
				condFailed = stateVal != trigVal
			case "includes", "contains":
				condFailed = !strings.Contains(stateVal, trigVal)
			}
			if condFailed {
				failed = append(failed, FailedCondition{
					VisibilityCondition: cond,
					TriggerQuestionName: *cond.TriggerStatePath,
				})
			}
			continue
		}

		// Condición de pregunta — lógica existente
		trigger := findQuestion(fs, cond.TriggerQuestionID)

		var pool []models.Answer
		if trigger.RepeaterGroupID != nil {
			pool = entryAnswers
		} else {
			pool = directAnswers
		}

		var answer *models.Answer
		for i, a := range pool {
			if a.QuestionID == cond.TriggerQuestionID {
				answer = &pool[i]
				break
			}
		}

		condFailed := false
		if answer == nil {
			condFailed = true
		} else {
			switch strings.ToLower(cond.Operator) {
			case "equals":
				condFailed = answer.Value != trigVal
			case "includes", "contains":
				condFailed = !strings.Contains(answer.Value, trigVal)
			}
		}

		if condFailed {
			failed = append(failed, FailedCondition{
				VisibilityCondition: cond,
				TriggerQuestionName: trigger.Description,
			})
		}
	}

	if len(failed) > 0 {
		return VisibilityResult{Name: item.Name, Visible: false, FailedConditions: failed}
	}
	return VisibilityResult{Name: item.Name, Visible: true}
}

// ─── IsAnsweredResult ─────────────────────────────────────────────────────────

// IsAnsweredResult es el resultado de isAnsweredQuestion, isAnsweredRepeater
// e isAnsweredSection.
type IsAnsweredResult struct {
	IsVisible  bool    `json:"isVisible"`
	IsAnswered bool    `json:"isAnswered"`
	Error      *string `json:"error,omitempty"`
}

func answeredError(msg string) IsAnsweredResult {
	return IsAnsweredResult{IsVisible: true, IsAnswered: false, Error: &msg}
}

// ─── IsAnsweredRepeaterResult ─────────────────────────────────────────────────

type RepeaterAnswerError struct {
	Iteration  int    `json:"iteration"`
	QuestionID string `json:"questionId"`
	Error      string `json:"error"`
}

type IsAnsweredRepeaterResult struct {
	IsVisible  bool                  `json:"isVisible"`
	IsAnswered bool                  `json:"isAnswered"`
	Errors     []RepeaterAnswerError `json:"errors,omitempty"`
}

// ─── isAnsweredQuestion ───────────────────────────────────────────────────────

// isAnsweredQuestion evalúa si una pregunta directa (no repeater) está respondida.
//
//   q            → pregunta a evaluar
//   sectionOrder → order de la sección que la contiene (para checkVisibility)
//   answer       → respuesta encontrada en el submission, nil si no existe
//   fs           → resultado de GetFormStructure (árbol completo del form)
//   sub          → resultado de GetFormSubmission, puede ser nil
func isAnsweredQuestion(
	q            QuestionStructure,
	sectionOrder int,
	answer       *models.Answer,
	fs           *FormStructure,
	sub          *SubmissionStructure,
	formState    map[string]interface{},
) IsAnsweredResult {
	// PASO 1 — evaluar visibilidad
	vis := checkVisibility(fs, sub, nil, directQuestionItem(q, sectionOrder), formState)

	// PASO 2 — no visible: no necesita respuesta, se considera ok
	if !vis.Visible {
		return IsAnsweredResult{IsVisible: false, IsAnswered: true}
	}

	// PASO 3 — visible: validar respuesta
	answerVal := models.Answer{}
	if answer != nil {
		answerVal = *answer
	}

	result := validateAnswer(answerVal, q)
	if result.Valid {
		return IsAnsweredResult{IsVisible: true, IsAnswered: true}
	}
	return answeredError(*result.Error)
}

// ─── isAnsweredRepeater ───────────────────────────────────────────────────────

// isAnsweredRepeater evalúa si un repeater group está respondido.
//
//   r            → repeater a evaluar (incluye .Questions y .MinRepetitions)
//   sectionOrder → order de la sección que lo contiene (para checkVisibility)
//   entries      → entries ya filtradas para este repeater group
//   fs           → resultado de GetFormStructure
//   sub          → resultado de GetFormSubmission, puede ser nil
func isAnsweredRepeater(
	r            RepeaterStructure,
	sectionOrder int,
	entries      []RepeaterEntryStructure,
	fs           *FormStructure,
	sub          *SubmissionStructure,
	formState    map[string]interface{},
) IsAnsweredRepeaterResult {
	// PASO 1 — visibilidad del repeater group
	vis := checkVisibility(fs, sub, nil, repeaterItem(r, sectionOrder), formState)
	if !vis.Visible {
		return IsAnsweredRepeaterResult{IsVisible: false, IsAnswered: true}
	}

	// PASO 2 — verificar mínimo de entries
	if len(entries) < r.MinRepetitions {
		errMsg := fmt.Sprintf("Se requieren al menos %d entrada(s)", r.MinRepetitions)
		return IsAnsweredRepeaterResult{
			IsVisible:  true,
			IsAnswered: false,
			Errors:     []RepeaterAnswerError{{Error: errMsg}},
		}
	}

	// PASO 3 — validar cada entry (solo llega aquí si len >= MinRepetitions)
	var errors []RepeaterAnswerError

	for _, entry := range entries {
		for _, q := range r.Questions {
			// 3.1 — visibilidad de la pregunta dentro del scope de esta entry
			qVis := checkVisibility(fs, sub, entry.Answers, repeaterQuestionItem(q, sectionOrder, r.Order, r.ID), formState)
			if !qVis.Visible {
				continue
			}

			// 3.2 — buscar respuesta en entry.Answers
			var answerVal models.Answer
			for i := range entry.Answers {
				if entry.Answers[i].QuestionID == q.ID {
					answerVal = entry.Answers[i]
					break
				}
			}

			// 3.3 — validar respuesta
			result := validateAnswer(answerVal, q)
			if !result.Valid {
				errors = append(errors, RepeaterAnswerError{
					Iteration:  entry.Iteration,
					QuestionID: q.ID,
					Error:      *result.Error,
				})
			}
		}
	}

	// PASO 4 — retornar
	return IsAnsweredRepeaterResult{
		IsVisible:  true,
		IsAnswered: len(errors) == 0,
		Errors:     errors,
	}
}

// ─── isAnsweredSection ────────────────────────────────────────────────────────

// isAnsweredSection evalúa si una sección completa está respondida.
// Retorna IsAnsweredResult reutilizando { IsVisible, IsAnswered }.
//
//   sec → sección a evaluar (incluye .Questions y .Repeaters)
//   fs  → resultado de GetFormStructure
//   sub → resultado de GetFormSubmission, puede ser nil
func isAnsweredSection(
	sec       SectionStructure,
	fs        *FormStructure,
	sub       *SubmissionStructure,
	formState map[string]interface{},
) IsAnsweredResult {
	// PASO 1 — visibilidad de la sección
	vis := checkVisibility(fs, sub, nil, sectionItem(sec), formState)
	if !vis.Visible {
		return IsAnsweredResult{IsVisible: false, IsAnswered: true}
	}

	// PASO 2 — construir índices desde sub (vacíos si sub == nil)
	answerByQuestion := map[string]models.Answer{}
	entriesByGroup   := map[string][]RepeaterEntryStructure{}

	if sub != nil {
		for _, a := range sub.DirectAnswers {
			answerByQuestion[a.QuestionID] = a
		}
		for _, e := range sub.RepeaterEntries {
			entriesByGroup[e.RepeaterGroupID] = append(entriesByGroup[e.RepeaterGroupID], e)
		}
	}

	// PASO 3 — iterar preguntas directas
	for _, q := range sec.Questions {
		var answer *models.Answer
		if a, ok := answerByQuestion[q.ID]; ok {
			answer = &a
		}
		result := isAnsweredQuestion(q, sec.Order, answer, fs, sub, formState)
		if !result.IsAnswered {
			return IsAnsweredResult{IsVisible: true, IsAnswered: false}
		}
	}

	// PASO 4 — iterar repeaters
	for _, r := range sec.Repeaters {
		entries := entriesByGroup[r.ID]
		result  := isAnsweredRepeater(r, sec.Order, entries, fs, sub, formState)
		if !result.IsAnswered {
			return IsAnsweredResult{IsVisible: true, IsAnswered: false}
		}
	}

	// PASO 5 — todo respondido
	return IsAnsweredResult{IsVisible: true, IsAnswered: true}
}

// ─── loadForm ─────────────────────────────────────────────────────────────────

// LoadFormResult es la respuesta de loadForm.
type LoadFormResult struct {
	FormStructure  *FormStructure       `json:"formStructure"`
	FormSubmission *SubmissionStructure `json:"formSubmission,omitempty"`
	CurrentSection *SectionStructure    `json:"currentSection"`
}

// loadForm carga la estructura del formulario, la submission si existe, y
// determina la currentSection: primera sección visible no respondida.
// Si todas están respondidas, currentSection = última sección visible.
func (s *formService) LoadForm(ctx context.Context, formID, submissionID string, formState map[string]interface{}) (*LoadFormResult, error) {
	// 1. Cargar estructura del formulario
	fs, err := s.GetFormStructure(ctx, formID)
	if err != nil {
		return nil, err
	}

	// 2. Cargar submission si se proporcionó
	var sub *SubmissionStructure
	if submissionID != "" {
		sub, err = s.GetFormSubmission(ctx, submissionID)
		if err != nil {
			return nil, err
		}
	}

	// 3. Encontrar currentSection: primera visible no respondida
	var currentSection *SectionStructure
	var lastVisible *SectionStructure

	for i := range fs.Sections {
		sec := &fs.Sections[i]
		result := isAnsweredSection(*sec, fs, sub, formState)
		sec.IsAnswered = result.IsAnswered
		sec.IsVisible  = result.IsVisible
		if !result.IsVisible {
			continue
		}
		lastVisible = sec
		if !result.IsAnswered && currentSection == nil {
			currentSection = sec
		}
	}

	// 4. Si todas respondidas → currentSection = última visible
	if currentSection == nil {
		currentSection = lastVisible
	}

	return &LoadFormResult{
		FormStructure:  fs,
		FormSubmission: sub,
		CurrentSection: currentSection,
	}, nil
}

// ─── saveSection ─────────────────────────────────────────────────────────────

type SaveAnswerInput struct {
	QuestionID string `json:"questionId"`
	Value      string `json:"value"`
}

type SaveRepeaterEntryInput struct {
	ID              string            `json:"id"`
	RepeaterGroupID string            `json:"repeaterGroupId"`
	Iteration       int               `json:"iteration"`
	IsTemp          bool              `json:"isTemp"`
	Answers         []SaveAnswerInput  `json:"answers"`
}

type SaveSectionInput struct {
	FormID           string                   `json:"formId"`
	FormSectionID    string                   `json:"formSectionId"`
	FormSubmissionID string                   `json:"formSubmissionId"` // vacío = crear nuevo
	DirectAnswers    []SaveAnswerInput         `json:"directAnswers"`
	RepeaterEntries  []SaveRepeaterEntryInput  `json:"repeaterEntries"`
	ActorID          string                   `json:"actorId"`          // general_user_i_code — inyectado por el controller desde la sesión
	FormState        map[string]interface{}   `json:"formState"`        // estado externo del padre para condiciones de visibilidad
}

func (s *formService) SaveSection(ctx context.Context, input SaveSectionInput) (*LoadFormResult, error) {
	// 1. Resolver FormSubmission
	submissionID := input.FormSubmissionID
	if submissionID == "" {
		fs := &models.FormSubmission{FormID: input.FormID}
		if err := s.submissionRepo.Create(ctx, fs); err != nil {
			return nil, fmt.Errorf("saveSection: crear submission: %w", err)
		}
		submissionID = fs.ID
	}

	// 2. Construir mapa de valores válidos por pregunta (single, dropdown, multiple)
	validOptionValues, err := s.buildValidOptionValues(ctx, input.FormSectionID)
	if err != nil {
		return nil, fmt.Errorf("saveSection: cargar opciones válidas: %w", err)
	}
	input.DirectAnswers   = filterAnswers(input.DirectAnswers, validOptionValues)
	input.RepeaterEntries = filterRepeaterAnswers(input.RepeaterEntries, validOptionValues)

	// 4. Procesar directAnswers
	existingDirect, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("saveSection: leer directAnswers: %w", err)
	}
	byQuestion := map[string][]models.Answer{}
	for _, a := range existingDirect {
		byQuestion[a.QuestionID] = append(byQuestion[a.QuestionID], a)
	}

	for _, da := range input.DirectAnswers {
		existing := byQuestion[da.QuestionID]
		if len(existing) == 0 {
			ans := &models.Answer{FormSubmissionID: submissionID, QuestionID: da.QuestionID, Value: da.Value}
			if err := s.answerRepo.Create(ctx, ans); err != nil {
				return nil, fmt.Errorf("saveSection: crear answer %s: %w", da.QuestionID, err)
			}
		} else {
			// Eliminar duplicados (dejar solo el primero)
			for i := 1; i < len(existing); i++ {
				if err := s.answerRepo.Delete(ctx, existing[i].ID); err != nil {
					return nil, fmt.Errorf("saveSection: eliminar answer duplicada: %w", err)
				}
			}
			// Actualizar el que queda
			existing[0].Value = da.Value
			if err := s.answerRepo.Update(ctx, &existing[0]); err != nil {
				return nil, fmt.Errorf("saveSection: actualizar answer %s: %w", da.QuestionID, err)
			}
		}
	}

	// 3. Procesar repeaterEntries
	// createdEntryIDs: IDs reales de entries recién creadas (para protegerlas en el paso 4)
	createdEntryIDs := map[string]bool{}
	for _, entry := range input.RepeaterEntries {
		if entry.IsTemp {
			// Crear nueva entry
			newEntry := &models.RepeaterEntry{
				FormSubmissionID: submissionID,
				RepeaterGroupID:  entry.RepeaterGroupID,
				Iteration:        entry.Iteration,
			}
			if err := s.repeaterEntryRepo.Create(ctx, newEntry); err != nil {
				return nil, fmt.Errorf("saveSection: crear repeaterEntry: %w", err)
			}
			createdEntryIDs[newEntry.ID] = true
			for _, ans := range entry.Answers {
				entryID := newEntry.ID
				a := &models.Answer{FormSubmissionID: submissionID, QuestionID: ans.QuestionID, Value: ans.Value, RepeaterEntryID: &entryID}
				if err := s.answerRepo.Create(ctx, a); err != nil {
					return nil, fmt.Errorf("saveSection: crear answer repeater: %w", err)
				}
			}
		} else {
			// Borrar answers existentes y recrear (clean upsert)
			existingAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
			if err != nil {
				return nil, fmt.Errorf("saveSection: leer answers de entry %s: %w", entry.ID, err)
			}
			for _, a := range existingAnswers {
				if err := s.answerRepo.Delete(ctx, a.ID); err != nil {
					return nil, fmt.Errorf("saveSection: eliminar answer de entry: %w", err)
				}
			}
			for _, ans := range entry.Answers {
				entryID := entry.ID
				a := &models.Answer{FormSubmissionID: submissionID, QuestionID: ans.QuestionID, Value: ans.Value, RepeaterEntryID: &entryID}
				if err := s.answerRepo.Create(ctx, a); err != nil {
					return nil, fmt.Errorf("saveSection: crear answer repeater existente: %w", err)
				}
			}
		}
	}

	// 4. Eliminar repeaterEntries obsoletas (existían en BD pero no llegaron en el request)
	groups, err := s.repeaterGroupRepo.FindBySectionID(ctx, input.FormSectionID)
	if err != nil {
		return nil, fmt.Errorf("saveSection: leer groups de sección: %w", err)
	}
	if len(groups) > 0 {
		groupIDs := make([]string, len(groups))
		for i, g := range groups {
			groupIDs[i] = g.ID
		}
		existingEntries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, groupIDs)
		if err != nil {
			return nil, fmt.Errorf("saveSection: leer entries existentes: %w", err)
		}
		// Proteger: entries del request (no-temp) + entries recién creadas en este mismo save
		requestIDs := map[string]bool{}
		for _, e := range input.RepeaterEntries {
			if !e.IsTemp {
				requestIDs[e.ID] = true
			}
		}
		for id := range createdEntryIDs {
			requestIDs[id] = true
		}
		for _, entry := range existingEntries {
			if requestIDs[entry.ID] {
				continue
			}
			// Eliminar answers de la entry obsoleta
			answers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
			if err != nil {
				return nil, fmt.Errorf("saveSection: leer answers de entry obsoleta: %w", err)
			}
			for _, a := range answers {
				if err := s.answerRepo.Delete(ctx, a.ID); err != nil {
					return nil, fmt.Errorf("saveSection: eliminar answer de entry obsoleta: %w", err)
				}
			}
			if err := s.repeaterEntryRepo.Delete(ctx, entry.ID); err != nil {
				return nil, fmt.Errorf("saveSection: eliminar entry obsoleta: %w", err)
			}
		}
	}

	// Hook genérico: se dispara SIEMPRE después de guardar las respuestas de
	// una sección (no solo la primera ni la última) -- análogo a
	// OnEndFormSubmission pero por cada saveSection en vez de solo al
	// completar. Síncrono (no goroutine): los guardados posteriores asumen
	// que cualquier efecto de esta sección ya se aplicó.
	if err := s.OnSectionUpdate(ctx, input.FormID, submissionID, input.ActorID); err != nil {
		log.Printf("[saveSection] OnSectionUpdate error para formId=%s submission=%s: %v", input.FormID, submissionID, err)
	}

	// 5. Retornar LoadForm con el estado actualizado (isAnswered/isVisible por sección)
	result, err := s.LoadForm(ctx, input.FormID, submissionID, input.FormState)
	if err != nil {
		return nil, err
	}

	// 6. Si todas las secciones visibles están respondidas → onEndFormSubmission
	allAnswered := true
	var pendingSections []string
	for _, sec := range result.FormStructure.Sections {
		if sec.IsVisible && !sec.IsAnswered {
			allAnswered = false
			pendingSections = append(pendingSections, fmt.Sprintf("%q (order=%d)", sec.Name, sec.Order))
		}
	}
	if allAnswered {
		log.Printf("[saveSection] submissionId=%s → formulario COMPLETO, disparando OnEndFormSubmission", submissionID)
		actorID := input.ActorID
		if err := s.OnEndFormSubmission(ctx, input.FormID, submissionID, actorID); err != nil {
			log.Printf("[OnEndFormSubmission] error: %v\n", err)
		}
	} else {
		log.Printf("[saveSection] submissionId=%s → formulario INCOMPLETO — secciones pendientes: %v", submissionID, pendingSections)
	}

	return result, nil
}

// ─── OnSectionUpdate ────────────────────────────────────────────────────────────

// OnSectionUpdate es llamado después de CADA guardado de sección (completo o
// no) -- a diferencia de OnEndFormSubmission, que solo se dispara cuando el
// formulario completo queda respondido. Delega a la función específica según
// el formID, igual que OnEndFormSubmission. Los formularios sin lógica
// registrada no hacen nada (no-op).
func (s *formService) OnSectionUpdate(ctx context.Context, formID, submissionID, actorID string) error {
	switch formID {
	case RegistroCasoFormID:
		if s.victimCaseFormSvc == nil {
			return fmt.Errorf("victimCaseFormSvc no inyectado")
		}
		return s.victimCaseFormSvc.UpdateCaseDraft(ctx, submissionID, actorID)
	default:
		return nil
	}
}

// ─── OnEndFormSubmission ──────────────────────────────────────────────────────

// IDs de formularios con lógica de end-submission.
const (
	// SeguimientoFormID es el formulario "vivo" de seguimiento; exportado para que
	// otros paquetes (p.ej. el controller) puedan validar permisos por formId.
	SeguimientoFormID   = "2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff"
	barrierUpdateFormID = "4d0aeb46-5af3-4c47-a0d5-c5bfc4d549ff"
	cierreCasoFormID    = "da8423ab-1a8c-47db-96b7-d10496df571a"
)

// OnEndFormSubmission es llamado cuando todas las secciones visibles han sido respondidas.
// Delega a la función específica según el formID.
func (s *formService) OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error {
	switch formID {
	case SeguimientoFormID:
		return s.processFollowUpSubmission(ctx, submissionID, actorID)
	case RegistroCasoFormID:
		if s.victimCaseFormSvc == nil {
			return fmt.Errorf("victimCaseFormSvc no inyectado")
		}
		return s.victimCaseFormSvc.Activate(ctx, submissionID, actorID)
	case barrierUpdateFormID:
		return s.processBarrierUpdateSubmission(ctx, submissionID)
	case cierreCasoFormID:
		return s.processCaseClosureSubmission(ctx, submissionID, actorID)
	case constants.FormIDPrimerContacto, constants.FormIDPrimeraAtencion,
		constants.FormIDAtencionPsicosocial, constants.FormIDCierre:
		return s.processPsicosocialSessionSubmission(ctx, formID, submissionID, actorID)
	}
	return nil
}

// ─── Formularios Psicosociales — evento E-02 (guardado) ──────────────────────
//
// Ver DocsMD/Screens/psicosocial-sesion/Flujos/flow-E02-cuando-se-guarda-formulario.md
// para el detalle de decisiones (IDs de pregunta reales capturados de Supabase, tabla de
// transición de estados, y preguntas resueltas con el líder).

// resolvePsicosocialSessionType determina el session_type resultante de completar uno de los
// 4 formularios psicosociales, según formID y las respuestas directas del submission.
// El segundo valor retornado es la "fecha nueva" (string YYYY-MM-DD) cuando sessionType ==
// SessionTypeContactoSinAtencion; vacío en cualquier otro caso.
func resolvePsicosocialSessionType(formID string, answerMap map[string]string) (sessionType string, fechaNueva string) {
	const (
		// Form Primer Contacto
		qPCContinuarPA    = "aa46351a-f6c0-4666-af3e-e0f1f7b6dc8c" // S1 — Continuar Primera Atención (boolean)
		qPCConsentimiento = "c3296c87-d7e3-4ea1-8a0f-cbf77c561d0f" // S4 — Consentimiento Informado (si/no)

		// Form Primera Atención
		qPAEsAtencion     = "3117fd01-6a31-4595-9888-dcdcc96d84e7" // S1 — ¿Es atención o solo contacto?
		qPAConsentimiento = "7256b91e-861b-48b3-9216-2ac13f0ae889" // S4 — Consentimiento Informado (si/no)
		qPAFechaNueva     = "64d63b79-edee-464b-be56-1104efd31a46" // S1 — Fecha nueva

		// Form Atención Psicosocial (antes "Seguimiento")
		qSEGEsAtencion = "16de7674-4349-4acc-9d0d-656d1d5e5f10" // S1 — ¿Es atención o solo contacto?
		qSEGFechaNueva = "72ce49d2-f853-4f4a-9f1a-f795f4d514c3" // S1 — Fecha nueva

		// Form Cierre
		qCIEEsAtencion     = "be185336-e025-4768-a1a2-3ab9e2281450" // S1 — ¿Es atención o solo contacto?
		qCIEFechaNueva     = "1b4d09f0-e5ca-4b4e-9d74-428478a113c6" // S1 — Fecha nueva
		qCIECerrarRemision = "2edd39af-45da-4d24-ac67-fdc8b9b0b6ac" // S4 — Cerrar remisión (boolean)
	)

	switch formID {
	case constants.FormIDPrimerContacto:
		if answerMap[qPCContinuarPA] != "true" {
			return models.SessionTypePrimerContacto, ""
		}
		if answerMap[qPCConsentimiento] == "si" {
			return models.SessionTypePrimerContactoConAtencion, ""
		}
		// Decisión del líder (Jul 2026): sin consentimiento avanza igual que un primer
		// contacto regular — no queda "atrapada" en en_devolucion.
		return models.SessionTypePrimerContactoSinConsentimiento, ""

	case constants.FormIDPrimeraAtencion:
		if answerMap[qPAEsAtencion] == "solo_contacto" {
			return models.SessionTypeContactoSinAtencion, answerMap[qPAFechaNueva]
		}
		if answerMap[qPAConsentimiento] == "si" {
			return models.SessionTypePrimeraAtencion, ""
		}
		return models.SessionTypeCierreNoConsentimiento, ""

	case constants.FormIDAtencionPsicosocial:
		if answerMap[qSEGEsAtencion] == "solo_contacto" {
			return models.SessionTypeContactoSinAtencion, answerMap[qSEGFechaNueva]
		}
		return models.SessionTypeAtencionPsicosocial, ""

	case constants.FormIDCierre:
		if answerMap[qCIEEsAtencion] == "solo_contacto" {
			return models.SessionTypeContactoSinAtencion, answerMap[qCIEFechaNueva]
		}
		if answerMap[qCIECerrarRemision] == "true" {
			return models.SessionTypeCierre, ""
		}
		// "Cerrar remisión = No" → actúa como un seguimiento más dentro del Form de Cierre.
		return models.SessionTypeAtencionPsicosocial, ""
	}
	return "", ""
}

// psicosocialSessionTimelineDescription construye el tipo de evento de timeline según el
// sessionType resultante.
func psicosocialSessionTimelineType(sessionType string) string {
	switch sessionType {
	case models.SessionTypeCierre:
		return "Remisión Psicosocial Cerrada"
	case models.SessionTypeCierreNoConsentimiento:
		return "Cierre sin Consentimiento"
	case models.SessionTypeContactoSinAtencion:
		return "Contacto sin Atención"
	default:
		return "Sesión Psicosocial Realizada"
	}
}

// psicosocialAgendaAnswers agrupa las respuestas de agenda condicional (Aug 2026).
// ShouldSchedule es true solo si "¿Agendar nueva sesión?" = sí/true Y hay fecha Y hora.
type psicosocialAgendaAnswers struct {
	ShouldSchedule bool
	Fecha          string
	Hora           string
	Mode           string // "individual" | "dupla"
}

// extractAgendaAnswers busca Agendar / Fecha / Hora entre candidatos S1 y S4 del formulario.
// Reemplaza extractFechaProximaAtencion: la fecha sola ya no agenda; hace falta Agendar=Sí + hora.
func extractAgendaAnswers(formID string, answerMap map[string]string) psicosocialAgendaAnswers {
	type slot struct {
		agendar, fecha, hora string
	}
	var slots []slot
	switch formID {
	case constants.FormIDPrimerContacto:
		slots = []slot{
			{constants.QPCAgendarNuevaSesionS1, constants.QPCFechaProximaS1, constants.QPCHoraProximaAtencionS1},
			{constants.QPCAgendarNuevaSesionS4, constants.QPCFechaProximaS4, constants.QPCHoraProximaAtencionS4},
		}
	case constants.FormIDPrimeraAtencion:
		slots = []slot{
			{constants.QPAAgendarNuevaSesionS1, constants.QPAFechaNuevaS1, constants.QPAHoraProximaAtencionS1},
			{constants.QPAAgendarNuevaSesionS4, constants.QPAFechaProximaS4, constants.QPAHoraProximaAtencionS4},
		}
	case constants.FormIDAtencionPsicosocial:
		slots = []slot{
			{constants.QSEGAgendarNuevaSesionS1, constants.QSEGFechaNuevaS1, constants.QSEGHoraProximaAtencionS1},
			{constants.QSEGAgendarNuevaSesionS4, constants.QSEGFechaProximaS4, constants.QSEGHoraProximaAtencionS4},
		}
	case constants.FormIDCierre:
		slots = []slot{
			{constants.QCIEAgendarNuevaSesionS1, constants.QCIEFechaNuevaS1, constants.QCIEHoraProximaAtencionS1},
			{constants.QCIEAgendarNuevaSesionS4, constants.QCIEFechaProximaS4, constants.QCIEHoraProximaAtencionS4},
		}
	}

	for _, sl := range slots {
		agendar := strings.TrimSpace(answerMap[sl.agendar])
		fecha := strings.TrimSpace(answerMap[sl.fecha])
		hora := strings.TrimSpace(answerMap[sl.hora])
		if answerIsTrueish(agendar) && fecha != "" && hora != "" {
			return psicosocialAgendaAnswers{ShouldSchedule: true, Fecha: fecha, Hora: hora}
		}
	}
	return psicosocialAgendaAnswers{}
}

// resolvePsicosocialAgendaMode lee "¿La atención es individual o en dupla?" (preferencia)
// o infiere de ps.DuplaID / ProfessionalID.
func resolvePsicosocialAgendaMode(answerMap map[string]string, questions []models.Question, ps *models.PsychosocialSupport) string {
	for _, q := range questions {
		if q.RepeaterGroupID != nil {
			continue
		}
		desc := strings.ToLower(q.Description)
		if strings.Contains(desc, "individual o en dupla") || strings.Contains(desc, "individual o en dúpla") {
			v := strings.ToLower(strings.TrimSpace(answerMap[q.ID]))
			if strings.Contains(v, "dupla") || strings.Contains(v, "dúpla") {
				return "dupla"
			}
			if strings.Contains(v, "individual") {
				return "individual"
			}
		}
	}
	if ps.DuplaID != nil && *ps.DuplaID != "" {
		return "dupla"
	}
	return "individual"
}

// scheduleNextPsicosocialContact crea un nuevo team_contact PENDIENTE (is_completed = false)
// para la fecha/hora indicadas. Antes de crear valida disponibilidad (ventana 2h); si no hay
// cupo retorna nil sin error (solo log warn) para no fallar el guardado del formulario.
func (s *formService) scheduleNextPsicosocialContact(ctx context.Context, ps *models.PsychosocialSupport, fechaStr, horaStr, mode string) error {
	fecha, err := parseAgendaDate(fechaStr)
	if err != nil {
		return fmt.Errorf("fecha próxima atención no parseable %q: %w", fechaStr, err)
	}
	hora := normalizeScheduledTime(horaStr)

	available, msg, err := s.checkPsicosocialSessionAvailability(ctx, ps, fecha, hora, mode)
	if err != nil {
		log.Printf("[scheduleNextPsicosocialContact] advertencia: error chequeando disponibilidad: %v", err)
	} else if !available {
		log.Printf("[scheduleNextPsicosocialContact] horario no disponible psicosocialId=%s fecha=%s hora=%s — %s (no se agenda)",
			ps.ID, fechaStr, hora, msg)
		return nil
	}

	status := "agendada"
	next := &models.TeamContact{
		CaseID:         ps.CaseID,
		PsicosocialID:  &ps.ID,
		IsPsicoSession: true,
		IsCompleted:    false,
		Status:         &status,
		ScheduledDate:  &fecha,
		ScheduledTime:  &hora,
	}
	switch {
	case mode == "dupla" && ps.DuplaID != nil && *ps.DuplaID != "":
		next.DuplaID = ps.DuplaID
	case ps.DuplaID != nil && *ps.DuplaID != "" && mode != "individual":
		next.DuplaID = ps.DuplaID
	case ps.ProfessionalID != nil && *ps.ProfessionalID != "":
		next.ProfessionalID = ps.ProfessionalID
	}

	if err := s.teamContactRepo.Create(ctx, next); err != nil {
		return fmt.Errorf("crear team_contact agendado: %w", err)
	}
	log.Printf("[scheduleNextPsicosocialContact] psicosocialId=%s fecha=%s hora=%s mode=%s nuevo team_contact=%s",
		ps.ID, fechaStr, hora, mode, next.ID)
	return nil
}

// checkPsicosocialSessionAvailability valida solape de 2h para el profesional o la dupla.
func (s *formService) checkPsicosocialSessionAvailability(ctx context.Context, ps *models.PsychosocialSupport, date time.Time, timeStr, mode string) (bool, string, error) {
	if s.teamContactRepo == nil {
		return true, "Horario disponible", nil
	}
	contacts, err := s.teamContactRepo.FindScheduledOnDate(ctx, date)
	if err != nil {
		return false, "", err
	}
	var psychID, swID string
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		if ps.DuplaID != nil && *ps.DuplaID != "" {
			mode = "dupla"
		} else {
			mode = "individual"
		}
	}
	if mode == "dupla" && ps.DuplaID != nil && *ps.DuplaID != "" && s.duplaRepo != nil {
		if d, err := s.duplaRepo.FindEnrichedByID(ctx, *ps.DuplaID); err == nil && d != nil {
			psychID = d.PsychologistID
			swID = d.SocialWorkerID
		}
	}
	ok, msg := evaluatePsicosocialAvailability(contacts, ps, timeStr, mode, psychID, swID)
	return ok, msg, nil
}

// processPsicosocialSessionSubmission implementa el evento E-02 para los 4 formularios
// psicosociales: resuelve el team_contact asociado al submission, determina el session_type
// según las respuestas, actualiza team_contact + psychosocial_support + timeline, crea un
// BarrierV2 (con sus tareas/oficios derivados) por cada entrada de "Identificación de Barreras"
// y procesa "Seguimiento a Barreras" (§9.6) — ver psicosocial_barreras.go.
func (s *formService) processPsicosocialSessionSubmission(ctx context.Context, formID, submissionID, actorID string) error {
	if s.teamContactRepo == nil || s.psychosocialSupportRepo == nil {
		log.Printf("⚠️  [processPsicosocialSessionSubmission] teamContactRepo/psychosocialSupportRepo es nil — inyectar TeamContactRepo/PsychosocialSupportRepo en FormServiceDeps (main.go)")
		return nil
	}

	tc, err := s.teamContactRepo.FindByFormSubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processPsicosocialSessionSubmission: buscar team_contact: %w", err)
	}

	actorName := actorID
	if s.agentLightRepo != nil && actorID != "" {
		if agent, err := s.agentLightRepo.FindByICode(ctx, actorID); err == nil && agent != nil {
			actorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
		}
	}

	psID := ""
	if tc.PsicosocialID != nil {
		psID = *tc.PsicosocialID
	}

	// Idempotente — si ya fue procesado, solo registrar la edición y salir.
	if tc.IsCompleted {
		if s.caseTimelineRepo == nil {
			log.Printf("⚠️  [processPsicosocialSessionSubmission] caseTimelineRepo es nil — evento 'Sesión Editada' NO creado para team_contact=%s", tc.ID)
			return nil
		}
		now := time.Now()
		if err := s.caseTimelineRepo.Create(ctx, &models.CaseTimelineEvent{
			CaseID:                tc.CaseID,
			Category:              models.TimelineCategoryPsicosocial,
			Type:                  "Sesión Editada",
			Icon:                  models.TimelineIconPospuesto,
			Color:                 models.TimelineColorTeal,
			Date:                  now,
			Description:           "Sesión psicosocial editada después de su ejecución",
			EventUserID:           actorID,
			ActorName:             actorName,
			PsychosocialSupportID: psID,
			CreatedAt:             now,
		}); err != nil {
			log.Printf("[processPsicosocialSessionSubmission] advertencia: no se pudo crear evento 'Sesión Editada': %v", err)
		}
		return nil
	}

	if psID == "" {
		return fmt.Errorf("processPsicosocialSessionSubmission: team_contact %s sin psicosocial_id", tc.ID)
	}
	ps, err := s.psychosocialSupportRepo.FindByID(ctx, psID)
	if err != nil {
		return fmt.Errorf("processPsicosocialSessionSubmission: buscar psychosocial_support: %w", err)
	}

	answers, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processPsicosocialSessionSubmission: leer respuestas: %w", err)
	}
	answerMap := make(map[string]string, len(answers))
	for _, a := range answers {
		answerMap[a.QuestionID] = a.Value
	}

	sessionType, fechaNueva := resolvePsicosocialSessionType(formID, answerMap)
	if sessionType == "" {
		return fmt.Errorf("processPsicosocialSessionSubmission: no se pudo determinar sessionType para formID=%s", formID)
	}

	// ── Actualizar team_contact ──────────────────────────────────────────────
	now := time.Now()
	tc.SessionType = &sessionType
	tc.IsCompleted = true
	tc.CompletedAt = &now
	if sessionType == models.SessionTypeContactoSinAtencion {
		tc.IsPsicoSession = false
	}
	if err := s.teamContactRepo.Update(ctx, tc); err != nil {
		return fmt.Errorf("processPsicosocialSessionSubmission: actualizar team_contact: %w", err)
	}

	// ── Actualizar psychosocial_support según sessionType ────────────────────
	switch sessionType {
	case models.SessionTypePrimerContacto, models.SessionTypePrimerContactoSinConsentimiento:
		ps.YaHizoPrimerContacto = true
		ps.Status = models.PsychosocialSupportStatusEnGestion
	case models.SessionTypePrimerContactoConAtencion:
		ps.YaHizoPrimerContacto = true
		ps.YaHizoPrimeraAtencion = true
		ps.SessionCount++
		ps.Status = models.PsychosocialSupportStatusEnGestion
	case models.SessionTypePrimeraAtencion:
		ps.YaHizoPrimeraAtencion = true
		ps.SessionCount++
		ps.Status = models.PsychosocialSupportStatusEnGestion
	case models.SessionTypeAtencionPsicosocial:
		ps.SessionCount++
	case models.SessionTypeCierre:
		ps.SessionCount++
		ps.Status = models.PsychosocialSupportStatusCerrado
	case models.SessionTypeCierreNoConsentimiento:
		ps.Status = models.PsychosocialSupportStatusEnDevolucion
	case models.SessionTypeContactoSinAtencion:
		if fechaNueva != "" {
			if t, err := time.Parse("2006-01-02", fechaNueva); err == nil {
				ps.ScheduledAt = &t
			} else {
				log.Printf("[processPsicosocialSessionSubmission] advertencia: fecha nueva no parseable %q: %v", fechaNueva, err)
			}
		}
	}

	if err := s.psychosocialSupportRepo.Update(ctx, ps); err != nil {
		return fmt.Errorf("processPsicosocialSessionSubmission: actualizar psychosocial_support: %w", err)
	}

	// ── Agendar próxima sesión (condicional: Agendar=Sí + fecha + hora + disponibilidad) ──
	// No aplica si la remisión se está cerrando en esta misma sesión (sessionType == CIERRE).
	if sessionType != models.SessionTypeCierre {
		agenda := extractAgendaAnswers(formID, answerMap)
		var formQuestions []models.Question
		if s.questionRepo != nil {
			if qs, err := s.questionRepo.FindByFormID(ctx, formID); err == nil {
				formQuestions = qs
			} else {
				log.Printf("[processPsicosocialSessionSubmission] advertencia: FindByFormID: %v", err)
			}
		}
		agenda.Mode = resolvePsicosocialAgendaMode(answerMap, formQuestions, ps)
		if agenda.ShouldSchedule {
			if err := s.scheduleNextPsicosocialContact(ctx, ps, agenda.Fecha, agenda.Hora, agenda.Mode); err != nil {
				log.Printf("[processPsicosocialSessionSubmission] advertencia: no se pudo agendar próxima sesión: %v", err)
			}
		}
	}

	// ── Barreras: crear barrier_v2 + case_task/entity_letter por cada entrada del repeater
	// "Identificación de Barreras" ────────────────────────────────────────────────────────
	barrierCount, err := s.processPsicosocialBarrierEntries(ctx, formID, submissionID, actorID, tc.ID, ps)
	if err != nil {
		log.Printf("[processPsicosocialSessionSubmission] advertencia: error procesando barreras: %v", err)
	}

	// ── Seguimiento a Barreras (§9.6): actualizar/cerrar barreras activas ────────────────
	if err := s.processPsicosocialBarrierFollowUps(ctx, formID, submissionID, actorID, ps); err != nil {
		log.Printf("[processPsicosocialSessionSubmission] advertencia: error en seguimiento a barreras: %v", err)
	}

	// ── Timeline de sesión ──────────────────────────────────────────────────────────────
	timelineDesc := "Sesión psicosocial registrada — tipo: " + sessionType
	if barrierCount > 0 {
		timelineDesc = fmt.Sprintf("%s | %d barrera(s) identificada(s)", timelineDesc, barrierCount)
	}
	if s.caseTimelineRepo != nil {
		if err := s.caseTimelineRepo.Create(ctx, &models.CaseTimelineEvent{
			CaseID:                tc.CaseID,
			Category:              models.TimelineCategoryPsicosocial,
			Type:                  psicosocialSessionTimelineType(sessionType),
			Icon:                  models.TimelineIconSeguimiento,
			Color:                 models.TimelineColorGreen,
			Date:                  now,
			Description:           timelineDesc,
			EventUserID:           actorID,
			ActorName:             actorName,
			PsychosocialSupportID: ps.ID,
			CreatedAt:             now,
		}); err != nil {
			log.Printf("[processPsicosocialSessionSubmission] advertencia: no se pudo crear evento de timeline: %v", err)
		}
	} else {
		log.Printf("⚠️  [processPsicosocialSessionSubmission] caseTimelineRepo es nil — evento de timeline NO creado para psicosocialId=%s", ps.ID)
	}

	// ── Timeline hechos de violencia (paridad PASO 7b de processFollowUpSubmission) ─────
	if err := s.processPsicosocialHechosViolenciaTimeline(ctx, formID, answerMap, ps, actorID); err != nil {
		log.Printf("[processPsicosocialSessionSubmission] advertencia: hechos de violencia: %v", err)
	}

	log.Printf("[processPsicosocialSessionSubmission] psicosocialId=%s sessionType=%s sessionCount=%d status=%s ya_hizo_pc=%v ya_hizo_pa=%v",
		ps.ID, sessionType, ps.SessionCount, ps.Status, ps.YaHizoPrimerContacto, ps.YaHizoPrimeraAtencion)

	return nil
}

// processPsicosocialHechosViolenciaTimeline crea un CaseTimelineEvent TimelineTypeHechosCaso
// si "Hay nuevos hechos de violencia" = true. Resuelve preguntas por description (no IDs fijos),
// preferiendo las que no pertenecen a un repeater (secciones de contacto, no S5 cierre).
func (s *formService) processPsicosocialHechosViolenciaTimeline(ctx context.Context, formID string, answerMap map[string]string, ps *models.PsychosocialSupport, actorID string) error {
	if s.caseTimelineRepo == nil || s.questionRepo == nil {
		return nil
	}
	questions, err := s.questionRepo.FindByFormID(ctx, formID)
	if err != nil {
		return err
	}

	findFirst := func(substr string) string {
		var fallback string
		for _, q := range questions {
			if !strings.Contains(q.Description, substr) {
				continue
			}
			if q.RepeaterGroupID == nil {
				return q.ID
			}
			if fallback == "" {
				fallback = q.ID
			}
		}
		return fallback
	}

	qHechos := findFirst("Hay nuevos hechos de violencia")
	qDesc := findFirst("Descripción de los hechos")
	qFecha := findFirst("Fecha (de los hechos)")
	if qFecha == "" {
		qFecha = findFirst("Fecha de los hechos")
	}

	if qHechos == "" || !answerIsTrueish(answerMap[qHechos]) {
		return nil
	}

	descripcion := ""
	if qDesc != "" {
		descripcion = strings.TrimSpace(answerMap[qDesc])
	}
	fechaStr := ""
	if qFecha != "" {
		fechaStr = strings.TrimSpace(answerMap[qFecha])
	}

	fechaHechos := time.Now()
	if fechaStr != "" {
		if t, err := time.Parse("2006-01-02", fechaStr); err == nil {
			fechaHechos = t
		}
	}

	hechoEvent := &models.CaseTimelineEvent{
		CaseID:                ps.CaseID,
		FollowUpID:            ps.FollowUpID,
		Category:              models.TimelineCategoryGeneral,
		Type:                  models.TimelineTypeHechosCaso,
		Icon:                  models.TimelineIconHechosCaso,
		Color:                 "#f87171",
		Description:           descripcion,
		EventUserID:           actorID,
		Date:                  fechaHechos,
		PsychosocialSupportID: ps.ID,
		CreatedAt:             time.Now(),
	}
	if err := s.caseTimelineRepo.Create(ctx, hechoEvent); err != nil {
		return err
	}
	log.Printf("[processPsicosocialSessionSubmission] ✅ evento 'Nuevos hechos de violencia' creado → fecha=%s", fechaStr)
	return nil
}

// processCaseClosureSubmission marca el seguimiento como REALIZADO y si se confirma el cierre
// del caso, cambia el estado del caso a cerrado, cierra todos los seguimientos pendientes/reprogramados,
// y registra el evento de cierre en el timeline del caso.
func (s *formService) processCaseClosureSubmission(ctx context.Context, submissionID, actorID string) error {
	// 1. Buscar el follow_up asociado al submission
	fu, err := s.followUpRepo.FindByFormSubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processCaseClosureSubmission: buscar followUp: %w", err)
	}

	// Si ya fue procesado → no re-procesar
	if fu.Status == models.FollowUpStatusRealizado {
		return nil
	}

	// 2. Marcar el follow_up como REALIZADO y registrar completed_at
	if err := s.followUpRepo.UpdateStatus(ctx, fu.ID, models.FollowUpStatusRealizado); err != nil {
		return fmt.Errorf("processCaseClosureSubmission: actualizar estado followUp: %w", err)
	}

	// 3. Leer las respuestas directas del submission
	answers, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processCaseClosureSubmission: leer respuestas: %w", err)
	}

	answerMap := make(map[string]string, len(answers))
	for _, a := range answers {
		answerMap[a.QuestionID] = a.Value
	}

	const (
		qMotivoCierre       = "d2c6e1af-651f-43b4-8e61-aecafd07443d" // Motivo del cierre (single)
		qCausaCierre        = "90375500-a316-4cf5-b7ec-f5c402da92c2" // Describa la causa del cierre (text)
		qAccionesCierre     = "4a7d0110-b06d-4569-9572-cd0a5e9ef2c1" // ¿Realizó acciones institucionales? (boolean)
		qRutaAtencionPrevia = "f4b162fd-adf5-4341-ab4f-162fdadf5341" // ¿Se activó la ruta de atención o se generó algún oficio previamente en este caso? (boolean)
	)

	// 4. Cambiar el estado del caso a "cd" (cerrado) en la base de datos
	if err := s.caseRepo.UpdateStatus(ctx, fu.CaseID, "cd"); err != nil {
		return fmt.Errorf("processCaseClosureSubmission: actualizar estado de caso: %w", err)
	}

	// 5. Cerrar los demás seguimientos pendientes/reprogramados del caso
	if err := s.followUpRepo.CloseCaseFollowUps(ctx, fu.ID); err != nil {
		return fmt.Errorf("processCaseClosureSubmission: cerrar seguimientos del caso: %w", err)
	}

	// 5.1 Crear Oficio (EntityLetter) y Tarea (CaseTask) si cumple condición crítica de cierre
	motivo := answerMap[qMotivoCierre]
	rutaVal := answerMap[qRutaAtencionPrevia]
	if motivo == "perdida_contacto" && rutaVal == "false" {
		letter := &models.EntityLetter{
			CaseID:     fu.CaseID,
			State:      models.EntityLetterStatePorProyectar, // "por_proyectar"
			Priority:   "normal",
			AgentID:    &actorID,
			RegisterBy: &actorID,
		}
		if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
			log.Printf("[WARN] No se pudo crear EntityLetter de cierre obligatorio: %v", err)
			return fmt.Errorf("crear entity_letter: %w", err)
		}

		task := &models.CaseTask{
			Category:       "Oficios",
			Type:           "Escribir oficio",
			Description:    "Debido a que no se registran acciones previas y se ha perdido el contacto, el protocolo exige la activación de la ruta de emergencia al cierre. Se generará el oficio de cierre obligatorio para proteger el estado de la usuaria.",
			AssignedUserID: actorID,
			Status:         models.CaseTaskStatusToDo, // "ToDo"
			CaseID:         fu.CaseID,
			FollowUpID:     &fu.ID,
			EntityLetterID: &letter.ID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := s.caseTaskRepo.Create(ctx, task); err != nil {
			log.Printf("[WARN] No se pudo crear CaseTask de cierre obligatorio: %v", err)
		} else {
			log.Printf("[CaseTask] Creado oficio de cierre obligatorio exitosamente para caso %s", fu.CaseID)
		}
	}

	// 6. Crear evento en el timeline
	if s.caseTimelineRepo != nil {
		motivo := answerMap[qMotivoCierre]
		motivoTraducido := motivo
		switch motivo {
		case "perdida_contacto":
			motivoTraducido = "Pérdida de contacto"
		case "solicitud_ciudadana":
			motivoTraducido = "Solicitud expresa de la ciudadana de finalizar el proceso"
		case "cumplimiento_plan":
			motivoTraducido = "Cumplimiento del plan de atención"
		case "no_corresponde":
			motivoTraducido = "No corresponde al ámbito, población o naturaleza de la atención"
		case "otro":
			motivoTraducido = "Otro"
		}

		causa := answerMap[qCausaCierre]
		accionesVal := answerMap[qAccionesCierre]
		accionesTraducido := "No"
		if accionesVal == "true" {
			accionesTraducido = "Sí"
		}

		rutaVal := answerMap[qRutaAtencionPrevia]
		rutaTraducida := "No"
		if rutaVal == "true" {
			rutaTraducida = "Sí"
		}

		description := fmt.Sprintf(
			"Cierre de caso registrado. Motivo: %s. Causa: %s. ¿Acciones institucionales realizadas?: %s. ¿Ruta de atención o algún oficio previamente activo?: %s",
			motivoTraducido,
			causa,
			accionesTraducido,
			rutaTraducida,
		)

		now := time.Now()
		actorName := actorID
		if s.agentLightRepo != nil && actorID != "" {
			agent, err := s.agentLightRepo.FindByICode(ctx, actorID)
			if err == nil && agent != nil {
				actorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}

		event := &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			FollowUpID:  fu.ID,
			Category:    models.TimelineCategoryGeneral,
			Type:        models.TimelineTypeCierreCaso,
			Icon:        models.TimelineIconCierre,
			Date:        now,
			Description: description,
			EventUserID: actorID,
			ActorName:   actorName,
			Color:       models.TimelineColorRed,
			CreatedAt:   now,
		}
		if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
			log.Printf("[processCaseClosureSubmission] advertencia: no se pudo crear evento timeline de cierre: %v", err)
		}
	}

	return nil
}

// processFollowUpSubmission marca el follow_up_v2 asociado como REALIZADO y crea las entidades
// derivadas (barreras, medidas de emergencia, derivaciones) a partir de las respuestas.
func (s *formService) processFollowUpSubmission(ctx context.Context, submissionID, actorID string) error {
	// 1. Buscar el follow_up asociado al submission
	fu, err := s.followUpRepo.FindByFormSubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processFollowUpSubmission: buscar followUp: %w", err)
	}
	// Si ya fue procesado → registrar edición en el timeline y salir
	if fu.Status == models.FollowUpStatusRealizado {
		if s.caseTimelineRepo == nil {
			log.Printf("⚠️  [processFollowUpSubmission] caseTimelineRepo es nil — el evento 'Seguimiento Editado' NO fue creado para followUp=%s", fu.ID)
			return nil
		}
		now := time.Now()

		// Resolver nombre del agente
		actorName := actorID
		if s.agentLightRepo != nil && actorID != "" {
			agent, err := s.agentLightRepo.FindByICode(ctx, actorID)
			if err == nil && agent != nil {
				actorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}

		editEvent := &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			FollowUpID:  fu.ID,
			Category:    models.TimelineCategorySeguimientos,
			Type:        models.TimelineTypeSeguimientoEditado,
			Icon:        models.TimelineIconPospuesto,
			Date:        now,
			Description: "Seguimiento editado después de su ejecución",
			EventUserID: actorID,
			ActorName:   actorName,
			Color:       models.TimelineColorTeal,
			CreatedAt:   now,
		}
		if err := s.caseTimelineRepo.Create(ctx, editEvent); err != nil {
			log.Printf("[processFollowUpSubmission] advertencia: no se pudo crear evento 'Seguimiento Editado': %v", err)
		}
		return nil
	}

	// 2. Leer las respuestas directas del submission
	answers, err := s.answerRepo.FindDirectBySubmissionID(ctx, submissionID)
	if err != nil {
		return fmt.Errorf("processFollowUpSubmission: leer respuestas: %w", err)
	}

	// Indexar respuestas por questionId para acceso rápido
	answerMap := make(map[string]string, len(answers))
	for _, a := range answers {
		answerMap[a.QuestionID] = a.Value
	}

	// ── IDs de las preguntas del formulario de seguimiento ────────────────────
	const (
		// Repeater group de Sección 3 — Seguimiento a Barreras
		rgSeguimientoBarreras      = "b536f16c-67b3-4370-810d-7cc8c9d5463e"
		qSBPersiste                = "4325a514-332a-42a9-9a89-a1d249b87c71" // Q2 boolean
		qSBRespuestaInstitucional  = "a234a906-7b02-4a43-b18d-b802d430d272" // Q3 single
		qSBGestion                 = "c785a994-3a38-4337-bd96-67329144affa" // Q4 multiple
		qSBActuaciones             = "c808b590-128c-43cf-af94-c191a4bc31ac" // Q5 text
		qSBCierra                  = "7442394d-256a-45f4-a63a-ca1164e865c2" // Q6 boolean
		qSBMotivoCierre            = "d2d3cee8-b5d3-422f-8ae0-e5a40effc983" // Q7 single
		// Repeater group de barreras
		rgBarreras = "5fd3ecdc-2e5f-4b31-97ef-8a994580586a"
		// Preguntas dentro del repeater de barreras — bloque sector
		qBarreraSector          = "f19378b6-55c5-4fdf-b765-7ebcc3978741" // Q1  Sector (dropdown)
		qBarreraSalud           = "5fc1f2af-cc30-41f0-aa31-4731e5cb674c" // Q2  Barreras Salud (multiple)
		qInstitucionSalud       = "e88fb2bf-7196-4bea-91e6-fe78ccdedac1" // Q3  Institución Salud (multiple)
		qOtraBarreraSalud       = "bebf6e6c-0200-4b53-886c-c01f591eaa62" // Q4  Otra barrera Salud (text)
		qBarreraJusticia        = "2bec977e-97c7-42d7-a00a-b536af8038eb" // Q5  Barreras Justicia (multiple)
		qInstitucionJusticia    = "4b4997fa-67a0-4479-852d-df9e7bfb2b3e" // Q6  Institución Justicia (multiple)
		qOtraBarreraJusticia    = "96f0c507-64c1-446b-b07d-83d5ad1d7172" // Q7  Otra barrera Justicia (text)
		qBarreraProteccion      = "66c9fc1e-9b5e-4ad4-999f-7aeb483d84dc" // Q8  Barreras Protección (multiple)
		qInstitucionProteccion  = "78474c82-61b9-4a0c-beca-439274813c03" // Q9  Institución Protección (multiple)
		qOtraBarreraProteccion  = "68f6bf06-a6a8-443f-9a1e-001148d4eb45" // Q10 Otra barrera Protección (text)
		qBarreraOtraInstitucion = "a1573bc3-28f6-485e-8d8b-b3567db42ae3" // Q11 Nombre institución (text)
		// Preguntas comunes del repeater de barreras
		qBarreraDepartamento           = "31c7f8ba-880e-4c9a-89f0-1a6a1e43b7b9" // Q12 Departamento (dropdown)
		qBarreraCiudad                 = "c6f2d54a-3c61-4ae7-b2bf-50a884031fa4" // Q13 Ciudad (dropdown)
		qBarreraMunicipio              = "8225d03f-8de9-4ff0-9345-67b5e71bf02d" // Q14 Municipio (dropdown)
		qBarreraEstructuralInstitucional = "64754fb5-04a8-43e4-ac86-d60f0a84010b" // Q15 (multiple)
		qBarreraEstructuralEconomico   = "94dfc417-b59f-4afe-b63b-c6839a8dfecc" // Q16 (multiple)
		qBarreraEstructuralTerritorial = "cf162686-0327-4b28-9639-e4d947272483" // Q17 (multiple)
		qBarreraEstructuralDiferencial = "81275638-2837-4ecc-8686-675dfeab4f32" // Q18 (multiple)
		qBarreraFecha                  = "4592d85f-8c11-4785-ad32-06aa810cc491" // Q19 date
		qBarreraFuncionario            = "33c3961e-4252-4b02-b491-615b11d6bf56" // Q20 text
		qBarreraDescripcion            = "d2be610f-44eb-4526-b4ec-86eaaaba08c8" // Q21 text
		qBarreraGestion                = "572ad72a-8174-4ff3-9c56-5c8c65ac63ac" // Q22 (multiple)
		// Preguntas directas del form
		qEquipos            = "e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc" // Derivaciones a equipos (multi-select)
		qMedidasEmergencia  = "1a36260c-33a4-4ebd-bffb-e387d7964b96" // Medidas de emergencia (multi-select)
		qCriteriosPsico     = "71c42c4a-f640-47ad-b2c1-5d4c18480449" // Criterios de remisión — Atención Psicosocial (multiple)
		qCriteriosHombres        = "f7edf4fc-d1cd-4591-a358-31566c806365" // Criterios de remisión — Atención Hombres (multiple)
		qCriteriosEstabilizacion = "28accaa6-99dc-4ec4-967b-f26045ad707c" // Criterios de remisión — Estabilización (multiple)
		qServiciosDiscapacidad   = "47b122b1-0151-4b58-a007-d4afb722c1b9" // Servicios del Equipo de Discapacidad (multiple)
		// Sección 1 — Valoración del Riesgo: nuevos hechos de violencia
		qNuevosHechosViolencia = "69ccecbe-8fd5-44a4-901a-17084ea7134d" // Q2 boolean — ¿Se registraron nuevos hechos?
		qDescripcionHechos     = "f8453544-2a8c-461d-9986-302ee4719492" // Q3 text   — Descripción (visible si Q2=true)
		qFechaHechos           = "7373aeab-b6d8-44e4-b1ce-4d40196f17bc" // Q4 date   — Fecha en que ocurrieron
		// Sección 5 — Cierre del caso
		qCierraCaso         = "08950a38-3db3-4dc7-852c-3b06b4b1ed72" // boolean — ¿Realiza cierre del caso?
		qCierreMotivo       = "95fb963e-99de-4a1d-a170-8e30845d1f7d" // single  — Motivo del cierre
		qCierreOtroMotivo   = "e259ff16-d049-41d1-b92a-b0f09ee6fa91" // text    — Otro motivo ¿cuál?
		qCierreDescripcion  = "50ab05f3-95d9-42e4-bfed-6abfcd0a8050" // text    — Describa la causa
		qCierreAccionesInst = "0e7b61c8-9401-4429-81ed-61c8fc97808e" // boolean — ¿Realizó acciones institucionales?
	)

	// Mapas sector → pregunta de barreras, instituciones y "otra barrera"
	sectorBarrierQ := map[string]string{
		"salud":      qBarreraSalud,
		"justicia":   qBarreraJusticia,
		"proteccion": qBarreraProteccion,
	}
	sectorInstitutionQ := map[string]string{
		"salud":      qInstitucionSalud,
		"justicia":   qInstitucionJusticia,
		"proteccion": qInstitucionProteccion,
	}
	sectorOtherBarrierQ := map[string]string{
		"salud":      qOtraBarreraSalud,
		"justicia":   qOtraBarreraJusticia,
		"proteccion": qOtraBarreraProteccion,
	}

	// ── Helper: parsea valor comma-separated y filtra vacíos ─────────────────
	splitValues := func(v string) []string {
		var out []string
		for _, part := range strings.Split(v, ",") {
			if t := strings.TrimSpace(part); t != "" {
				out = append(out, t)
			}
		}
		return out
	}

	// Contadores para la descripción del evento en el timeline
	barrierCount  := 0
	remisionCount := 0

	// 3. Crear un BarrierV2 por cada entrada del repeater de barreras
	barrierEntries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, []string{rgBarreras})
	if err != nil {
		return fmt.Errorf("processFollowUpSubmission: leer entradas de barreras: %w", err)
	}
	log.Printf("[processFollowUp] barreras: %d entradas encontradas", len(barrierEntries))
	for _, entry := range barrierEntries {
		entryAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
		if err != nil {
			return fmt.Errorf("processFollowUpSubmission: leer respuestas entrada [%s]: %w", entry.ID, err)
		}
		entryMap := make(map[string]string, len(entryAnswers))
		for _, a := range entryAnswers {
			entryMap[a.QuestionID] = a.Value
		}

		sector := strings.TrimSpace(entryMap[qBarreraSector])

		b := &models.BarrierV2{
			CaseID:      fu.CaseID,
			FollowUpID:  fu.ID,
			CreatedByID: actorID,
			Status:      models.BarrierV2StatusOpen,

			Sector:               sector,
			SpecificBarriers:     entryMap[sectorBarrierQ[sector]],
			SpecificInstitutions: entryMap[sectorInstitutionQ[sector]],
			OtherBarrierDesc:     entryMap[sectorOtherBarrierQ[sector]],
			InstitutionName:      entryMap[qBarreraOtraInstitucion],

			DepartmentID: entryMap[qBarreraDepartamento],
			CityID:       entryMap[qBarreraCiudad],
			TownID:       entryMap[qBarreraMunicipio],

			StructuralInstitutional: entryMap[qBarreraEstructuralInstitucional],
			StructuralEconomic:      entryMap[qBarreraEstructuralEconomico],
			StructuralTerritorial:   entryMap[qBarreraEstructuralTerritorial],
			StructuralDifferential:  entryMap[qBarreraEstructuralDiferencial],

			BarrierDate:        entryMap[qBarreraFecha],
			OfficialDependency: entryMap[qBarreraFuncionario],
			Description:        entryMap[qBarreraDescripcion],
			ManagementActions:  entryMap[qBarreraGestion],
		}
		log.Printf("[processFollowUp] creando barrera sector=%s entry=%s", sector, entry.ID)
		if err := s.barrierV2Repo.Create(ctx, b); err != nil {
			return fmt.Errorf("processFollowUpSubmission: crear barrera entry [%s]: %w", entry.ID, err)
		}
		barrierCount++

		// 3.3 Crear tareas (y oficios si aplica) por cada opción seleccionada en gestión
		gestionRaw := strings.TrimSpace(entryMap[qBarreraGestion])
		if gestionRaw != "" {
			// Opciones que no generan ningún registro (ni tarea ni oficio)
			gestionSinTarea := map[string]bool{
				"orientacion_llamada": true,
			}
			// Opciones de gestión que requieren un oficio (entity_letter)
			gestionConOficio := map[string]bool{
				"activacion_ruta_interinstitucional": true,
				"articulacion_institucional":         true,
				"escalamiento_organismo_control":     true,
			}
			// Label legible por opción de gestión
			gestionLabels := map[string]string{
				"gestion_llamada":                    "Gestión administrativa - Llamada",
				"activacion_ruta_interinstitucional": "Activación de ruta interinstitucional",
				"articulacion_institucional":         "Articulación institucional",
				"escalamiento_organismo_control":     "Escalamiento a organismo de control",
				"alerta_barreras":                    "Alerta por barreras",
			}

			for _, gVal := range strings.Split(gestionRaw, ",") {
				gVal = strings.TrimSpace(gVal)
				if gVal == "" || gestionSinTarea[gVal] {
					continue
				}

				label, ok := gestionLabels[gVal]
				if !ok {
					label = gVal
				}

				institution := b.InstitutionName
				if institution == "" {
					institution = b.SpecificInstitutions
				}
				label = label + " (Barreras)"
				if sector != "" || institution != "" {
					label = fmt.Sprintf("%s | Sector: %s | Institución: %s", label, sector, institution)
				}

				barrierID := b.ID
				followUpID := fu.ID
				var entityLetterID *string

				if gestionConOficio[gVal] {
					// Crear entity_letter en estado por_proyectar
					letter := &models.EntityLetter{
						BarrierID: barrierID,
						CaseID:    fu.CaseID,
						State:     models.EntityLetterStatePorProyectar,
						AgentID:   &actorID,
					}
					if s.entityLetterRepo != nil {
						if err := s.entityLetterRepo.Create(ctx, letter); err != nil {
							log.Printf("[processFollowUp] WARN: no se pudo crear entity_letter para barrera %s gestion=%s: %v", barrierID, gVal, err)
						} else {
							entityLetterID = &letter.ID
						}
					}
				}

				taskType := gVal
				if gestionConOficio[gVal] {
					taskType = "proyectar_oficio"
				} else if gVal == "alerta_barreras" {
					taskType = "comite_caso"
				}

				task := &models.CaseTask{
					Category:       "Barreras",
					Type:           taskType,
					Description:    label,
					AssignedUserID: actorID,
					Status:         models.CaseTaskStatusToDo,
					CaseID:         fu.CaseID,
					FollowUpID:     &followUpID,
					BarrierID:      &barrierID,
					EntityLetterID: entityLetterID,
				}
				if s.caseTaskRepo != nil {
					if err := s.caseTaskRepo.Create(ctx, task); err != nil {
						log.Printf("[processFollowUp] WARN: no se pudo crear case_task para barrera %s gestion=%s: %v", barrierID, gVal, err)
					}
				}
			}
		}
	}

	// 3b. Procesar Seguimiento a Barreras (Sección 3)
	// Cada entry del repeater corresponde por posición a un ID en fu.ActiveBarrierIDs.
	sbEntries, err := s.repeaterEntryRepo.FindBySubmissionIDAndGroupIDs(ctx, submissionID, []string{rgSeguimientoBarreras})
	if err != nil {
		return fmt.Errorf("processFollowUpSubmission: leer entradas seguimiento barreras: %w", err)
	}
	log.Printf("[processFollowUp] seguimiento barreras: %d entries encontradas", len(sbEntries))
	if len(sbEntries) > 0 {
		// Obtener IDs de barreras activas del follow-up para relacionar por posición
		var activeIDs []string
		if fu.ActiveBarrierIDs != nil && *fu.ActiveBarrierIDs != "" {
			for _, id := range strings.Split(*fu.ActiveBarrierIDs, ",") {
				if t := strings.TrimSpace(id); t != "" {
					activeIDs = append(activeIDs, t)
				}
			}
		}
		log.Printf("[processFollowUp] active_barrier_ids disponibles: %d → %v", len(activeIDs), activeIDs)
		for idx, entry := range sbEntries {
			sbAnswers, err := s.answerRepo.FindByRepeaterEntryID(ctx, entry.ID)
			if err != nil {
				return fmt.Errorf("processFollowUpSubmission: leer respuestas seguimiento barrera [%s]: %w", entry.ID, err)
			}
			sbMap := make(map[string]string, len(sbAnswers))
			for _, a := range sbAnswers {
				sbMap[a.QuestionID] = a.Value
			}

			// Relacionar entry con barrera por posición
			barrierID := ""
			if idx < len(activeIDs) {
				barrierID = activeIDs[idx]
			}

			log.Printf("[processFollowUp] seguimiento barrera [%d] → barrierID=%s | persiste=%s | cierra=%s | respInstitucional=%s",
				idx, barrierID, sbMap[qSBPersiste], sbMap[qSBCierra], sbMap[qSBRespuestaInstitucional])

			// Si la barrera fue cerrada → actualizar status
			if sbMap[qSBCierra] == "true" && barrierID != "" {
				if err := s.barrierV2Repo.UpdateStatus(ctx, barrierID, models.BarrierV2StatusManaged); err != nil {
					log.Printf("[processFollowUp] ❌ no se pudo cerrar barrera %s: %v", barrierID, err)
				} else {
					log.Printf("[processFollowUp] ✅ barrera cerrada → id=%s", barrierID)
				}
			}

			// Crear evento en timeline por este seguimiento de barrera
			if s.caseTimelineRepo != nil {
				resumen := buildBarrierFollowUpSummary(sbMap[qSBPersiste], sbMap[qSBRespuestaInstitucional], sbMap[qSBActuaciones])
				log.Printf("[processFollowUp] creando evento timeline barrera [%d] → resumen=%q", idx, resumen)
				tlEvent := &models.CaseTimelineEvent{
					CaseID:      fu.CaseID,
					FollowUpID:  fu.ID,
					Category:    "Barreras",
					Type:        "Seguimiento a Barrera",
					Icon:        "shield-halved",
					Color:       "#6366f1",
					Description: resumen,
					EventUserID: actorID,
					Date:        time.Now(),
				}
				if err := s.caseTimelineRepo.Create(ctx, tlEvent); err != nil {
					log.Printf("[processFollowUp] ❌ evento timeline barrera %s: %v", barrierID, err)
				} else {
					log.Printf("[processFollowUp] ✅ evento timeline creado → barrera [%d] id=%s", idx, barrierID)
				}
			}

			// Persistir respuestas estructuradas del seguimiento a barrera
			if s.barrierFollowUpRepo != nil && barrierID != "" {
				bfu := &models.BarrierFollowUp{
					BarrierID:             barrierID,
					FollowUpID:            fu.ID,
					CreatedByID:           actorID,
					Persists:              sbMap[qSBPersiste] == "true",
					InstitutionalResponse: sbMap[qSBRespuestaInstitucional],
					ManagementActions:     sbMap[qSBGestion],
					Actions:               sbMap[qSBActuaciones],
					ClosesBarrier:         sbMap[qSBCierra] == "true",
					ClosureReason:         sbMap[qSBMotivoCierre],
				}
				if err := s.barrierFollowUpRepo.Create(ctx, bfu); err != nil {
					log.Printf("[processFollowUp] ❌ barrier_follow_up barrera %s: %v", barrierID, err)
				} else {
					log.Printf("[processFollowUp] ✅ barrier_follow_up creado → barrera %s", barrierID)
				}
			}
		}
	}

	// 4. Crear registros de derivación por equipo
	if equiposVal, ok := answerMap[qEquipos]; ok && equiposVal != "" {
		log.Printf("[processFollowUp] derivaciones a equipos: %s", equiposVal)
		for _, equipo := range splitValues(equiposVal) {
			switch equipo {
			case "atencion_psico":
				// Regla de exclusión: si también se seleccionó medidas_emergencia, solo se crea
				// la derivación de emergencia. La remisión psicosocial se omite.
				equiposSeleccionados := splitValues(equiposVal)
				tieneMedidasEmergencia := false
				for _, e := range equiposSeleccionados {
					if e == "medidas_emergencia" {
						tieneMedidasEmergencia = true
						break
					}
				}
				if tieneMedidasEmergencia {
					log.Printf("[processFollowUp] atencion_psico omitida — incompatible con medidas_emergencia seleccionada al mismo tiempo")
					break
				}
				// Validar criterios de remisión psicosocial antes de crear la derivación.
				// Regla: debe estar marcado "criterio_obligatorio" Y sumar mínimo 3 puntos.
				criteriosVal := answerMap[qCriteriosPsico]
				criteriosSelected := splitValues(criteriosVal)
				tieneCriterioObligatorio := false
				for _, c := range criteriosSelected {
					if c == "criterio_obligatorio" {
						tieneCriterioObligatorio = true
						break
					}
				}
				puntajesPsico := map[string]int{
					"conducta_suicida":         3,
					"interseccionalidad":       2,
					"sin_ruta":                 1,
					"condiciones_territoriales": 1,
					"sin_acceso_psico":         1,
					"naturalizacion_vbg":       1,
				}
				totalPuntos := 0
				for _, c := range criteriosSelected {
					totalPuntos += puntajesPsico[c]
				}
				if !tieneCriterioObligatorio || totalPuntos < 3 {
					log.Printf("[processFollowUp] derivacion psicosocial NO cumple criterios (obligatorio=%v, puntos=%d) — omitida", tieneCriterioObligatorio, totalPuntos)
					break
				}
				log.Printf("[processFollowUp] creando derivacion -> atencion_psico (obligatorio=%v, puntos=%d)", tieneCriterioObligatorio, totalPuntos)
				ps := &models.PsychosocialSupport{
					CaseID:     fu.CaseID,
					FollowUpID: fu.ID,
					Type:       "derivacion",
					Status:     "ACTIVE",
				}
				if err := s.psychosocialSupportRepo.Create(ctx, ps); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear derivacion psicosocial: %w", err)
				}
				log.Printf("[processFollowUp] ✅ derivacion creada -> atencion_psico (id=%s)", ps.ID)
				remisionCount++
			case "atencion_hombres":
				// Validar criterio de remisión al equipo de hombres antes de crear la derivación.
				// Regla: debe estar marcado "criterio_hombres".
				criteriosHombresVal := answerMap[qCriteriosHombres]
				criteriosHombresSelected := splitValues(criteriosHombresVal)
				tieneCriterioHombres := false
				for _, c := range criteriosHombresSelected {
					if c == "criterio_hombres" {
						tieneCriterioHombres = true
						break
					}
				}
				if !tieneCriterioHombres {
					log.Printf("[processFollowUp] derivacion atencion_hombres NO cumple criterio — omitida")
					break
				}
				log.Printf("[processFollowUp] creando derivacion -> atencion_hombres")
				mtr := &models.MenTeamRemision{
					CaseID:     fu.CaseID,
					FollowUpID: fu.ID,
					Type:       "derivacion",
					Status:     "ACTIVE",
				}
				if err := s.menTeamRemisionRepo.Create(ctx, mtr); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear derivacion atencion_hombres: %w", err)
				}
				log.Printf("[processFollowUp] ✅ derivacion creada -> atencion_hombres (id=%s)", mtr.ID)
				remisionCount++
			case "discapacidad":
				// Crear un registro por cada servicio seleccionado del Equipo de Discapacidad.
				serviciosVal := answerMap[qServiciosDiscapacidad]
				if serviciosVal == "" {
					log.Printf("[processFollowUp] derivacion discapacidad sin servicios seleccionados — omitida")
					break
				}
				for _, svc := range splitValues(serviciosVal) {
					log.Printf("[processFollowUp] creando derivacion -> discapacidad servicio=%s", svc)
					dr := &models.DiscapacidadRemision{
						CaseID:     fu.CaseID,
						FollowUpID: fu.ID,
						Service:    svc,
						Status:     "ACTIVE",
					}
					if err := s.discapacidadRemisionRepo.Create(ctx, dr); err != nil {
						return fmt.Errorf("processFollowUpSubmission: crear derivacion discapacidad [%s]: %w", svc, err)
					}
					log.Printf("[processFollowUp] ✅ derivacion creada -> discapacidad servicio=%s (id=%s)", svc, dr.ID)
					remisionCount++
				}
			case "estabilizacion":
				// Validar que se haya seleccionado al menos 1 criterio de estabilización.
				criteriosEstVal := answerMap[qCriteriosEstabilizacion]
				if criteriosEstVal == "" || len(splitValues(criteriosEstVal)) == 0 {
					log.Printf("[processFollowUp] derivacion estabilizacion sin criterios seleccionados — omitida")
					break
				}
				log.Printf("[processFollowUp] creando derivacion -> estabilizacion (criterios: %s)", criteriosEstVal)
				ec := &models.EconomicStabilization{
					CaseID:     fu.CaseID,
					FollowUpID: fu.ID,
					Type:       "derivacion",
					Status:     "ACTIVE",
				}
				if err := s.economicStabilizationRepo.Create(ctx, ec); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear derivacion economica: %w", err)
				}
				log.Printf("[processFollowUp] ✅ derivacion creada -> estabilizacion (id=%s)", ec.ID)
				remisionCount++
			}
		}
	}

	// 5. Crear una EmergencyMeasure por cada medida de emergencia seleccionada
	if medidasVal, ok := answerMap[qMedidasEmergencia]; ok && medidasVal != "" {
		log.Printf("[processFollowUp] medidas de emergencia: %s", medidasVal)
		for _, medida := range splitValues(medidasVal) {
			log.Printf("[processFollowUp] creando medida emergencia tipo=%s", medida)
			em := &models.EmergencyMeasure{
				CaseID:     fu.CaseID,
				FollowUpID: fu.ID,
				Type:       medida,
				Status:     "ACTIVE",
			}
			if err := s.emergencyMeasureRepo.Create(ctx, em); err != nil {
				return fmt.Errorf("processFollowUpSubmission: crear medida emergencia [%s]: %w", medida, err)
			}
			log.Printf("[processFollowUp] ✅ medida emergencia creada tipo=%s (id=%s)", medida, em.ID)
			remisionCount++
		}
	}

	// 6. Marcar el follow_up como REALIZADO
	if err := s.followUpRepo.UpdateStatus(ctx, fu.ID, models.FollowUpStatusRealizado); err != nil {
		return fmt.Errorf("processFollowUpSubmission: actualizar estado followUp: %w", err)
	}

	// 7. Crear evento en el timeline del caso
	if s.caseTimelineRepo == nil {
		log.Printf("⚠️  [processFollowUpSubmission] caseTimelineRepo es nil — el evento de timeline NO fue creado para followUp=%s. Inyectar CaseTimelineEventRepo en FormServiceDeps (main.go)", fu.ID)
	}
	if s.caseTimelineRepo != nil {
		description := fmt.Sprintf(
			"Seguimiento ejecutado. Se identificaron %d %s y se realizaron %d %s a Equipos Salvia",
			barrierCount,
			pluralize(barrierCount, "barrera", "barreras"),
			remisionCount,
			pluralize(remisionCount, "remisión", "remisiones"),
		)
		now := time.Now()

		// Resolver nombre del agente que ejecutó el seguimiento
		actorName := actorID
		if s.agentLightRepo != nil && actorID != "" {
			agent, err := s.agentLightRepo.FindByICode(ctx, actorID)
			if err == nil && agent != nil {
				actorName = strings.TrimSpace(agent.Names + " " + agent.LastNames)
			}
		}

		event := &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			FollowUpID:  fu.ID,
			Category:    models.TimelineCategorySeguimientos,
			Type:        models.TimelineTypeSeguimientoEjecutado,
			Icon:        models.TimelineIconSeguimiento,
			Date:        now,
			Description: description,
			EventUserID: actorID,
			ActorName:   actorName,
			Color:       models.TimelineColorGreen,
			CreatedAt:   now,
		}
		if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
			log.Printf("[processFollowUp] advertencia: no se pudo crear evento timeline: %v", err)
		}
	}

	// 7b. Crear evento de nuevos hechos de violencia (si se reportaron en Sección 1)
	if answerMap[qNuevosHechosViolencia] == "true" && s.caseTimelineRepo != nil {
		descripcion := strings.TrimSpace(answerMap[qDescripcionHechos])
		fechaStr := strings.TrimSpace(answerMap[qFechaHechos])

		fechaHechos := time.Now()
		if fechaStr != "" {
			if t, err := time.Parse("2006-01-02", fechaStr); err == nil {
				fechaHechos = t
			}
		}

		hechoEvent := &models.CaseTimelineEvent{
			CaseID:      fu.CaseID,
			FollowUpID:  fu.ID,
			Category:    models.TimelineCategoryGeneral,
			Type:        models.TimelineTypeHechosCaso,
			Icon:        models.TimelineIconHechosCaso,
			Color:       "#f87171",
			Description: descripcion,
			EventUserID: actorID,
			Date:        fechaHechos,
			CreatedAt:   time.Now(),
		}
		if err := s.caseTimelineRepo.Create(ctx, hechoEvent); err != nil {
			log.Printf("[processFollowUp] advertencia: no se pudo crear evento nuevos hechos: %v", err)
		} else {
			log.Printf("[processFollowUp] ✅ evento 'Nuevos hechos de violencia' creado → fecha=%s", fechaStr)
		}
	}

	// 8. Cierre del caso: si el profesional marcó ciorre en la Sección 5, delegar al servicio modular
	if answerMap[qCierraCaso] == "true" && s.casoCierreService != nil {
		input := CerrarCasoInput{
			CaseICode:               fu.CaseID,
			Motivo:                  answerMap[qCierreMotivo],
			OtroMotivo:              answerMap[qCierreOtroMotivo],
			Descripcion:             answerMap[qCierreDescripcion],
			AccionesInstitucionales: answerMap[qCierreAccionesInst] == "true",
			ActorID:                 actorID,
		}
		if err := s.casoCierreService.CerrarCaso(ctx, input); err != nil {
			// No retornamos error: el cierre fallido no debe revertir el seguimiento ya completado
			log.Printf("[processFollowUpSubmission] advertencia: no se pudo cerrar el caso %s: %v", fu.CaseID, err)
		}
	} else if answerMap[qCierraCaso] == "true" && s.casoCierreService == nil {
		log.Printf("⚠️  [processFollowUpSubmission] casoCierreService es nil — cierre del caso %s no ejecutado. Inyectar CasoCierreService en FormServiceDeps (main.go)", fu.CaseID)
	}

	// 9. Reasignación de caso: evalúa nivel de riesgo y reasigna calendario o agente si corresponde
	s.reasignarCaso(ctx, fu, answerMap, actorID)

	return nil
}

// reasignarCaso evalúa las respuestas del formulario de seguimiento y, si corresponde,
// actualiza el nivel de riesgo del caso y reasigna el calendario o el agente.
//
// Niveles de destino:
//   - Extremo automático (caso bajo + factor extremo seleccionado) → 4
//   - Confirmación a alto  (caso bajo + ≥4 factores de riesgo)     → 3
//   - Confirmación a bajo  (caso alto + ≥3 factores protectores)   → 2
func (s *formService) reasignarCaso(ctx context.Context, fu *models.FollowUpV2, answerMap map[string]string, actorID string) {
	const (
		qProtectores = "a0fdcf67-b05b-4d92-9c19-14753a32bbe3" // Factores protectores (multiple)
		qRiesgos     = "ec5bb242-6f86-4c64-8b9f-afabd5a51878" // Factores de riesgo (multiple)
		qExtremo     = "65f2d582-a39d-4c93-9519-1a20efbebb03" // Factores de riesgo extremo (multiple)
		qConfirmHigh = "df7a0293-e2ca-49e4-88a8-66a70519a9e5" // ¿Confirmar reasignación a riesgo alto? (boolean)
		qConfirmLow  = "7ec8d66d-7015-470a-a596-edebe574b5c2" // ¿Confirmar reasignación a riesgo bajo? (boolean)
	)

	splitCSV := func(v string) []string {
		if v == "" {
			return nil
		}
		var out []string
		for _, p := range strings.Split(v, ",") {
			if t := strings.TrimSpace(p); t != "" {
				out = append(out, t)
			}
		}
		return out
	}

	sinNinguno := func(vals []string) []string {
		out := vals[:0]
		for _, v := range vals {
			if v != "ninguno" {
				out = append(out, v)
			}
		}
		return out
	}
	protectores := sinNinguno(splitCSV(answerMap[qProtectores]))
	riesgos     := sinNinguno(splitCSV(answerMap[qRiesgos]))
	extremos    := sinNinguno(splitCSV(answerMap[qExtremo]))
	confirmHigh := answerMap[qConfirmHigh] == "true"
	confirmLow  := answerMap[qConfirmLow]  == "true"

	// Leer nivel de riesgo actual del caso desde victim_case_form2
	currentLevel, err := s.caseRepo.FindRiskLevelByICode(ctx, fu.CaseID)
	if err != nil {
		log.Printf("[reasignarCaso] advertencia: no se pudo leer risk_level del caso %s: %v", fu.CaseID, err)
		return
	}
	log.Printf("[reasignarCaso] caseID=%s currentLevel=%d extremos=%d riesgos=%d protectores=%d confirmHigh=%v confirmLow=%v",
		fu.CaseID, currentLevel, len(extremos), len(riesgos), len(protectores), confirmHigh, confirmLow)

	// Determinar nuevo nivel
	newLevel := 0
	var razon string
	if currentLevel >= 1 && currentLevel <= 2 {
		if len(extremos) > 0 {
			newLevel = 4
			razon = "caso bajo + factor extremo → automático nivel 4"
		} else if confirmHigh && len(riesgos) >= 4 {
			newLevel = 3
			razon = fmt.Sprintf("caso bajo + %d riesgos + confirmHigh → nivel 3", len(riesgos))
		} else {
			razon = fmt.Sprintf("caso bajo — sin condición: extremos=%d riesgos=%d confirmHigh=%v", len(extremos), len(riesgos), confirmHigh)
		}
	} else if currentLevel >= 3 {
		if confirmLow && len(extremos) == 0 && len(protectores) >= 3 {
			newLevel = 2
			razon = fmt.Sprintf("caso alto + %d protectores + sin extremo + confirmLow → nivel 2", len(protectores))
		} else {
			razon = fmt.Sprintf("caso alto — sin condición: extremos=%d protectores=%d confirmLow=%v", len(extremos), len(protectores), confirmLow)
		}
	} else {
		razon = fmt.Sprintf("nivel fuera de rango: %d", currentLevel)
	}

	log.Printf("[reasignarCaso] decisión: newLevel=%d | %s", newLevel, razon)

	if newLevel == 0 {
		return
	}

	log.Printf("[reasignarCaso] iniciando reasignación caseID=%s %d→%d", fu.CaseID, currentLevel, newLevel)

	// Actualizar nivel de riesgo en victim_case_form2
	if err := s.caseRepo.UpdateRiskLevelByICode(ctx, fu.CaseID, newLevel); err != nil {
		log.Printf("[reasignarCaso] advertencia: no se pudo actualizar risk_level caseID=%s: %v", fu.CaseID, err)
		return
	}

	// Reasignar calendario o agente según el nuevo nivel
	if s.followUpV2Svc == nil {
		log.Printf("⚠️  [reasignarCaso] followUpV2Svc es nil — reasignación de calendario no ejecutada. Inyectar FollowUpV2Svc en FormServiceDeps (main.go)")
		return
	}
	if err := s.followUpV2Svc.ReasignarCalendario(ctx, fu.CaseID, newLevel); err != nil {
		log.Printf("[reasignarCaso] advertencia: ReasignarCalendario falló caseID=%s: %v", fu.CaseID, err)
		return
	}

	// Actualizar equipo y agente en victim_case
	team := "Riesgo bajo"
	if newLevel >= 3 {
		team = "Riesgo alto"
	}
	// Leer el agente del primer follow-up PENDIENTE recién asignado
	newAgentID := ""
	if pendientes, err := s.followUpRepo.FindPendingByCaseID(ctx, fu.CaseID); err == nil && len(pendientes) > 0 {
		if pendientes[0].AgentID != nil {
			newAgentID = *pendientes[0].AgentID
		}
	}
	if err := s.caseRepo.UpdateTeamAndAgent(ctx, fu.CaseID, team, newAgentID); err != nil {
		log.Printf("[reasignarCaso] advertencia: no se pudo actualizar team/agent en victim_case: %v", err)
	} else {
		log.Printf("[reasignarCaso] ✅ victim_case actualizado → team=%q agentID=%s", team, newAgentID)
	}

	// Evento de timeline "Reasignación de Caso"
	now := time.Now()
	event := &models.CaseTimelineEvent{
		CaseID:      fu.CaseID,
		FollowUpID:  fu.ID,
		Category:    models.TimelineCategoryGeneral,
		Type:        "Reasignación de Caso",
		Icon:        "arrows-rotate",
		Color:       "#f59e0b",
		Date:        now,
		EventUserID: actorID,
		CreatedAt:   now,
	}
	if err := s.caseTimelineRepo.Create(ctx, event); err != nil {
		log.Printf("[reasignarCaso] advertencia: no se pudo crear evento timeline: %v", err)
	}

	log.Printf("[reasignarCaso] ✅ reasignación completada caseID=%s %d→%d", fu.CaseID, currentLevel, newLevel)
}

// buildValidOptionValues carga las preguntas de una sección y retorna un mapa
// questionId → set de valores válidos, solo para tipos single/dropdown/multiple.
// Las preguntas de otros tipos no aparecen en el mapa (no se filtran).
func (s *formService) buildValidOptionValues(ctx context.Context, sectionID string) (map[string]map[string]bool, error) {
	questions, err := s.questionRepo.FindBySectionID(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	optionTypes := map[string]bool{"single": true, "dropdown": true, "multiple": true}
	var questionIDs []string
	for _, q := range questions {
		if optionTypes[q.QuestionTypeID] {
			questionIDs = append(questionIDs, q.ID)
		}
	}
	if len(questionIDs) == 0 {
		return map[string]map[string]bool{}, nil
	}

	options, err := s.optionRepo.FindByQuestionIDs(ctx, questionIDs)
	if err != nil {
		return nil, err
	}

	result := map[string]map[string]bool{}
	for _, o := range options {
		if result[o.QuestionID] == nil {
			result[o.QuestionID] = map[string]bool{}
		}
		result[o.QuestionID][o.Value] = true
	}
	return result, nil
}

// filterAnswerValue filtra el value de una respuesta según los valores válidos.
// Para multiple (CSV) elimina los valores inválidos. Para single/dropdown vacía si es inválido.
// Si la pregunta no está en el mapa (no tiene opciones), devuelve el valor sin tocar.
func filterAnswerValue(questionID, value string, validOptions map[string]map[string]bool) string {
	valid, hasOptions := validOptions[questionID]
	if !hasOptions {
		return value
	}
	// Detectar si es multiple por si tiene comas
	parts := strings.Split(value, ",")
	if len(parts) <= 1 {
		// single / dropdown
		if valid[strings.TrimSpace(value)] {
			return value
		}
		return ""
	}
	// multiple — filtrar cada parte
	var kept []string
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" && valid[t] {
			kept = append(kept, t)
		}
	}
	return strings.Join(kept, ",")
}

// filterAnswers aplica filterAnswerValue a un slice de SaveAnswerInput.
func filterAnswers(answers []SaveAnswerInput, validOptions map[string]map[string]bool) []SaveAnswerInput {
	out := make([]SaveAnswerInput, 0, len(answers))
	for _, a := range answers {
		a.Value = filterAnswerValue(a.QuestionID, a.Value, validOptions)
		out = append(out, a)
	}
	return out
}

// filterRepeaterAnswers aplica filterAnswerValue a las answers de cada entry.
func filterRepeaterAnswers(entries []SaveRepeaterEntryInput, validOptions map[string]map[string]bool) []SaveRepeaterEntryInput {
	for i := range entries {
		filtered := make([]SaveAnswerInput, 0, len(entries[i].Answers))
		for _, a := range entries[i].Answers {
			a.Value = filterAnswerValue(a.QuestionID, a.Value, validOptions)
			filtered = append(filtered, a)
		}
		entries[i].Answers = filtered
	}
	return entries
}

// pluralize retorna singular o plural según el conteo.
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

// processBarrierUpdateSubmission ejecuta la lógica de fin de formulario para actualización de barreras.
func (s *formService) processBarrierUpdateSubmission(ctx context.Context, submissionID string) error {
	// TODO: implementar lógica específica de barrier update
	return nil
}

// ─── validateAnswer ───────────────────────────────────────────────────────────

// ValidateAnswerResult es el resultado de validateAnswer.
type ValidateAnswerResult struct {
	Valid bool    `json:"valid"`
	Error *string `json:"error,omitempty"`
}

func validationError(msg string) ValidateAnswerResult {
	return ValidateAnswerResult{Valid: false, Error: &msg}
}

var validateAnswerOK = ValidateAnswerResult{Valid: true}

// validateAnswer evalúa si la respuesta de una pregunta es válida.
// Si la pregunta no es requerida, siempre retorna válido.
// Si es requerida, valida según el tipo de pregunta.
func validateAnswer(answer models.Answer, question QuestionStructure) ValidateAnswerResult {
	// PASO 1 — si no es requerida, siempre válida
	if !question.Required {
		return validateAnswerOK
	}

	val := strings.TrimSpace(answer.Value)

	// PASO 2 — validar según tipo
	switch question.QuestionTypeID {

	case "text":
		if val == "" {
			return validationError("Este campo es requerido")
		}

	case "number":
		if val == "" {
			return validationError("Este campo es requerido")
		}
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return validationError("Debe ser un número válido")
		}

	case "date":
		if val == "" {
			return validationError("Este campo es requerido")
		}
		if _, err := time.Parse("2006-01-02", val); err != nil {
			return validationError("Formato de fecha inválido (esperado YYYY-MM-DD)")
		}

	case "datetime":
		if val == "" {
			return validationError("Este campo es requerido")
		}
		// Intentar RFC3339 primero, luego el formato de <input datetime-local>
		if _, err := time.Parse(time.RFC3339, val); err != nil {
			if _, err2 := time.Parse("2006-01-02T15:04", val); err2 != nil {
				return validationError("Formato de fecha y hora inválido")
			}
		}

	case "single", "dropdown":
		if val == "" {
			return validationError("Debes seleccionar una opción")
		}

	case "boolean":
		if answer.Value != "true" && answer.Value != "false" {
			return validationError("Debes responder Sí o No")
		}

	case "multiple":
		if val == "" {
			return validationError("Selecciona al menos una opción")
		}

	case "info":
		// Los banners informativos no tienen respuesta — siempre válidos
		return validateAnswerOK

	// DEFAULT — tipo desconocido, no bloquear
	}

	return validateAnswerOK
}

// FormSection fue movido a form_section_service.go

// buildBarrierFollowUpSummary genera un texto resumen del seguimiento a una barrera.
func buildBarrierFollowUpSummary(persiste, respuestaInstitucional, actuaciones string) string {
	persisteStr := "No persiste"
	if persiste == "true" {
		persisteStr = "Persiste"
	}
	parts := []string{persisteStr}
	if respuestaInstitucional != "" {
		parts = append(parts, "Respuesta institucional: "+respuestaInstitucional)
	}
	if actuaciones != "" {
		parts = append(parts, actuaciones)
	}
	return strings.Join(parts, " · ")
}
