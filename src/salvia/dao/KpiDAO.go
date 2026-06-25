// Package salvia_daos — KpiDAO: DTOs para métricas del embudo analítico (Salvia_KPIs_v1).
// No es entidad CRUD; solo estructuras de respuesta para el dashboard.
package salvia_daos

// FunnelMetrics agrupa los 5 niveles del embudo.
type FunnelMetrics struct {
	Level1 Level1Metrics `json:"level1"`
	Level2 Level2Metrics `json:"level2"`
	Level3 Level3Metrics `json:"level3"`
	Level4 Level4Metrics `json:"level4"`
	Level5 Level5Metrics `json:"level5"`
	Range  DateRange     `json:"range"`
}

type DateRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Level1 — Entrada y Asignación.
type Level1Metrics struct {
	TotalCases       int            `json:"totalCases"`
	AssignedCases    int            `json:"assignedCases"`
	PendingCases     int            `json:"pendingCases"`
	PendingRatio     float64        `json:"pendingRatio"`
	SeverityIngress  map[string]int `json:"severityIngress"` // h/m/l
	CriticalAlerts   int            `json:"criticalAlerts"`   // alert.alert_priority = 1
	ImmediateActions int            `json:"immediateActions"` // urgent_emotional_crisis = sí
	AvgRiskScore     float64        `json:"avgRiskScore"`     // promedio victim_case_form2_risk_score
}

// Level2 — Contacto y Gestión Operativa.
type Level2Metrics struct {
	TotalFollowUpEntries int     `json:"totalFollowUpEntries"`
	EffectiveEntries     int     `json:"effectiveEntries"`
	EffectivenessRatio   float64 `json:"effectivenessRatio"`
	AvgIterationsPerCase float64 `json:"avgIterationsPerCase"`
}

// Level3 — Conversión y Ruta Institucional.
type Level3Metrics struct {
	AttendedInstitution         int            `json:"attendedInstitution"`
	NotAttendedInstitution      int            `json:"notAttendedInstitution"`
	InstitutionalAttendanceRate float64        `json:"institutionalAttendanceRate"`
	ReceivedAttention           int            `json:"receivedAttention"`
	EffectiveAttentionRate      float64        `json:"effectiveAttentionRate"`
	AttentionBySector           map[string]int `json:"attentionBySector"`
}

// Level4 — Fricción y Bloqueos.
type Level4Metrics struct {
	TotalBarriers      int                       `json:"totalBarriers"`
	BarriersBySector   map[string]int            `json:"barriersBySector"`
	TopBarriers        []BarrierCount            `json:"topBarriers"`
	CriticalNeedsCount map[string]int            `json:"criticalNeedsCount"` // transporte, alimentación
}

type BarrierCount struct {
	BarrierName string `json:"name"`
	SectorName  string `json:"sector"`
	Count       int    `json:"count"`
}

// Level5 — Mitigación de Riesgo y Resolución.
type Level5Metrics struct {
	CasesByStatus       map[string]int `json:"casesByStatus"`     // fc/r/ra/c/ex/is/cd/iv
	EffectiveClosureRate float64       `json:"effectiveClosureRate"` // (c+cd) / total
	RiskDelta           RiskDeltaInfo  `json:"riskDelta"`
}

type RiskDeltaInfo struct {
	CasesWithImprovement int `json:"casesWithImprovement"`
	CasesWithWorsening   int `json:"casesWithWorsening"`
	CasesUnchanged       int `json:"casesUnchanged"`
}
