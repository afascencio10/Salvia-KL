package salvia_daos

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

// T1 (REP-13*): la hora de recontacto del reportante llega como "HH:MM:SS"
// sin zona desde el wizard V2; antes encoding/json exigía RFC 3339.
func TestVictimContactForm2V2DTO_UnmarshalBestContactTime(t *testing.T) {
	cases := []struct {
		name     string
		payload  string
		wantZero bool
		wantHH   int
		wantMM   int
		wantErr  bool
	}{
		{name: "HH:MM:SS del wizard", payload: `{"reporterBestContactTime":"14:30:00"}`, wantHH: 14, wantMM: 30},
		{name: "HH:MM", payload: `{"reporterBestContactTime":"09:05"}`, wantHH: 9, wantMM: 5},
		{name: "medianoche es una hora válida", payload: `{"reporterBestContactTime":"00:00:00"}`, wantHH: 0, wantMM: 0},
		{name: "RFC 3339 sigue funcionando", payload: `{"reporterBestContactTime":"2026-01-02T07:45:00Z"}`, wantHH: 7, wantMM: 45},
		{name: "null", payload: `{"reporterBestContactTime":null}`, wantZero: true},
		{name: "vacío", payload: `{"reporterBestContactTime":""}`, wantZero: true},
		{name: "ausente", payload: `{}`, wantZero: true},
		{name: "fuera de rango", payload: `{"reporterBestContactTime":"25:99:00"}`, wantErr: true},
		{name: "texto no hora", payload: `{"reporterBestContactTime":"mañana"}`, wantErr: true},
		{name: "número", payload: `{"reporterBestContactTime":1430}`, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var dto VictimContactForm2V2DTO
			err := json.Unmarshal([]byte(c.payload), &dto)

			if c.wantErr {
				if !errors.Is(err, ErrReporterBestContactTimeInvalid) {
					t.Fatalf("esperaba ErrReporterBestContactTimeInvalid, obtuve %v", err)
				}
				if !dto.VictimContactForm2V2ReporterBestContactTime.IsZero() {
					t.Fatalf("con hora inválida el campo debe quedar en cero")
				}
				return
			}
			if err != nil {
				t.Fatalf("no esperaba error: %v", err)
			}

			got := dto.VictimContactForm2V2ReporterBestContactTime
			if c.wantZero {
				if !got.IsZero() {
					t.Fatalf("esperaba hora en cero, obtuve %v", got)
				}
				return
			}
			if got.IsZero() {
				t.Fatalf("la hora válida no debe quedar en cero (IsZero = true)")
			}
			if got.Hour() != c.wantHH || got.Minute() != c.wantMM {
				t.Fatalf("esperaba %02d:%02d, obtuve %02d:%02d", c.wantHH, c.wantMM, got.Hour(), got.Minute())
			}
		})
	}
}

// Una hora inválida no debe impedir que el resto del DTO se decodifique: el
// controller necesita los demás campos para seguir validando.
func TestVictimContactForm2V2DTO_InvalidTimeKeepsOtherFields(t *testing.T) {
	var dto VictimContactForm2V2DTO
	err := json.Unmarshal([]byte(`{"reporterLastNames":"Pérez","reporterEmail":"a@b.co","reporterBestContactTime":"nunca"}`), &dto)
	if !errors.Is(err, ErrReporterBestContactTimeInvalid) {
		t.Fatalf("esperaba ErrReporterBestContactTimeInvalid, obtuve %v", err)
	}
	if dto.VictimContactForm2V2ReporterLastNames != "Pérez" || dto.VictimContactForm2V2ReporterEmail != "a@b.co" {
		t.Fatalf("los demás campos deben decodificarse: %+v", dto)
	}
}

// Las fechas que MarshalJSON escribe como "dd/mm/yyyy hh:mm:ss" no son RFC
// 3339: UnmarshalJSON no debe fallar si un cliente las devuelve.
func TestVictimContactForm2V2DTO_RoundTripWithDates(t *testing.T) {
	src := VictimContactForm2V2DTO{
		VictimContactForm2V2CreationDate:            time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC),
		VictimContactForm2V2UpdateDate:              time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC),
		VictimContactForm2V2ReporterLastNames:       "Gómez",
		VictimContactForm2V2ReporterBestContactTime: time.Date(0, 1, 1, 16, 20, 0, 0, time.UTC),
	}
	raw, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var back VictimContactForm2V2DTO
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("unmarshal del JSON que produce MarshalJSON: %v (json=%s)", err, raw)
	}
	if back.VictimContactForm2V2ReporterLastNames != "Gómez" {
		t.Fatalf("apellidos perdidos: %+v", back)
	}
	if back.VictimContactForm2V2ReporterBestContactTime.Hour() != 16 || back.VictimContactForm2V2ReporterBestContactTime.Minute() != 20 {
		t.Fatalf("hora perdida: %v", back.VictimContactForm2V2ReporterBestContactTime)
	}
}
