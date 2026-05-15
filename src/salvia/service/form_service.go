package service

import (
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
	Options    []models.Option              `json:"options"`
	Conditions []models.VisibilityCondition `json:"conditions"`
}

type RepeaterStructure struct {
	models.RepeaterGroup
	Questions  []QuestionStructure          `json:"questions"`
	Conditions []models.VisibilityCondition `json:"conditions"`
}

type SectionStructure struct {
	models.FormSection
	Questions  []QuestionStructure          `json:"questions"`
	Repeaters  []RepeaterStructure          `json:"repeaters"`
	Conditions []models.VisibilityCondition `json:"conditions"`
	IsAnswered bool                         `json:"isAnswered"`
	IsVisible  bool                         `json:"isVisible"`
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
	LoadForm(ctx context.Context, formID, submissionID string) (*LoadFormResult, error)
	SaveSection(ctx context.Context, input SaveSectionInput) (*LoadFormResult, error)
	OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error
	TestFunction(ctx context.Context, fn, id, submissionID string) (interface{}, error)
}

type FormServiceDeps struct {
	FormRepo                   repository.FormRepository
	FormSectionRepo            repository.FormSectionRepository
	QuestionRepo               repository.QuestionRepository
	RepeaterGroupRepo          repository.RepeaterGroupRepository
	OptionRepo                 repository.OptionRepository
	VisibilityCondRepo         repository.VisibilityConditionRepository
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
}

type formService struct {
	repo                   repository.FormRepository
	formSectionRepo        repository.FormSectionRepository
	questionRepo           repository.QuestionRepository
	repeaterGroupRepo      repository.RepeaterGroupRepository
	optionRepo             repository.OptionRepository
	visibilityCondRepo     repository.VisibilityConditionRepository
	submissionRepo         repository.FormSubmissionRepository
	repeaterEntryRepo      repository.RepeaterEntryRepository
	answerRepo             repository.AnswerRepository
	followUpRepo           repository.FollowUpRepository
	emergencyMeasureRepo   repository.EmergencyMeasureRepository
	psychosocialSupportRepo repository.PsychosocialSupportRepository
	economicStabilizationRepo repository.EconomicStabilizationRepository
	barrierV2Repo          repository.BarrierV2Repository
	caseTimelineRepo       repository.CaseTimelineEventRepository
	agentLightRepo         repository.AgentLightRepository
}

func NewFormService(deps FormServiceDeps) FormService {
	return &formService{
		repo:                      deps.FormRepo,
		formSectionRepo:           deps.FormSectionRepo,
		questionRepo:              deps.QuestionRepo,
		repeaterGroupRepo:         deps.RepeaterGroupRepo,
		optionRepo:                deps.OptionRepo,
		visibilityCondRepo:        deps.VisibilityCondRepo,
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

	// ── Build lookup maps ──────────────────────────────────────────────────────

	// options by questionID
	optionsByQuestion := map[string][]models.Option{}
	for _, o := range options {
		optionsByQuestion[o.QuestionID] = append(optionsByQuestion[o.QuestionID], o)
	}

	// VCs by composite key "TargetType:TargetID" — validates both fields
	vcByKey := map[string][]models.VisibilityCondition{}
	for _, vc := range vcs {
		key := vc.TargetType + ":" + vc.TargetID
		vcByKey[key] = append(vcByKey[key], vc)
	}

	// questions by sectionID (section-level) and by repeaterGroupID
	sectionQuestions  := map[string][]QuestionStructure{}
	repeaterQuestions := map[string][]QuestionStructure{}
	for _, q := range questions {
		qs := QuestionStructure{
			Question:   q,
			Options:    optionsByQuestion[q.ID],
			Conditions: vcByKey["QUESTION:"+q.ID],
		}
		if qs.Options == nil    { qs.Options = []models.Option{} }
		if qs.Conditions == nil { qs.Conditions = []models.VisibilityCondition{} }

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
		}
		if rs.Questions == nil  { rs.Questions = []QuestionStructure{} }
		if rs.Conditions == nil { rs.Conditions = []models.VisibilityCondition{} }
		repeatersBySection[rg.FormSectionID] = append(repeatersBySection[rg.FormSectionID], rs)
	}

	// ── Assemble result ────────────────────────────────────────────────────────
	sectionStructures := make([]SectionStructure, len(sections))
	for i, sec := range sections {
		qs := sectionQuestions[sec.ID]
		rs := repeatersBySection[sec.ID]
		cs := vcByKey["SECTION:"+sec.ID]
		if qs == nil { qs = []QuestionStructure{} }
		if rs == nil { rs = []RepeaterStructure{} }
		if cs == nil { cs = []models.VisibilityCondition{} }
		sectionStructures[i] = SectionStructure{
			FormSection: sec,
			Questions:   qs,
			Repeaters:   rs,
			Conditions:  cs,
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
		secVis := checkVisibility(fs, sub, nil, sectionItem(sec))

		// Preguntas directas con su answer y visibilidad
		questions := make([]QuestionStructured, len(sec.Questions))
		for j, q := range sec.Questions {
			qVis := checkVisibility(fs, sub, nil, directQuestionItem(q, sec.Order))
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
			rVis := checkVisibility(fs, sub, nil, repeaterItem(r, sec.Order))

			rawEntries := entriesByGroup[r.ID]
			entries := make([]RepeaterEntryStructured, len(rawEntries))
			for l, e := range rawEntries {
				entryQuestions := make([]RepeaterEntryQuestion, len(r.Questions))
				for m, q := range r.Questions {
					qVis := checkVisibility(fs, sub, e.Answers, repeaterQuestionItem(q, sec.Order, r.Order, r.ID))
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
			result["section:"+sec.ID] = checkVisibility(fs, submission, nil, sectionItem(sec))
			for _, q := range sec.Questions {
				result["question:"+q.ID] = checkVisibility(fs, submission, nil, directQuestionItem(q, sec.Order))
			}
			for _, r := range sec.Repeaters {
				result["repeater:"+r.ID] = checkVisibility(fs, submission, nil, repeaterItem(r, sec.Order))
				for _, q := range r.Questions {
					entries := entriesByGroup[r.ID]
					if len(entries) == 0 {
						// Sin entries: evaluar sin entryAnswers
						result["repeaterQ:"+q.ID+":noEntry"] = checkVisibility(fs, submission, nil, repeaterQuestionItem(q, sec.Order, r.Order, r.ID))
					} else {
						// Evaluar para cada entry
						for _, entry := range entries {
							key := "repeaterQ:" + q.ID + ":entry" + entry.ID
							result[key] = checkVisibility(fs, submission, entry.Answers, repeaterQuestionItem(q, sec.Order, r.Order, r.ID))
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
			got := isAnsweredQuestion(qwo.q, qwo.order, c.answer, fs, c.sub)
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
			got  := isAnsweredRepeater(r, sectionOrder, entries, fs, sub)
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
			got  := isAnsweredSection(secByOrder[secOrder], fs, sub)
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
			res, e := s.LoadForm(ctx, id, subID)
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
func checkVisibility(
	fs *FormStructure,
	submission *SubmissionStructure,
	entryAnswers []models.Answer,
	item visibilityItem,
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

		trigVal := ""
		if cond.TriggerValue != nil {
			trigVal = *cond.TriggerValue
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
) IsAnsweredResult {
	// PASO 1 — evaluar visibilidad
	vis := checkVisibility(fs, sub, nil, directQuestionItem(q, sectionOrder))

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
) IsAnsweredRepeaterResult {
	// PASO 1 — visibilidad del repeater group
	vis := checkVisibility(fs, sub, nil, repeaterItem(r, sectionOrder))
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
			qVis := checkVisibility(fs, sub, entry.Answers, repeaterQuestionItem(q, sectionOrder, r.Order, r.ID))
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
	sec SectionStructure,
	fs  *FormStructure,
	sub *SubmissionStructure,
) IsAnsweredResult {
	// PASO 1 — visibilidad de la sección
	vis := checkVisibility(fs, sub, nil, sectionItem(sec))
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
		result := isAnsweredQuestion(q, sec.Order, answer, fs, sub)
		if !result.IsAnswered {
			return IsAnsweredResult{IsVisible: true, IsAnswered: false}
		}
	}

	// PASO 4 — iterar repeaters
	for _, r := range sec.Repeaters {
		entries := entriesByGroup[r.ID]
		result  := isAnsweredRepeater(r, sec.Order, entries, fs, sub)
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
func (s *formService) LoadForm(ctx context.Context, formID, submissionID string) (*LoadFormResult, error) {
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
		result := isAnsweredSection(*sec, fs, sub)
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
	ActorID          string                   `json:"actorId"` // general_user_i_code — inyectado por el controller desde la sesión
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

	// 5. Retornar LoadForm con el estado actualizado (isAnswered/isVisible por sección)
	result, err := s.LoadForm(ctx, input.FormID, submissionID)
	if err != nil {
		return nil, err
	}

	// 6. Si todas las secciones visibles están respondidas → onEndFormSubmission
	allAnswered := true
	for _, sec := range result.FormStructure.Sections {
		if sec.IsVisible && !sec.IsAnswered {
			allAnswered = false
			break
		}
	}
	if allAnswered {
		actorID := input.ActorID
		go func() {
			if err := s.OnEndFormSubmission(context.Background(), input.FormID, submissionID, actorID); err != nil {
				fmt.Printf("[OnEndFormSubmission] error: %v\n", err)
			}
		}()
	}

	return result, nil
}

// ─── OnEndFormSubmission ──────────────────────────────────────────────────────

// IDs de formularios con lógica de end-submission.
const (
	seguimientoFormID    = "2d0aeb46-1af3-4c47-a0d5-c5bfc4d549ff"
	barrierUpdateFormID  = "4d0aeb46-5af3-4c47-a0d5-c5bfc4d549ff"
)

// OnEndFormSubmission es llamado cuando todas las secciones visibles han sido respondidas.
// Delega a la función específica según el formID.
func (s *formService) OnEndFormSubmission(ctx context.Context, formID, submissionID, actorID string) error {
	switch formID {
	case seguimientoFormID:
		return s.processFollowUpSubmission(ctx, submissionID, actorID)
	case barrierUpdateFormID:
		return s.processBarrierUpdateSubmission(ctx, submissionID)
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
		// Repeater group de barreras
		rgBarreras = "5fd3ecdc-2e5f-4b31-97ef-8a994580586a"
		// Preguntas dentro del repeater de barreras
		qBarreraSector     = "f19378b6-55c5-4fdf-b765-7ebcc3978741" // Sector (dropdown)
		qBarreraSalud      = "5fc1f2af-cc30-41f0-aa31-4731e5cb674c" // Barreras Salud (multiple)
		qBarreraJusticia   = "2bec977e-97c7-42d7-a00a-b536af8038eb" // Barreras Justicia (multiple)
		qBarreraProteccion = "66c9fc1e-9b5e-4ad4-999f-7aeb483d84dc" // Barreras Protección (multiple)
		// Preguntas directas del form
		qEquipos           = "e0d38cf5-fe3f-45cb-9fd3-f5b8f7b2f7dc" // Derivaciones a equipos (multi-select)
		qMedidasEmergencia = "1a36260c-33a4-4ebd-bffb-e387d7964b96" // Medidas de emergencia (multi-select)
	)

	// Mapa: questionID → nombre del sector para las preguntas de barrera
	barrierQuestionSector := map[string]string{
		qBarreraSalud:      "salud",
		qBarreraJusticia:   "justicia",
		qBarreraProteccion: "proteccion",
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
	barrierCount   := 0
	remisionCount  := 0

	// 3. Crear BarrierV2 por cada entrada del repeater de barreras
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
		// Indexar respuestas de esta entrada
		entryMap := make(map[string]string, len(entryAnswers))
		for _, a := range entryAnswers {
			entryMap[a.QuestionID] = a.Value
		}
		// Sector explícito (dropdown); si no hay, se inferirá de la pregunta con respuesta
		sectorExplicit := strings.TrimSpace(entryMap[qBarreraSector])

		for qID, sectorFallback := range barrierQuestionSector {
			val, ok := entryMap[qID]
			if !ok || val == "" {
				continue
			}
			sector := sectorExplicit
			if sector == "" {
				sector = sectorFallback
			}
			for _, option := range splitValues(val) {
				log.Printf("[processFollowUp] creando barrera sector=%s descripcion=%s", sector, option)
				b := &models.BarrierV2{
					CaseID:      fu.CaseID,
					FollowUpID:  fu.ID,
					Sector:      sector,
					Description: option,
					Status:      "OPEN",
				}
				if err := s.barrierV2Repo.Create(ctx, b); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear barrera [%s/%s]: %w", sector, option, err)
				}
				barrierCount++
			}
		}
	}

	// 4. Crear registros de derivación por equipo
	if equiposVal, ok := answerMap[qEquipos]; ok && equiposVal != "" {
		log.Printf("[processFollowUp] derivaciones a equipos: %s", equiposVal)
		for _, equipo := range splitValues(equiposVal) {
			switch equipo {
			case "atencion_psico":
				log.Printf("[processFollowUp] creando derivacion -> atencion_psico")
				ps := &models.PsychosocialSupport{
					CaseID:     fu.CaseID,
					FollowUpID: fu.ID,
					Type:       "derivacion",
					Status:     "ACTIVE",
				}
				if err := s.psychosocialSupportRepo.Create(ctx, ps); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear derivacion psicosocial: %w", err)
				}
				remisionCount++
			case "estabilizacion":
				log.Printf("[processFollowUp] creando derivacion -> estabilizacion")
				ec := &models.EconomicStabilization{
					CaseID:     fu.CaseID,
					FollowUpID: fu.ID,
					Type:       "derivacion",
					Status:     "ACTIVE",
				}
				if err := s.economicStabilizationRepo.Create(ctx, ec); err != nil {
					return fmt.Errorf("processFollowUpSubmission: crear derivacion economica: %w", err)
				}
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

	return nil
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

	// DEFAULT — tipo desconocido, no bloquear
	}

	return validateAnswerOK
}

// FormSection fue movido a form_section_service.go
