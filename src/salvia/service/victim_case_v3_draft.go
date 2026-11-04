package service

// Borrador en JSON de Registro V3 (WP-4, REG-38*): mientras el Caso está en
// Borrador ('bo') las entidades y las barreras -y los metadatos que Seguimiento
// no guarda en columnas- viven en salvia.victim_case.victim_case_v3_draft_data.
// Cuando el Caso deja de ser borrador ('bo' -> 'in') flushDraft los materializa
// en entity_case, entity_case_obligation y barrier_v2, crea las tareas, los
// oficios y las remisiones (victim_case_v3_tasks.go) y deja la sección como un
// snapshot de solo lectura. Diseño completo:
// DocsMD/Screens/Registro V3/ (planeación WP-4, «Borrador en JSON»).

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ErrSnapshotReadOnly -- el Caso ya dejó de ser borrador: la sección de
// entidades, barreras y remisiones es un snapshot de solo lectura (D20). El
// controller la mapea a 409.
var ErrSnapshotReadOnly = errors.New("la sección de entidades y remisiones ya no es editable: el Caso dejó de ser borrador")

const (
	draftIDPrefixV3  = "d~"
	statusBorradorV3 = "bo"
	statusIniciadoV3 = "in"
)

// ── Forma del JSON ──────────────────────────────────────────────────────────

type DraftBranchV3 struct {
	ICode      string `json:"icode"`
	Name       string `json:"name"`
	Sector     string `json:"sector"`
	Address    string `json:"address"`
	Department string `json:"department"`
	City       string `json:"city"`
	Town       string `json:"town"`
}

type DraftEntityV3 struct {
	ID                 string        `json:"id"`
	Kind               string        `json:"kind"` // acudida | activacion
	EntityBranchID     int64         `json:"entityBranchId"`
	Branch             DraftBranchV3 `json:"branch"`
	FileNumber         string        `json:"fileNumber"`
	ReportDate         string        `json:"reportDate"`
	Objetivo           string        `json:"objetivo"`
	Solicitud          string        `json:"solicitud"`
	Pending            []string      `json:"pending"`
	Completed          []string      `json:"completed"`
	Agreement          bool          `json:"agreement"`
	DueDiligence       bool          `json:"dueDiligence"`
	ActivationChannels []string      `json:"activationChannels"`
	ManagementNotes    string        `json:"managementNotes"`
	CreatedAt          string        `json:"createdAt"`
	By                 string        `json:"by"`
}

type DraftBarrierV3 struct {
	ID                      string   `json:"id"`
	EntityBranchID          int64    `json:"entityBranchId"`
	Sector                  string   `json:"sector"`
	Specific                []string `json:"specific"`
	Institution             string   `json:"institution"`
	InstitutionName         string   `json:"institutionName"`
	OtherDetail             string   `json:"otherDetail"`
	Institutional           []string `json:"institutional"`
	Economic                []string `json:"economic"`
	Territorial             []string `json:"territorial"`
	Differential            []string `json:"differential"`
	Date                    string   `json:"date"`
	Official                string   `json:"official"`
	Department              string   `json:"department"`
	City                    string   `json:"city"`
	Town                    string   `json:"town"`
	Zone                    string   `json:"zone"`
	Description             string   `json:"description"`
	Management              []string `json:"management"`
	ProfessionalPerspective string   `json:"professionalPerspective"`
	ManagementDescription   string   `json:"managementDescription"`
	Competences             []string `json:"competences,omitempty"`
	CreatedAt               string   `json:"createdAt"`
	By                      string   `json:"by"`
}

// EntityMetaV3 / BarrierMetaV3 -- datos del formulario de Seguimiento v4 que no
// tienen columna en entity_case / barrier_v2 (D1 y D18). Se conservan después de
// materializar.
type EntityMetaV3 struct {
	Kind               string   `json:"kind"`
	Agreement          bool     `json:"agreement"`
	DueDiligence       bool     `json:"dueDiligence"`
	ActivationChannels []string `json:"activationChannels"`
	ReportDate         string   `json:"reportDate"`
	ManagementNotes    string   `json:"managementNotes"`
}

type BarrierMetaV3 struct {
	Zone                    string `json:"zone"`
	ProfessionalPerspective string `json:"professionalPerspective"`
	ManagementDescription   string `json:"managementDescription"`
}

type DraftV3 struct {
	V              int                      `json:"v"`
	Entities       []DraftEntityV3          `json:"entities"`
	Barriers       []DraftBarrierV3         `json:"barriers"`
	EntityMeta     map[string]EntityMetaV3  `json:"entityMeta"`
	BarrierMeta    map[string]BarrierMetaV3 `json:"barrierMeta"`
	MaterializedAt *string                  `json:"materializedAt"`
}

func newDraftV3() *DraftV3 {
	d := &DraftV3{V: 1}
	d.normalize()
	return d
}

func (d *DraftV3) normalize() {
	if d.V == 0 {
		d.V = 1
	}
	if d.Entities == nil {
		d.Entities = []DraftEntityV3{}
	}
	if d.Barriers == nil {
		d.Barriers = []DraftBarrierV3{}
	}
	if d.EntityMeta == nil {
		d.EntityMeta = map[string]EntityMetaV3{}
	}
	if d.BarrierMeta == nil {
		d.BarrierMeta = map[string]BarrierMetaV3{}
	}
}

// ── Dependencias y acceso a la columna JSON ─────────────────────────────────

// V3DraftDeps -- dependencias del borrador, de las tareas y de las remisiones.
// Se inyectan con WithVictimCaseV3Draft para no cambiar la firma de
// NewVictimCaseV3Service (la usan otras sesiones y los tests).
type V3DraftDeps struct {
	DB                           *gorm.DB
	CaseTaskRepo                 repository.CaseTaskRepository
	EntityLetterRepo             repository.EntityLetterRepository
	PsychosocialSupportRepo      repository.PsychosocialSupportRepository
	EconomicStabilizationRepo    repository.EconomicStabilizationRepository
	MenTeamRemisionRepo          repository.MenTeamRemisionRepository
	DiscapacidadRemisionRepo     repository.DiscapacidadRemisionRepository
	EmergencyMeasureReferralRepo repository.EmergencyMeasureReferralRepository
	AgentLightRepo               repository.AgentLightRepository
}

type victimCaseV3Draft struct {
	deps V3DraftDeps
}

// VictimCaseV3Option -- opción del constructor del servicio.
type VictimCaseV3Option func(*victimCaseV3Service)

// WithVictimCaseV3Draft activa el borrador en JSON: las entidades y barreras de
// un Caso en 'bo' van a victim_case_v3_draft_data y se materializan al pasar a
// 'in'. Sin esta opción el servicio conserva el comportamiento anterior.
func WithVictimCaseV3Draft(deps V3DraftDeps) VictimCaseV3Option {
	return func(s *victimCaseV3Service) {
		if deps.DB == nil {
			return
		}
		s.draft = &victimCaseV3Draft{deps: deps}
		// Decorador: el vaciado corre justo antes de cambiar bo -> in, sin editar
		// SaveSection (que es de WP-3). Si falla, el estado no cambia.
		s.repo = &draftFlushRepoV3{
			VictimCaseV3Repository: s.repo,
			flush: func(ctx context.Context, icode string) error {
				return s.flushDraft(ctx, icode, true)
			},
		}
	}
}

type draftFlushRepoV3 struct {
	repository.VictimCaseV3Repository
	flush func(ctx context.Context, icode string) error
}

func (r *draftFlushRepoV3) UpdateStatusIfCurrent(ctx context.Context, icode, expectedCurrent, newStatus string) error {
	if expectedCurrent == statusBorradorV3 && newStatus == statusIniciadoV3 && r.flush != nil {
		if err := r.flush(ctx, icode); err != nil {
			return fmt.Errorf("materializar el borrador: %w", err)
		}
	}
	return r.VictimCaseV3Repository.UpdateStatusIfCurrent(ctx, icode, expectedCurrent, newStatus)
}

type draftSnapshotV3 struct {
	Status  string
	AgentID string
	Draft   *DraftV3
}

func (dr *victimCaseV3Draft) load(tx *gorm.DB, icode string, lock bool) (*draftSnapshotV3, error) {
	var row struct {
		Status  string
		Data    string
		AgentID string
	}
	q := `SELECT COALESCE(victim_case_status, '') AS status,
	             COALESCE(victim_case_v3_draft_data::text, '') AS data,
	             COALESCE(agent_id, '') AS agent_id
	      FROM salvia.victim_case WHERE victim_case_i_code = ?`
	if lock {
		q += " FOR UPDATE"
	}
	res := tx.Raw(q, icode).Scan(&row)
	if res.Error != nil {
		return nil, fmt.Errorf("leer el borrador: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	d := newDraftV3()
	if strings.TrimSpace(row.Data) != "" {
		if err := json.Unmarshal([]byte(row.Data), d); err != nil {
			return nil, fmt.Errorf("leer el borrador: JSON inválido: %w", err)
		}
		d.normalize()
	}
	return &draftSnapshotV3{Status: row.Status, AgentID: row.AgentID, Draft: d}, nil
}

func (dr *victimCaseV3Draft) save(tx *gorm.DB, icode string, d *DraftV3) error {
	d.normalize()
	b, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("guardar el borrador: %w", err)
	}
	return tx.Exec(`UPDATE salvia.victim_case SET victim_case_v3_draft_data = ?::jsonb WHERE victim_case_i_code = ?`, string(b), icode).Error
}

// mutate lee el borrador con bloqueo de fila, exige que el Caso siga en 'bo' y
// aplica fn. Un Caso fuera de borrador devuelve ErrSnapshotReadOnly.
func (dr *victimCaseV3Draft) mutate(ctx context.Context, icode string, fn func(snap *draftSnapshotV3) error) error {
	return dr.deps.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		snap, err := dr.load(tx, icode, true)
		if err != nil {
			return err
		}
		if snap.Status != statusBorradorV3 {
			return ErrSnapshotReadOnly
		}
		if err := fn(snap); err != nil {
			return err
		}
		return dr.save(tx, icode, snap.Draft)
	})
}

func (dr *victimCaseV3Draft) read(ctx context.Context, icode string) (*draftSnapshotV3, error) {
	return dr.load(dr.deps.DB.WithContext(ctx), icode, false)
}

func (dr *victimCaseV3Draft) lookupBranch(ctx context.Context, branchID int64) (*DraftBranchV3, error) {
	var row struct {
		ICode      string
		Name       string
		Sector     string
		Address    string
		Department string
		City       string
		Town       string
	}
	res := dr.deps.DB.WithContext(ctx).Raw(`
SELECT eb.entity_branch_i_code AS i_code, eb.entity_branch_name AS name,
       COALESCE(e.entity_sector, '') AS sector, COALESCE(eb.entity_branch_address, '') AS address,
       COALESCE(d.department_name, '') AS department, COALESCE(c.city_name, '') AS city, COALESCE(t.town_name, '') AS town
FROM salvia.entity_branch eb
JOIN salvia.entity e ON e.entity_id = eb.entity_id
LEFT JOIN security.town t ON t.town_code = eb.entity_branch_town_code
LEFT JOIN security.city c ON c.city_id = t.city_id
LEFT JOIN security.department d ON d.department_id = c.department_id
WHERE eb.entity_branch_id = ?`, branchID).Scan(&row)
	if res.Error != nil {
		return nil, fmt.Errorf("buscar la sede %d: %w", branchID, res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("la sede %d no existe", branchID)
	}
	return &DraftBranchV3{ICode: row.ICode, Name: row.Name, Sector: row.Sector, Address: row.Address,
		Department: row.Department, City: row.City, Town: row.Town}, nil
}

// ── Helpers de normalización ────────────────────────────────────────────────

// canalesActivacionV3 -- las opciones de canal de Seguimiento v4 (D11): los 4 canales y
// «Ninguna», que se conserva para mostrarla pero no crea tareas (v4ProcesarCanales la ignora).
var canalesActivacionV3 = map[string]bool{ //nolint:gochecknoglobals
	"notificacion": true, "articulacion": true, "denuncia": true, "misalvia": true, "ninguna": true,
}

// canalActivacionLegadoV3 -- códigos del desplegable anterior de Registro V3
// (no, ll, ms) que un cliente viejo todavía podría enviar.
var canalActivacionLegadoV3 = map[string]string{ //nolint:gochecknoglobals
	"no": "notificacion", "ll": "articulacion", "ms": "misalvia",
}

func normalizarCanalesV3(channels []string, legacyCSV string) []string {
	raw := append([]string{}, channels...)
	if len(raw) == 0 && strings.TrimSpace(legacyCSV) != "" {
		raw = splitCSVV3(legacyCSV)
	}
	out := make([]string, 0, len(raw))
	seen := map[string]bool{}
	for _, c := range raw {
		c = strings.TrimSpace(c)
		if mapped, ok := canalActivacionLegadoV3[c]; ok {
			c = mapped
		}
		if !canalesActivacionV3[c] || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

func limpiarIDsV3(ids []string) []string {
	out := make([]string, 0, len(ids))
	seen := map[string]bool{}
	for _, raw := range ids {
		id := strings.TrimSpace(raw)
		if id == "" || id == "otra" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return out
}

func splitCSVV3(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func nuevoIDBorradorV3(icode string) string {
	return draftIDPrefixV3 + icode + "~" + utils.GetUUID()
}

func esIDBorradorV3(id string) bool { return strings.HasPrefix(id, draftIDPrefixV3) }

func nowRFC3339() string { return time.Now().UTC().Format(time.RFC3339) }

func ptrStr(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// ── Operaciones del borrador (Caso en 'bo') ─────────────────────────────────

func (s *victimCaseV3Service) addEntityDraft(ctx context.Context, e repository.VictimCaseEntityInput) (string, error) {
	if e.EntityBranchID == 0 {
		return "", fmt.Errorf("addEntity: falta la sede de la entidad")
	}
	// Primero el estado: un Caso que ya dejó de ser borrador es solo lectura
	// aunque la sede además esté repetida.
	current, err := s.draft.read(ctx, e.VictimCaseICode)
	if err != nil {
		return "", fmt.Errorf("addEntity: %w", err)
	}
	if current.Status != statusBorradorV3 {
		return "", ErrSnapshotReadOnly
	}
	// Una sede que ya está en las tablas (filas anteriores a este cambio) cuenta
	// como duplicada igual que antes.
	exists, err := s.entityCaseRepo.ExistsActive(ctx, e.VictimCaseICode, e.EntityBranchID)
	if err != nil {
		return "", fmt.Errorf("addEntity: verificar duplicado: %w", err)
	}
	if exists {
		return "", ErrEntityCaseDuplicate
	}
	branch, err := s.draft.lookupBranch(ctx, e.EntityBranchID)
	if err != nil {
		return "", fmt.Errorf("addEntity: %w", err)
	}
	kind := "acudida"
	if strings.TrimSpace(e.Kind) == "activacion" {
		kind = "activacion"
	}
	completed := limpiarIDsV3(e.CompletedActions)
	if kind != "acudida" {
		completed = []string{}
	}
	entity := DraftEntityV3{
		ID:                 nuevoIDBorradorV3(e.VictimCaseICode),
		Kind:               kind,
		EntityBranchID:     e.EntityBranchID,
		Branch:             *branch,
		FileNumber:         strings.TrimSpace(e.FileNumber),
		ReportDate:         strings.TrimSpace(e.ReportDate),
		Objetivo:           strings.TrimSpace(e.OtherInfo),
		Solicitud:          strings.TrimSpace(e.Solicitud),
		Pending:            limpiarIDsV3(e.PendingActions),
		Completed:          completed,
		Agreement:          e.Agreement,
		DueDiligence:       e.DueDiligence,
		ActivationChannels: normalizarCanalesV3(e.ActivationChannels, e.ActivationChannel),
		ManagementNotes:    strings.TrimSpace(e.ManagementNotes),
		CreatedAt:          nowRFC3339(),
		By:                 e.Agent,
	}
	err = s.draft.mutate(ctx, e.VictimCaseICode, func(snap *draftSnapshotV3) error {
		for _, existing := range snap.Draft.Entities {
			if existing.EntityBranchID == e.EntityBranchID {
				return ErrEntityCaseDuplicate
			}
		}
		snap.Draft.Entities = append(snap.Draft.Entities, entity)
		return nil
	})
	if err != nil {
		return "", err
	}
	return entity.ID, nil
}

func (s *victimCaseV3Service) removeEntityDraft(ctx context.Context, id string) error {
	icode, ok := icodeDeIDBorradorV3(id)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	return s.draft.mutate(ctx, icode, func(snap *draftSnapshotV3) error {
		idx := -1
		var branch int64
		for i, ent := range snap.Draft.Entities {
			if ent.ID == id {
				idx, branch = i, ent.EntityBranchID
				break
			}
		}
		if idx < 0 {
			return gorm.ErrRecordNotFound
		}
		snap.Draft.Entities = append(snap.Draft.Entities[:idx], snap.Draft.Entities[idx+1:]...)
		// Cascada solo en el JSON (D13): se van las barreras de esa sede y, con
		// ellas, todo lo que se iba a crear (tareas, oficios y acciones).
		kept := snap.Draft.Barriers[:0]
		for _, b := range snap.Draft.Barriers {
			if b.EntityBranchID != branch {
				kept = append(kept, b)
			}
		}
		snap.Draft.Barriers = kept
		return nil
	})
}

func icodeDeIDBorradorV3(id string) (string, bool) {
	if !esIDBorradorV3(id) {
		return "", false
	}
	parts := strings.Split(id, "~")
	if len(parts) < 3 || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func (s *victimCaseV3Service) addBarrierDraft(ctx context.Context, b repository.VictimCaseBarrierInput) (string, error) {
	management := make([]string, 0, len(b.Management))
	for _, code := range b.Management {
		code = strings.TrimSpace(code)
		if code == "" {
			continue
		}
		// Cliente anterior: códigos 01..06 del desplegable viejo -> slugs.
		if slug, ok := gestionBarreraV3CodeToSlug[code]; ok {
			code = slug
		}
		management = append(management, code)
	}
	barrier := DraftBarrierV3{
		ID:                      nuevoIDBorradorV3(b.VictimCaseICode),
		EntityBranchID:          b.EntityBranchID,
		Sector:                  b.SectorCode,
		Specific:                append([]string{}, b.Specific...),
		Institution:             strings.TrimSpace(b.Institution),
		InstitutionName:         strings.TrimSpace(b.EntityName),
		OtherDetail:             strings.TrimSpace(b.OtherDetail),
		Institutional:           append([]string{}, b.Institutional...),
		Economic:                append([]string{}, b.Economic...),
		Territorial:             append([]string{}, b.Territorial...),
		Differential:            append([]string{}, b.Differential...),
		Date:                    strings.TrimSpace(b.Date),
		Official:                strings.TrimSpace(b.Official),
		Department:              strings.TrimSpace(b.Department),
		City:                    strings.TrimSpace(b.City),
		Town:                    strings.TrimSpace(b.Town),
		Zone:                    strings.TrimSpace(b.Zone),
		Description:             strings.TrimSpace(b.Description),
		Management:              management,
		ProfessionalPerspective: strings.TrimSpace(b.ProfessionalPerspective),
		ManagementDescription:   strings.TrimSpace(b.ManagementDescription),
		Competences:             append([]string{}, b.Competences...),
		CreatedAt:               nowRFC3339(),
		By:                      b.Agent,
	}
	if err := s.draft.mutate(ctx, b.VictimCaseICode, func(snap *draftSnapshotV3) error {
		snap.Draft.Barriers = append(snap.Draft.Barriers, barrier)
		return nil
	}); err != nil {
		return "", err
	}
	return barrier.ID, nil
}

func (s *victimCaseV3Service) removeBarrierDraft(ctx context.Context, id string) error {
	icode, ok := icodeDeIDBorradorV3(id)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	return s.draft.mutate(ctx, icode, func(snap *draftSnapshotV3) error {
		for i, br := range snap.Draft.Barriers {
			if br.ID == id {
				snap.Draft.Barriers = append(snap.Draft.Barriers[:i], snap.Draft.Barriers[i+1:]...)
				return nil
			}
		}
		return gorm.ErrRecordNotFound
	})
}

// ── Lectura: une el JSON con las tablas ─────────────────────────────────────

// reintentarVaciadoV3 -- si el Caso ya salió de borrador pero quedó JSON
// pendiente (un vaciado anterior falló, o una escritura llegó justo antes del
// cambio de estado), lo vuelve a materializar. Idempotente; no crea remisiones.
func (s *victimCaseV3Service) reintentarVaciadoV3(ctx context.Context, icode string, snap *draftSnapshotV3) *draftSnapshotV3 {
	if snap.Status == statusBorradorV3 || (len(snap.Draft.Entities) == 0 && len(snap.Draft.Barriers) == 0) {
		return snap
	}
	if err := s.flushDraft(ctx, icode, false); err != nil {
		log.Printf("[victimCaseV3] WARN: no se pudo reintentar el vaciado del borrador de %s: %v", icode, err)
		return snap
	}
	fresh, err := s.draft.read(ctx, icode)
	if err != nil {
		return snap
	}
	return fresh
}

func (s *victimCaseV3Service) draftEntitiesAsRows(ctx context.Context, icode string) ([]models.EntityCaseWithRelations, error) {
	snap, err := s.draft.read(ctx, icode)
	if err != nil {
		return nil, err
	}
	snap = s.reintentarVaciadoV3(ctx, icode, snap)
	barrierCount := map[int64]int64{}
	for _, b := range snap.Draft.Barriers {
		barrierCount[b.EntityBranchID]++
	}
	rows := make([]models.EntityCaseWithRelations, 0, len(snap.Draft.Entities))
	for _, e := range snap.Draft.Entities {
		created, _ := time.Parse(time.RFC3339, e.CreatedAt)
		rows = append(rows, models.EntityCaseWithRelations{
			RelID: e.ID, EntityBranchID: e.EntityBranchID, EntityBranchICode: e.Branch.ICode,
			EntityBranchName: e.Branch.Name, Sector: e.Branch.Sector, Address: e.Branch.Address,
			DepartmentName: e.Branch.Department, CityName: e.Branch.City, TownName: e.Branch.Town,
			Objetivo: ptrStr(e.Objetivo), Solicitud: ptrStr(e.Solicitud), Radicado: ptrStr(e.FileNumber),
			BarrerasActivasCount: barrierCount[e.EntityBranchID],
			CreatedByID:          e.By, CreatedAt: created,
		})
	}
	return rows, nil
}

func (s *victimCaseV3Service) draftBarriersAsRows(ctx context.Context, icode string) ([]models.BarrierV2, error) {
	snap, err := s.draft.read(ctx, icode)
	if err != nil {
		return nil, err
	}
	snap = s.reintentarVaciadoV3(ctx, icode, snap)
	rows := make([]models.BarrierV2, 0, len(snap.Draft.Barriers))
	for _, b := range snap.Draft.Barriers {
		rows = append(rows, barrierModelFromDraftV3(icode, b))
	}
	return rows, nil
}

func barrierModelFromDraftV3(icode string, b DraftBarrierV3) models.BarrierV2 {
	var branch *int64
	if b.EntityBranchID != 0 {
		id := b.EntityBranchID
		branch = &id
	}
	created, _ := time.Parse(time.RFC3339, b.CreatedAt)
	return models.BarrierV2{
		ID: b.ID, CaseID: icode, CreatedByID: b.By, Status: models.BarrierV2StatusOpen,
		Sector: b.Sector, SpecificBarriers: strings.Join(b.Specific, ","), SpecificInstitutions: b.Institution,
		OtherBarrierDesc: b.OtherDetail, InstitutionName: b.InstitutionName,
		DepartmentID: b.Department, CityID: b.City, TownID: b.Town,
		StructuralInstitutional: strings.Join(b.Institutional, ","), StructuralEconomic: strings.Join(b.Economic, ","),
		StructuralTerritorial: strings.Join(b.Territorial, ","), StructuralDifferential: strings.Join(b.Differential, ","),
		BarrierDate: b.Date, OfficialDependency: b.Official, Description: b.Description,
		ManagementActions: strings.Join(b.Management, ","), Competences: strings.Join(b.Competences, ","),
		EntityBranchID: branch, CreatedAt: created,
	}
}

// draftMetaV3 devuelve los metadatos que GetCase entrega junto a las entidades y
// barreras, y si la sección es de solo lectura (el Caso ya no es borrador).
func (s *victimCaseV3Service) draftMetaV3(ctx context.Context, icode string) (map[string]EntityMetaV3, map[string]BarrierMetaV3, bool, error) {
	snap, err := s.draft.read(ctx, icode)
	if err != nil {
		return nil, nil, false, err
	}
	entityMeta := map[string]EntityMetaV3{}
	for k, v := range snap.Draft.EntityMeta {
		entityMeta[k] = v
	}
	// Las entidades que todavía están en el JSON también llevan su meta.
	for _, e := range snap.Draft.Entities {
		entityMeta[strconv.FormatInt(e.EntityBranchID, 10)] = entityMetaFromDraftV3(e)
	}
	barrierMeta := map[string]BarrierMetaV3{}
	for k, v := range snap.Draft.BarrierMeta {
		barrierMeta[k] = v
	}
	for _, b := range snap.Draft.Barriers {
		barrierMeta[b.ID] = BarrierMetaV3{Zone: b.Zone, ProfessionalPerspective: b.ProfessionalPerspective, ManagementDescription: b.ManagementDescription}
	}
	return entityMeta, barrierMeta, snap.Status != statusBorradorV3, nil
}

func entityMetaFromDraftV3(e DraftEntityV3) EntityMetaV3 {
	return EntityMetaV3{Kind: e.Kind, Agreement: e.Agreement, DueDiligence: e.DueDiligence,
		ActivationChannels: e.ActivationChannels, ReportDate: e.ReportDate, ManagementNotes: e.ManagementNotes}
}

// ── Vaciado («flush») ───────────────────────────────────────────────────────

type flushedEntityV3 struct {
	Draft        DraftEntityV3
	EntityCaseID string
}

type flushedBarrierV3 struct {
	Draft     DraftBarrierV3
	BarrierID string
}

// flushDraft materializa el borrador.
//
//	transition=true  (decorador de UpdateStatusIfCurrent bo -> in): exige que el
//	  Caso siga en 'bo'; materializa entidades y barreras pendientes y, la primera
//	  vez, crea las remisiones a otros equipos con las respuestas guardadas.
//	transition=false (reintento al leer un Caso que ya salió de borrador): solo
//	  materializa entidades y barreras pendientes.
//
// Las filas se crean en una sola transacción; las tareas, los oficios y las
// remisiones se crean después del COMMIT con mejor esfuerzo (como Seguimiento).
func (s *victimCaseV3Service) flushDraft(ctx context.Context, icode string, transition bool) error {
	if s.draft == nil {
		return nil
	}
	var (
		entities []flushedEntityV3
		barriers []flushedBarrierV3
		doRemis  bool
		agentID  string
	)
	err := s.draft.deps.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		snap, err := s.draft.load(tx, icode, true)
		if err != nil {
			return err
		}
		if transition && snap.Status != statusBorradorV3 {
			return nil // otro guardado ya hizo la transición
		}
		d := snap.Draft
		doRemis = transition && d.MaterializedAt == nil
		agentID = snap.AgentID
		if len(d.Entities) == 0 && len(d.Barriers) == 0 && !doRemis {
			return nil
		}
		for _, e := range d.Entities {
			entityCaseID, err := s.materializarEntidadV3(tx, icode, e)
			if err != nil {
				return err
			}
			d.EntityMeta[strconv.FormatInt(e.EntityBranchID, 10)] = entityMetaFromDraftV3(e)
			entities = append(entities, flushedEntityV3{Draft: e, EntityCaseID: entityCaseID})
		}
		for _, b := range d.Barriers {
			model := barrierModelFromDraftV3(icode, b)
			model.ID = "" // lo genera la BD
			model.CreatedAt = time.Time{}
			if err := tx.Create(&model).Error; err != nil {
				return fmt.Errorf("crear barrier_v2: %w", err)
			}
			d.BarrierMeta[model.ID] = BarrierMetaV3{Zone: b.Zone, ProfessionalPerspective: b.ProfessionalPerspective, ManagementDescription: b.ManagementDescription}
			barriers = append(barriers, flushedBarrierV3{Draft: b, BarrierID: model.ID})
		}
		d.Entities = []DraftEntityV3{}
		d.Barriers = []DraftBarrierV3{}
		if doRemis {
			now := nowRFC3339()
			d.MaterializedAt = &now
		}
		return s.draft.save(tx, icode, d)
	})
	if err != nil {
		return err
	}
	s.crearTareasMaterializacionV3(ctx, icode, entities, barriers)
	if doRemis {
		s.crearRemisionesV3(ctx, icode, agentID)
	}
	return nil
}

// materializarEntidadV3 crea el entity_case (o reutiliza el existente de esa
// sede) y sus obligaciones, dentro de la transacción del vaciado.
func (s *victimCaseV3Service) materializarEntidadV3(tx *gorm.DB, icode string, e DraftEntityV3) (string, error) {
	var existingID string
	res := tx.Raw(`SELECT id FROM salvia.entity_case WHERE case_id = ? AND entity_branch_id = ? AND deleted_at IS NULL LIMIT 1`,
		icode, e.EntityBranchID).Scan(&existingID)
	if res.Error != nil {
		return "", fmt.Errorf("buscar entity_case existente: %w", res.Error)
	}
	entityCaseID := existingID
	if entityCaseID == "" {
		ec := &models.EntityCase{
			CaseID: icode, EntityBranchID: e.EntityBranchID,
			Objetivo: ptrStr(e.Objetivo), Solicitud: ptrStr(e.Solicitud), Radicado: ptrStr(e.FileNumber),
			CreatedByID: e.By,
		}
		if t, err := time.Parse(time.RFC3339, e.CreatedAt); err == nil {
			ec.CreatedAt = t
		}
		if err := tx.Create(ec).Error; err != nil {
			return "", fmt.Errorf("crear entity_case: %w", err)
		}
		entityCaseID = ec.ID
	}
	// Si la sede ya estaba en las tablas se reutiliza (no se duplica) y solo se
	// agregan las acciones que todavía no tenga.
	crear := func(ids []string, status string) error {
		for _, id := range limpiarIDsV3(ids) {
			var n int64
			if err := tx.Raw(`SELECT COUNT(*) FROM salvia.entity_case_obligation WHERE entity_case_id = ? AND entity_obligation_id = ?`,
				entityCaseID, id).Scan(&n).Error; err != nil {
				return fmt.Errorf("verificar entity_case_obligation: %w", err)
			}
			if n > 0 {
				continue
			}
			oid := id
			eco := &models.EntityCaseObligation{
				EntityCaseID: entityCaseID, EntityObligationID: &oid, Status: status, CreatedByID: e.By,
			}
			if err := tx.Create(eco).Error; err != nil {
				return fmt.Errorf("crear entity_case_obligation: %w", err)
			}
		}
		return nil
	}
	if err := crear(e.Completed, models.EntityCaseObligationStatusCompletada); err != nil {
		return "", err
	}
	if err := crear(e.Pending, models.EntityCaseObligationStatusPendiente); err != nil {
		return "", err
	}
	return entityCaseID, nil
}
