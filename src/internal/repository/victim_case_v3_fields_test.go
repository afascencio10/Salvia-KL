package repository

import "testing"

// REG-12 / REG-34 — columnas del Tamizaje que SaveTamizajeBranch escribe o reinicia.

func TestTamizajeGeneralesYResto(t *testing.T) {
	if got := len(victimCaseV3TamizajeGenerales); got != 8 {
		t.Errorf("generales del Tamizaje: %d columnas, se esperaban 8 (TZ_1..TZ_8)", got)
	}
	// resto = TZ_5, TZ_6, TZ_8 + 24 de la rama pareja + 24 de la rama no pareja
	if got := len(victimCaseV3TamizajeResto); got != 3+24+24 {
		t.Errorf("resto del Tamizaje: %d columnas, se esperaban %d", got, 3+24+24)
	}
	want := map[string]bool{
		"victim_case_form2_tzv3_g5": true, "victim_case_form2_tzv3_g6": true, "victim_case_form2_tzv3_vpp": true,
		"victim_case_form2_tzv3_p1": true, "victim_case_form2_tzv3_p24": true,
		"victim_case_form2_tzv3_np1": true, "victim_case_form2_tzv3_np24": true,
	}
	seen := map[string]bool{}
	for _, c := range victimCaseV3TamizajeResto {
		seen[c] = true
	}
	for c := range want {
		if !seen[c] {
			t.Errorf("falta %s en las columnas del resto del Tamizaje", c)
		}
	}
	// Las cinco iniciales NO son parte del resto.
	for _, c := range []string{"g1", "g2", "g3", "g4", "g7"} {
		if seen["victim_case_form2_tzv3_"+c] {
			t.Errorf("la inicial %s no debe reiniciarse como parte del resto", c)
		}
	}
}

func TestTamizajeTieneResto(t *testing.T) {
	cases := []struct {
		name string
		in   map[string]string
		want bool
	}{
		{"solo las cinco iniciales", map[string]string{"TZ_1": "false", "TZ_2": "false", "TZ_3": "false", "TZ_4": "false", "TZ_7": "false"}, false},
		{"vacío", map[string]string{}, false},
		{"con TZ_5", map[string]string{"TZ_1": "true", "TZ_5": "false"}, true},
		{"con TZ_8", map[string]string{"TZ_8": "true"}, true},
		{"con la rama no pareja", map[string]string{"TZNP_1": "false"}, true},
		{"con la rama pareja", map[string]string{"TZP_24": "true"}, true},
	}
	for _, c := range cases {
		if got := tamizajeTieneResto(c.in); got != c.want {
			t.Errorf("%s: tamizajeTieneResto = %v, se esperaba %v", c.name, got, c.want)
		}
	}
}
