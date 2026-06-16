package utils

import "testing"

func TestNormalizeLocationName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Bogotá D.C.", "BOGOTA D C"},
		{"bogota", "BOGOTA"},
		{"MEDELLÍN", "MEDELLIN"},
		{"  Santa Fé de Antioquia ", "SANTA FE DE ANTIOQUIA"},
		{"BRICEÑO", "BRICENO"},
	}

	for _, tt := range tests {
		got := NormalizeLocationName(tt.in)
		if got != tt.want {
			t.Errorf("NormalizeLocationName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestScoreCityNameMatch_Bogota(t *testing.T) {
	input := NormalizeLocationName("Bogota")
	city := NormalizeLocationName("BOGOTÁ D.C.")
	score := ScoreCityNameMatch(input, city)
	if score < 850 {
		t.Fatalf("expected strong match for Bogota vs BOGOTÁ D.C., got score %d", score)
	}
}

func TestScoreCityNameMatch_Medellin(t *testing.T) {
	input := NormalizeLocationName("Medellin")
	city := NormalizeLocationName("MEDELLÍN")
	if ScoreCityNameMatch(input, city) != 1000 {
		t.Fatalf("expected exact normalized match for Medellin")
	}
}
