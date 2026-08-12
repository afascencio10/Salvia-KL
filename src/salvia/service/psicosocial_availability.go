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

// evaluatePsicosocialAvailability filtra contactos del día y determina si el slot
// [timeStr, timeStr+2h) está libre para el profesional o la dupla.
//
// mode: "individual" | "dupla" (si vacío, se infiere de ps.DuplaID).
// psychologistID / socialWorkerID solo aplican en modo dupla.
func evaluatePsicosocialAvailability(
	contacts []models.TeamContact,
	ps *models.PsychosocialSupport,
	timeStr, mode string,
	psychologistID, socialWorkerID string,
) (available bool, message string) {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		if ps.DuplaID != nil && *ps.DuplaID != "" {
			mode = "dupla"
		} else {
			mode = "individual"
		}
	}

	reqStart, ok := parseTimeToMinutes(timeStr)
	if !ok {
		return false, "Hora inválida; use formato HH:MM"
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
		switch mode {
		case "dupla":
			if ps.DuplaID != nil && tc.DuplaID != nil && *tc.DuplaID == *ps.DuplaID {
				relevant = true
			}
			if tc.ProfessionalID != nil {
				pid := *tc.ProfessionalID
				if (psychologistID != "" && pid == psychologistID) || (socialWorkerID != "" && pid == socialWorkerID) {
					relevant = true
				}
			}
		default: // individual
			if ps.ProfessionalID != nil && tc.ProfessionalID != nil && *tc.ProfessionalID == *ps.ProfessionalID {
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
