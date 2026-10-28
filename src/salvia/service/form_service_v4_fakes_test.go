package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type v4FollowFake struct {
	fu      *models.FollowUpV2
	updates int
}

func (f *v4FollowFake) Create(context.Context, *models.FollowUpV2) error { return nil }
func (f *v4FollowFake) Update(context.Context, *models.FollowUpV2) error { return nil }
func (f *v4FollowFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4FollowFake) Delete(context.Context, string) error { return nil }
func (f *v4FollowFake) FindByID(context.Context, string) (*models.FollowUpV2, error) {
	return f.fu, nil
}
func (f *v4FollowFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.FollowUpV2], error) {
	return repository.PageResult[models.FollowUpV2]{}, nil
}
func (f *v4FollowFake) FindByCaseID(context.Context, string) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (f *v4FollowFake) FindByCaseIDOrdered(context.Context, string) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (f *v4FollowFake) FindPendingByCaseID(context.Context, string) ([]models.FollowUpV2, error) {
	agent := "ag-1"
	return []models.FollowUpV2{{AgentID: &agent}}, nil
}
func (f *v4FollowFake) FindCompletedByCaseID(context.Context, string) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (f *v4FollowFake) BulkCreate(context.Context, *gorm.DB, []models.FollowUpV2) error { return nil }
func (f *v4FollowFake) SoftDeleteAndReprogramPending(context.Context, *gorm.DB, string) error {
	return nil
}
func (f *v4FollowFake) RunInTransaction(context.Context, func(*gorm.DB) error) error { return nil }
func (f *v4FollowFake) FindByAgentAndDate(context.Context, string, time.Time) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (f *v4FollowFake) FindRealizedTodayByAgent(context.Context, string, time.Time) ([]models.FollowUpV2, error) {
	return nil, nil
}
func (f *v4FollowFake) IncrementAttempt(context.Context, string) error { return nil }
func (f *v4FollowFake) FindByTeamPaginated(context.Context, string, repository.FollowUpFilters, int, int) ([]models.FollowUpV2, int64, error) {
	return nil, 0, nil
}
func (f *v4FollowFake) FindPendingByTeamGroupedByAgent(context.Context, string, string) ([]repository.AgentWorkload, error) {
	return nil, nil
}
func (f *v4FollowFake) FindAgentsByTeam(context.Context, string) ([]repository.AgentOption, error) {
	return nil, nil
}
func (f *v4FollowFake) Reschedule(context.Context, string, map[string]interface{}) error { return nil }
func (f *v4FollowFake) FindWorkloadByDates(context.Context, string, []time.Time) ([]repository.AgentDateWorkload, error) {
	return nil, nil
}
func (f *v4FollowFake) FindGlobalWorkloadByTeam(context.Context, string) ([]repository.AgentWorkload, error) {
	return nil, nil
}
func (f *v4FollowFake) CloseCaseFollowUps(context.Context, string) error    { return nil }
func (f *v4FollowFake) DeletePendingByCaseID(context.Context, string) error { return nil }
func (f *v4FollowFake) UpdateAgentForPendingByCaseID(context.Context, string, string) error {
	return nil
}
func (f *v4FollowFake) LoadVictimInfoByCaseID(context.Context, string) (*repository.VictimCaseInfo, error) {
	return nil, nil
}
func (f *v4FollowFake) UpdateFormSubmissionID(context.Context, string, string) error { return nil }
func (f *v4FollowFake) UpdateActiveBarrierIDs(context.Context, string, string) error { return nil }
func (f *v4FollowFake) UpdateFormIDAndSubmissionID(context.Context, string, string, string) error {
	return nil
}
func (f *v4FollowFake) FindByFormSubmissionID(context.Context, string) (*models.FollowUpV2, error) {
	return f.fu, nil
}
func (f *v4FollowFake) UpdateStatus(context.Context, string, string) error {
	f.updates++
	return nil
}
func (f *v4FollowFake) CreateTimelineEvent(context.Context, *models.CaseTimelineEvent) error {
	return nil
}

type v4CaseFake struct {
	level    int
	newLevel int
	levelSet bool
	team     string
	agent    string
}

func (f *v4CaseFake) FindByID(context.Context, string) (*models.VictimCaseLight, error) {
	return nil, nil
}
func (f *v4CaseFake) FindByICode(context.Context, string) (*models.VictimCaseLight, error) {
	return nil, nil
}
func (f *v4CaseFake) FindOpenByDocNumber(context.Context, string) (*models.VictimCaseLight, error) {
	return nil, nil
}
func (f *v4CaseFake) UpdateStatus(context.Context, string, string) error { return nil }
func (f *v4CaseFake) FindRiskLevelByICode(context.Context, string) (int, error) {
	return f.level, nil
}
func (f *v4CaseFake) UpdateRiskLevelByICode(_ context.Context, _ string, level int) error {
	f.levelSet = true
	f.newLevel = level
	return nil
}
func (f *v4CaseFake) UpdateTeamAndAgent(_ context.Context, _, team, agent string) error {
	f.team = team
	f.agent = agent
	return nil
}

type v4LetterFake struct{ rows []models.EntityLetter }

func (f *v4LetterFake) Create(_ context.Context, e *models.EntityLetter) error {
	e.ID = fmt.Sprintf("letter-%d", len(f.rows)+1)
	cp := *e
	if e.EntityBranchID != nil {
		v := *e.EntityBranchID
		cp.EntityBranchID = &v
	}
	f.rows = append(f.rows, cp)
	return nil
}
func (f *v4LetterFake) Update(context.Context, *models.EntityLetter) error { return nil }
func (f *v4LetterFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4LetterFake) Delete(context.Context, string) error { return nil }
func (f *v4LetterFake) FindByID(context.Context, string) (*models.EntityLetter, error) {
	return nil, nil
}
func (f *v4LetterFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.EntityLetter], error) {
	return repository.PageResult[models.EntityLetter]{}, nil
}
func (f *v4LetterFake) FindByCaseID(context.Context, string) ([]models.EntityLetter, error) {
	return nil, nil
}
func (f *v4LetterFake) FindByBarrierID(context.Context, string) ([]models.EntityLetter, error) {
	return nil, nil
}
func (f *v4LetterFake) FindByState(context.Context, string, int, int) (repository.PageResult[models.EntityLetter], error) {
	return repository.PageResult[models.EntityLetter]{}, nil
}
func (f *v4LetterFake) FindByAgentID(context.Context, string) ([]models.EntityLetter, error) {
	return nil, nil
}
func (f *v4LetterFake) FindByNotificationUserID(context.Context, string) ([]models.EntityLetter, error) {
	return nil, nil
}
func (f *v4LetterFake) FindByAgentIDWithRelations(context.Context, string) ([]models.EntityLetterWithRelations, error) {
	return nil, nil
}
func (f *v4LetterFake) FindByNotificationUserIDWithRelations(context.Context, string) ([]models.EntityLetterWithRelations, error) {
	return nil, nil
}
func (f *v4LetterFake) FindWithRelationsFilteredPaginated(context.Context, repository.EntityLetterListFilter, int, int) (repository.EntityLetterListResult, error) {
	return repository.EntityLetterListResult{}, nil
}
func (f *v4LetterFake) UpdateState(context.Context, string, string) error { return nil }

type v4TaskFake struct{ rows []models.CaseTask }

func (f *v4TaskFake) Create(_ context.Context, task *models.CaseTask) error {
	cp := *task
	copyPtr := func(src *string) *string {
		if src == nil {
			return nil
		}
		v := *src
		return &v
	}
	cp.EntityCaseID = copyPtr(task.EntityCaseID)
	cp.BarrierID = copyPtr(task.BarrierID)
	cp.EntityLetterID = copyPtr(task.EntityLetterID)
	cp.PsychosocialSupportID = copyPtr(task.PsychosocialSupportID)
	f.rows = append(f.rows, cp)
	return nil
}
func (f *v4TaskFake) Update(context.Context, *models.CaseTask) error { return nil }
func (f *v4TaskFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4TaskFake) Delete(context.Context, string) error { return nil }
func (f *v4TaskFake) FindByID(context.Context, string) (*models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.CaseTask], error) {
	return repository.PageResult[models.CaseTask]{}, nil
}
func (f *v4TaskFake) FindByAssignedUserID(context.Context, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByAssignedUserIDWithRelations(context.Context, string) ([]models.CaseTaskWithRelations, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByAssignedUserIDPaged(context.Context, string, string, string, int, int) (repository.PageResult[models.CaseTaskWithRelations], error) {
	return repository.PageResult[models.CaseTaskWithRelations]{}, nil
}
func (f *v4TaskFake) FindByCaseID(context.Context, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByCaseIDAndBarrier(context.Context, string, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByCaseIDAndPsychosocial(context.Context, string, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByReferralID(context.Context, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindByBarrierID(context.Context, string) ([]models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindTodoByEntityLetterID(context.Context, string) (*models.CaseTask, error) {
	return nil, nil
}
func (f *v4TaskFake) FindCommentsByCaseID(context.Context, string) ([]models.CaseTaskCommentItem, error) {
	return nil, nil
}

type v4BarrierFake struct {
	rows   []models.BarrierV2
	active []models.BarrierV2
	status map[string]string
}

func (f *v4BarrierFake) Create(_ context.Context, b *models.BarrierV2) error {
	b.ID = fmt.Sprintf("new-barrier-%d", len(f.rows)+1)
	cp := *b
	if b.EntityBranchID != nil {
		v := *b.EntityBranchID
		cp.EntityBranchID = &v
	}
	f.rows = append(f.rows, cp)
	return nil
}
func (f *v4BarrierFake) Update(context.Context, *models.BarrierV2) error { return nil }
func (f *v4BarrierFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4BarrierFake) Delete(context.Context, string) error { return nil }
func (f *v4BarrierFake) FindByID(context.Context, string) (*models.BarrierV2, error) {
	return nil, nil
}
func (f *v4BarrierFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.BarrierV2], error) {
	return repository.PageResult[models.BarrierV2]{}, nil
}
func (f *v4BarrierFake) FindByCaseID(context.Context, string) ([]models.BarrierV2, error) {
	return nil, nil
}
func (f *v4BarrierFake) FindByFollowUpID(context.Context, string) ([]models.BarrierV2, error) {
	return nil, nil
}
func (f *v4BarrierFake) FindActiveByCaseID(context.Context, string) ([]models.BarrierV2, error) {
	return f.active, nil
}
func (f *v4BarrierFake) FindByIDs(context.Context, []string) ([]models.BarrierV2, error) {
	return f.active, nil
}
func (f *v4BarrierFake) FindActiveByTeamContactIDs(context.Context, []string) ([]models.BarrierV2, error) {
	return nil, nil
}
func (f *v4BarrierFake) UpdateStatus(_ context.Context, id, status string) error {
	if f.status == nil {
		f.status = map[string]string{}
	}
	f.status[id] = status
	return nil
}
func (f *v4BarrierFake) FindByCreatedByIDWithRelations(context.Context, string) ([]models.BarrierV2WithRelations, error) {
	return nil, nil
}
func (f *v4BarrierFake) FindByDepartmentWithRelations(context.Context, string, repository.DepartmentBarrierFilters) ([]models.BarrierV2DepartmentItem, int64, error) {
	return nil, 0, nil
}
func (f *v4BarrierFake) FindMineWithRelations(context.Context, string) ([]models.BarrierV2MineItem, error) {
	return nil, nil
}
func (f *v4BarrierFake) GetDB() *gorm.DB { return nil }

type v4ReferralFake struct {
	referrals []models.EmergencyMeasureReferral
	types     [][]string
	riesgos   []bool
}

func (f *v4ReferralFake) CountByStatus(context.Context, string, string) (map[string]int, error) {
	return nil, nil
}
func (f *v4ReferralFake) List(context.Context, repository.EmergencyReferralListFilter) ([]repository.EmergencyReferralRow, error) {
	return nil, nil
}
func (f *v4ReferralFake) FindDetail(context.Context, string, string, string) (*repository.EmergencyReferralDetailRow, error) {
	return nil, nil
}
func (f *v4ReferralFake) UserHasActiveRole(context.Context, string, string) (bool, error) {
	return false, nil
}
func (f *v4ReferralFake) Reassign(context.Context, []string, string) error { return nil }
func (f *v4ReferralFake) Review(context.Context, string, string, string, string, string, string) (bool, bool, error) {
	return false, false, nil
}
func (f *v4ReferralFake) CreateWithMeasures(_ context.Context, referral *models.EmergencyMeasureReferral, measureTypes []string, riesgoInminente bool) error {
	if referral.ID == "" {
		referral.ID = fmt.Sprintf("ref-%d", len(f.referrals)+1)
	}
	cp := *referral
	if referral.FollowUpID != nil {
		id := *referral.FollowUpID
		cp.FollowUpID = &id
	}
	f.referrals = append(f.referrals, cp)
	copied := append([]string(nil), measureTypes...)
	f.types = append(f.types, copied)
	f.riesgos = append(f.riesgos, riesgoInminente)
	return nil
}

type v4MeasureFake struct{ rows []models.EmergencyMeasure }

func (f *v4MeasureFake) Create(_ context.Context, e *models.EmergencyMeasure) error {
	f.rows = append(f.rows, *e)
	return nil
}
func (f *v4MeasureFake) Update(context.Context, *models.EmergencyMeasure) error { return nil }
func (f *v4MeasureFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4MeasureFake) Delete(context.Context, string) error { return nil }
func (f *v4MeasureFake) FindByID(context.Context, string) (*models.EmergencyMeasure, error) {
	return nil, nil
}
func (f *v4MeasureFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.EmergencyMeasure], error) {
	return repository.PageResult[models.EmergencyMeasure]{}, nil
}
func (f *v4MeasureFake) FindByCaseID(context.Context, string) ([]models.EmergencyMeasure, error) {
	return nil, nil
}
func (f *v4MeasureFake) FindByFollowUpID(context.Context, string) ([]models.EmergencyMeasure, error) {
	return nil, nil
}

type v4PsicoFake struct{ rows []models.PsychosocialSupport }

func (f *v4PsicoFake) Create(_ context.Context, e *models.PsychosocialSupport) error {
	e.ID = fmt.Sprintf("ps-%d", len(f.rows)+1)
	f.rows = append(f.rows, *e)
	return nil
}
func (f *v4PsicoFake) Update(context.Context, *models.PsychosocialSupport) error { return nil }
func (f *v4PsicoFake) UpdateFields(context.Context, string, map[string]interface{}) error {
	return nil
}
func (f *v4PsicoFake) Delete(context.Context, string) error { return nil }
func (f *v4PsicoFake) FindByID(context.Context, string) (*models.PsychosocialSupport, error) {
	return nil, nil
}
func (f *v4PsicoFake) FindWithPagination(context.Context, int, int) (repository.PageResult[models.PsychosocialSupport], error) {
	return repository.PageResult[models.PsychosocialSupport]{}, nil
}
func (f *v4PsicoFake) FindByCaseID(context.Context, string) ([]models.PsychosocialSupport, error) {
	return nil, nil
}
func (f *v4PsicoFake) FindByFollowUpID(context.Context, string) ([]models.PsychosocialSupport, error) {
	return nil, nil
}
