package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"testing"

	"gorm.io/gorm"
)

// fakeDuplaRepoConflict embeds la interfaz para no implementar todos los métodos.
type fakeDuplaRepoConflict struct {
	repository.DuplaRepository
	occupied map[string]bool
}

func (f *fakeDuplaRepoConflict) ExistsActiveByName(ctx context.Context, name, excludeID string) (bool, error) {
	return false, nil
}

func (f *fakeDuplaRepoConflict) ExistsActiveMemberConflict(ctx context.Context, memberID, excludeID string) (bool, error) {
	return f.occupied[memberID], nil
}

type fakeReassignRepoConflict struct {
	repository.PsychosocialReassignRepository
	roles map[string]string
}

func (f *fakeReassignRepoConflict) FindProfessionalForReassign(ctx context.Context, icode string) (*models.PsTsProfessionalRow, error) {
	role, ok := f.roles[icode]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &models.PsTsProfessionalRow{ICode: icode, FullName: icode, Role: role}, nil
}

// TestValidateSave_ConflictAppliesToBothRoles reproduce el bug "Trabajador social puede estar
// en varias duplas": la trabajadora social ya asignada a otra dupla debe rechazarse igual que
// pasaría con una profesional en Psicología ya asignada.
func TestValidateSave_ConflictAppliesToBothRoles(t *testing.T) {
	svc := &duplaAdminService{
		duplaRepo: &fakeDuplaRepoConflict{
			occupied: map[string]bool{"TS-YA-ASIGNADA": true},
		},
		reassignRepo: &fakeReassignRepoConflict{
			roles: map[string]string{
				"PS-1":           "ps",
				"TS-YA-ASIGNADA": "ts",
			},
		},
	}

	err := svc.validateSave(context.Background(), "", DuplaAdminSaveInput{
		Name:           "Dupla X",
		PsychologistID: "PS-1",
		SocialWorkerID: "TS-YA-ASIGNADA",
	})

	if err != ErrDuplaMembersConflict {
		t.Fatalf("esperaba ErrDuplaMembersConflict al asignar una trabajadora social ya ocupada, obtuve: %v", err)
	}
}
