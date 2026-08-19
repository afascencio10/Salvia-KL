package service

import (
	"bitsflow/internal/models"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const psicosocialSessionDurationMin = 120 // ventana de 2 horas

// parseTimeToMinutes convierte "HH:MM" o "HH:MM:SS" a minutos desde medianoche.
func parseTimeToMinutes(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	parts := strings.Split(s, ":")
	if len(parts) < 2 {
		return 0, false
	}
	h, errH := strconv.Atoi(parts[0])
	m, errM := strconv.Atoi(parts[1])
	if errH != nil || errM != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, false
	}
	return h*60 + m, true
}

// sessionWindowsOverlap indica si [a, a+dur) se solapa con [b, b+dur).
func sessionWindowsOverlap(startA, startB, durationMin int) bool {
	endA := startA + durationMin
	endB := startB + durationMin
	return startA < endB && startB < endA
}

// seedAvailabilityActors arma personIDs y duplaIDs base desde la remisión:
//   - personIDs: professional_id y/o miembros (PS+TS) de ps.dupla_id
//   - duplaIDs:  ps.dupla_id si existe
// remisionPsychID / remisionSWID vienen de cargar la dupla de la remisión (pueden ir vacíos).
func seedAvailabilityActors(ps *models.PsychosocialSupport, remisionPsychID, remisionSWID string) (personIDs, duplaIDs []string) {
	seenP := map[string]struct{}{}
	addP := func(id string) {
		id = strings.TrimSpace(id)
		if id == "" {
			return
		}
		if _, ok := seenP[id]; ok {
			return
		}
		seenP[id] = struct{}{}
		personIDs = append(personIDs, id)
	}
	if ps != nil {
		if ps.ProfessionalID != nil {
			addP(*ps.ProfessionalID)
		}
		if ps.DuplaID != nil && strings.TrimSpace(*ps.DuplaID) != "" {
			duplaIDs = append(duplaIDs, strings.TrimSpace(*ps.DuplaID))
			addP(remisionPsychID)
			addP(remisionSWID)
		}
	}
	return personIDs, duplaIDs
}

// mergeDuplaIDs une seed + expansión (duplas activas del profesional) sin duplicados.
func mergeDuplaIDs(seed, expanded []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(seed)+len(expanded))
	for _, id := range append(append([]string{}, seed...), expanded...) {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// evaluatePsicosocialAvailabilityWithScope determina si [timeStr, timeStr+2h) está libre
// frente a contactos del día, considerando:
//   - sesiones con professional_id ∈ personIDs
//   - sesiones con dupla_id ∈ duplaIDs
// (duplaIDs debe incluir la dupla de la remisión y las duplas activas donde figuren
// los personIDs — expansión hecha por el caller).
func evaluatePsicosocialAvailabilityWithScope(
	contacts []models.TeamContact,
	timeStr string,
	personIDs, duplaIDs []string,
) (available bool, message string) {
	reqStart, ok := parseTimeToMinutes(timeStr)
	if !ok {
		return false, "Hora inválida; use formato HH:MM"
	}

	personSet := map[string]struct{}{}
	for _, id := range personIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			personSet[id] = struct{}{}
		}
	}
	duplaSet := map[string]struct{}{}
	for _, id := range duplaIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			duplaSet[id] = struct{}{}
		}
	}

	for _, tc := range contacts {
		if tc.ScheduledTime == nil || strings.TrimSpace(*tc.ScheduledTime) == "" {
			continue
		}
		existStart, ok := parseTimeToMinutes(*tc.ScheduledTime)
		if !ok {
			continue
		}
		if !sessionWindowsOverlap(reqStart, existStart, psicosocialSessionDurationMin) {
			continue
		}

		relevant := false
		if tc.ProfessionalID != nil {
			if _, ok := personSet[strings.TrimSpace(*tc.ProfessionalID)]; ok {
				relevant = true
			}
		}
		if !relevant && tc.DuplaID != nil {
			if _, ok := duplaSet[strings.TrimSpace(*tc.DuplaID)]; ok {
				relevant = true
			}
		}
		if relevant {
			existEnd := existStart + psicosocialSessionDurationMin
			msg := fmt.Sprintf(
				"Horario no disponible: hay una sesión de %02d:%02d a %02d:%02d (ventana de 2 horas)",
				existStart/60, existStart%60, existEnd/60, existEnd%60,
			)
			return false, msg
		}
	}

	return true, "Horario disponible"
}

// normalizeScheduledTime deja la hora en formato HH:MM (máx 8 chars en columna).
func normalizeScheduledTime(hora string) string {
	hora = strings.TrimSpace(hora)
	mins, ok := parseTimeToMinutes(hora)
	if !ok {
		return hora
	}
	return fmt.Sprintf("%02d:%02d", mins/60, mins%60)
}

// parseAgendaDate parsea YYYY-MM-DD.
func parseAgendaDate(fechaStr string) (time.Time, error) {
	return time.Parse("2006-01-02", strings.TrimSpace(fechaStr))
}
