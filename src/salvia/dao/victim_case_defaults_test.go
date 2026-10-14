package salvia_daos

import (
	common_dao "bitsflow/common/dao"
	"bitsflow/common/utils"
	"testing"
)

func TestSetVictimCaseDefaults_InsertAlwaysSetsStatus(t *testing.T) {
	roles := []struct {
		role   string
		status string
	}{
		{"op", "ra"},
		{"ro", "ra"},
		{"an", "ra"},
		{"sv", "ra"},
		{"", "ra"},
		{"et", "r"},
	}
	for _, tc := range roles {
		vc := VictimCaseDTO{VictimCaseStatus: "should-be-replaced"}
		SetVictimCaseDefaults(&vc, common_dao.SQL_INSERT, utils.CommonSession{CurrentRole: tc.role})
		if vc.VictimCaseStatus != tc.status {
			t.Errorf("rol %q: status %q, esperaba %q", tc.role, vc.VictimCaseStatus, tc.status)
		}
		if vc.VictimCaseICode == "" {
			t.Errorf("rol %q: no asignó icode", tc.role)
		}
	}
}

func TestSetVictimCaseDefaults_UpdateDoesNotTouchStatus(t *testing.T) {
	vc := VictimCaseDTO{VictimCaseStatus: "cd", VictimCaseICode: "keep"}
	SetVictimCaseDefaults(&vc, common_dao.SQL_UPDATE, utils.CommonSession{CurrentRole: "ro"})
	if vc.VictimCaseStatus != "cd" {
		t.Errorf("update cambió el estado a %q", vc.VictimCaseStatus)
	}
	if vc.VictimCaseICode != "keep" {
		t.Errorf("update cambió el icode")
	}
}
