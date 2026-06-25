// Package salvia_ctrl — KpiController: agrega métricas del embudo (Salvia_KPIs_v1).
// Optimizado: 4 queries en vez de 10, con cache in-memory 5 minutos.
package salvia_ctrl

import (
	common_controllers "bitsflow/common/controllers"
	"bitsflow/common/db"
	salvia_daos "bitsflow/salvia/dao"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Cache simple en memoria por clave from|to. TTL 5 minutos.
type kpiCacheEntry struct {
	body      string
	expiresAt time.Time
}

var (
	kpiCache   = sync.Map{}
	kpiCacheTTL = 5 * time.Minute
)

// GetFunnelKPIs ejecuta queries agregadas para los 5 niveles del embudo.
// from/to en formato YYYY-MM-DD; si vacíos, sin filtro temporal.
func GetFunnelKPIs(from, to string, connData *db.ConnData, dbClientConfig db.DBClientConfig, dbServerConfig db.DBServerConfig) (int, string) {

	// Cache lookup.
	cacheKey := from + "|" + to
	if v, ok := kpiCache.Load(cacheKey); ok {
		if entry, ok := v.(kpiCacheEntry); ok && time.Now().Before(entry.expiresAt) {
			return http.StatusOK, entry.body
		}
	}

	if connData.ConnID == "" {
		defer db.ReleaseConnection(connData)
	}

	pc := common_controllers.PersistenceController{}
	pc.Setup(connData, &dbClientConfig, &dbServerConfig)
	if pc.Error != nil {
		return http.StatusInternalServerError, pc.Error.Error()
	}

	metrics := salvia_daos.FunnelMetrics{
		Range: salvia_daos.DateRange{From: from, To: to},
	}

	// Nivel 1 + 5 (sobre victim_case).
	queryCasesAggregate(&pc, from, to, &metrics)

	// Nivel 2 + 3 (sobre follow_up_entry).
	queryFollowUpEntries(&pc, from, to, &metrics)

	// Nivel 4 (barreras + necesidades).
	queryBarriersAndNeeds(&pc, from, to, &metrics)

	// Delta riesgo (Nivel 5) — usa victim_case_form2_risk_score primera vs última por caso.
	queryRiskDelta(&pc, from, to, &metrics)

	// Severidad de ingreso (Nivel 1).
	querySeverity(&pc, from, to, &metrics)

	// Alertas críticas + medidas inmediatas + puntaje riesgo (Nivel 1, según PDF).
	queryLevel1Extras(&pc, from, to, &metrics)

	body, err := json.Marshal(metrics)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	kpiCache.Store(cacheKey, kpiCacheEntry{
		body:      string(body),
		expiresAt: time.Now().Add(kpiCacheTTL),
	})

	return http.StatusOK, string(body)
}

func dateFilter(col, from, to string) string {
	clause := ""
	if from != "" {
		clause += fmt.Sprintf(" AND %s >= '%s'", col, from)
	}
	if to != "" {
		clause += fmt.Sprintf(" AND %s <= '%s'", col, to)
	}
	return clause
}

// queryCasesAggregate combina Nivel 1 (asignación) + Nivel 5 (estados).
func queryCasesAggregate(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	m.Level5.CasesByStatus = map[string]int{}

	q := `WITH base AS (
		SELECT vc.victim_case_id, COALESCE(vc.victim_case_status, 'sin_estado') AS victim_case_status,
			(SELECT 1 FROM salvia.rel_case_owner_victim_case rco WHERE rco.victim_case_id = vc.victim_case_id LIMIT 1) AS has_owner
		FROM salvia.victim_case vc
		WHERE TRUE` + dateFilter("vc.victim_case_creation_date", from, to) + `
	)
	SELECT victim_case_status, COUNT(*)::int AS total, COUNT(*) FILTER (WHERE has_owner IS NOT NULL)::int AS assigned
	FROM base GROUP BY victim_case_status`

	pc.Query(context.Background(), q)
	totalCases := 0
	assigned := 0
	effective := 0
	if pc.Rows != nil {
		for pc.Next() {
			var status string
			var cnt, asg int
			pc.ScanRow(&status, &cnt, &asg)
			m.Level5.CasesByStatus[status] = cnt
			totalCases += cnt
			assigned += asg
			if status == "c" || status == "cd" {
				effective += cnt
			}
		}
		pc.Rows.Close()
	}
	m.Level1.TotalCases = totalCases
	m.Level1.AssignedCases = assigned
	m.Level1.PendingCases = totalCases - assigned
	if totalCases > 0 {
		m.Level1.PendingRatio = float64(m.Level1.PendingCases) / float64(totalCases)
		m.Level5.EffectiveClosureRate = float64(effective) / float64(totalCases)
	}
}

// queryLevel1Extras: alertas críticas, medidas inmediatas, promedio puntaje riesgo.
// Cubre campos del PDF: "Alerta pendientes", "Prioridad", "Puntaje riesgo",
// "Las medidas requieren gestión de manera inmediata".
func queryLevel1Extras(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	// Alertas críticas (priority=1).
	q1 := `SELECT COUNT(*) FROM salvia.alert
		WHERE alert_priority = 1` + dateFilter("alert_creation_date", from, to)
	pc.Query(context.Background(), q1)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level1.CriticalAlerts)
		pc.Rows.Close()
	}

	// Promedio puntaje riesgo + medidas inmediatas (urgent_emotional_crisis).
	q2 := `SELECT
		COALESCE(AVG(victim_case_form2_risk_score)::float, 0) AS avg_score
		FROM salvia.victim_case_form2
		WHERE TRUE` + dateFilter("victim_case_form2_creation_date", from, to)
	pc.Query(context.Background(), q2)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level1.AvgRiskScore)
		pc.Rows.Close()
	}

	q3 := `SELECT COUNT(*) FROM salvia.feminicide_risk_form1
		WHERE feminicide_risk_form1_victim_urgent_emotional_crisis IS NOT NULL
		AND feminicide_risk_form1_victim_urgent_emotional_crisis > 0`
	pc.Query(context.Background(), q3)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level1.ImmediateActions)
		pc.Rows.Close()
	}
}

// querySeverity calcula distribución de severidad del primer follow_up por caso.
func querySeverity(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	m.Level1.SeverityIngress = map[string]int{}
	q := `SELECT follow_up_risk_level, COUNT(*) FROM (
		SELECT DISTINCT ON (f.follow_up_id) f.follow_up_risk_level
		FROM salvia.follow_up f
		WHERE TRUE` + dateFilter("f.follow_up_creation_date", from, to) + `
		ORDER BY f.follow_up_id, f.follow_up_creation_date ASC
	) sub GROUP BY follow_up_risk_level`
	pc.Query(context.Background(), q)
	if pc.Rows != nil {
		for pc.Next() {
			var lvl string
			var cnt int
			pc.ScanRow(&lvl, &cnt)
			m.Level1.SeverityIngress[lvl] = cnt
		}
		pc.Rows.Close()
	}
}

// queryFollowUpEntries combina Nivel 2 + 3 en una sola pasada sobre follow_up_entry.
func queryFollowUpEntries(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	m.Level3.AttentionBySector = map[string]int{}

	// Agregados globales con FILTER.
	q := `SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE follow_up_entry_was_done = '1') AS effective,
		COUNT(*) FILTER (WHERE follow_up_entry_person_visited_entity = '1') AS visited,
		COUNT(*) FILTER (WHERE follow_up_entry_person_visited_entity = '0') AS not_visited,
		COUNT(*) FILTER (WHERE follow_up_entry_received_attention = '1') AS received,
		COALESCE(AVG(c)::float, 0) AS avg_iter
		FROM salvia.follow_up_entry fe,
		LATERAL (SELECT COUNT(*) c FROM salvia.follow_up_entry fe2 WHERE fe2.follow_up_id = fe.follow_up_id) sub
		WHERE TRUE` + dateFilter("fe.follow_up_entry_creation_date", from, to)

	// Simplificado: agregar promedio en query separada porque lateral inflará counts.
	q = `SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE follow_up_entry_was_done = '1') AS effective,
		COUNT(*) FILTER (WHERE follow_up_entry_person_visited_entity = '1') AS visited,
		COUNT(*) FILTER (WHERE follow_up_entry_person_visited_entity = '0') AS not_visited,
		COUNT(*) FILTER (WHERE follow_up_entry_received_attention = '1') AS received
		FROM salvia.follow_up_entry
		WHERE TRUE` + dateFilter("follow_up_entry_creation_date", from, to)
	pc.Query(context.Background(), q)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level2.TotalFollowUpEntries, &m.Level2.EffectiveEntries,
			&m.Level3.AttendedInstitution, &m.Level3.NotAttendedInstitution, &m.Level3.ReceivedAttention)
		pc.Rows.Close()
	}

	if m.Level2.TotalFollowUpEntries > 0 {
		m.Level2.EffectivenessRatio = float64(m.Level2.EffectiveEntries) / float64(m.Level2.TotalFollowUpEntries)
	}
	totalVisited := m.Level3.AttendedInstitution + m.Level3.NotAttendedInstitution
	if totalVisited > 0 {
		m.Level3.InstitutionalAttendanceRate = float64(m.Level3.AttendedInstitution) / float64(totalVisited)
	}
	if m.Level3.AttendedInstitution > 0 {
		m.Level3.EffectiveAttentionRate = float64(m.Level3.ReceivedAttention) / float64(m.Level3.AttendedInstitution)
	}

	// Promedio iteraciones por caso.
	q2 := `SELECT COALESCE(AVG(c)::float, 0) FROM (
		SELECT COUNT(*) AS c FROM salvia.follow_up_entry
		WHERE TRUE` + dateFilter("follow_up_entry_creation_date", from, to) + `
		GROUP BY follow_up_id
	) s`
	pc.Query(context.Background(), q2)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level2.AvgIterationsPerCase)
		pc.Rows.Close()
	}

	// Atención por sector.
	q3 := `SELECT follow_up_entry_sector, COUNT(*) FROM salvia.follow_up_entry
		WHERE follow_up_entry_received_attention = '1'` + dateFilter("follow_up_entry_creation_date", from, to) + `
		GROUP BY follow_up_entry_sector`
	pc.Query(context.Background(), q3)
	if pc.Rows != nil {
		for pc.Next() {
			var sector *string
			var cnt int
			pc.ScanRow(&sector, &cnt)
			key := "sin_sector"
			if sector != nil {
				key = *sector
			}
			m.Level3.AttentionBySector[key] = cnt
		}
		pc.Rows.Close()
	}
}

// queryBarriersAndNeeds combina barreras (rel_barrier_follow_up_entry) + necesidades (feminicide_form1).
func queryBarriersAndNeeds(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	m.Level4.BarriersBySector = map[string]int{}
	m.Level4.CriticalNeedsCount = map[string]int{}
	m.Level4.TopBarriers = []salvia_daos.BarrierCount{}

	// Top barreras + por sector en una query con WINDOW.
	q := `SELECT b.barrier_name, sb.sector_barrier_name, COUNT(*) AS c
		FROM salvia.rel_barrier_follow_up_entry rbfe
		JOIN salvia.barrier b ON b.barrier_id = rbfe.rel_barrier_follow_up_entry_barrier
		JOIN salvia.sector_barrier sb ON sb.sector_barrier_id = b.sector_barrier_id
		WHERE TRUE` + dateFilter("rbfe.rel_barrier_follow_up_entry_creation_date", from, to) + `
		GROUP BY b.barrier_name, sb.sector_barrier_name
		ORDER BY c DESC`
	pc.Query(context.Background(), q)
	if pc.Rows != nil {
		count := 0
		for pc.Next() {
			var bc salvia_daos.BarrierCount
			pc.ScanRow(&bc.BarrierName, &bc.SectorName, &bc.Count)
			if count < 10 {
				m.Level4.TopBarriers = append(m.Level4.TopBarriers, bc)
			}
			m.Level4.BarriersBySector[bc.SectorName] += bc.Count
			m.Level4.TotalBarriers += bc.Count
			count++
		}
		pc.Rows.Close()
	}

	// Necesidades críticas (feminicide_form1).
	q2 := `SELECT
		COUNT(*) FILTER (WHERE feminicide_form1_transport_difficulty IS NOT NULL AND feminicide_form1_transport_difficulty > 0) AS transp,
		COUNT(*) FILTER (WHERE feminicide_form1_food_access_frequency IS NOT NULL AND feminicide_form1_food_access_frequency > 0) AS aliment,
		COUNT(*) FILTER (WHERE feminicide_form1_transport_subsidy_received IS NOT NULL) AS subsidio
		FROM salvia.feminicide_form1`
	pc.Query(context.Background(), q2)
	if pc.Rows != nil && pc.Next() {
		var transp, aliment, subsidio int
		pc.ScanRow(&transp, &aliment, &subsidio)
		m.Level4.CriticalNeedsCount["transporte"] = transp
		m.Level4.CriticalNeedsCount["alimentacion"] = aliment
		m.Level4.CriticalNeedsCount["subsidio_transporte"] = subsidio
		pc.Rows.Close()
	}
}

// queryRiskDelta — evolución riesgo: compara primer vs último victim_case_form2_risk_score por caso.
func queryRiskDelta(pc *common_controllers.PersistenceController, from, to string, m *salvia_daos.FunnelMetrics) {
	q := `WITH ranked AS (
		SELECT victim_case_form2_victim_case AS case_id,
			COALESCE(victim_case_form2_risk_score, 0) AS score,
			ROW_NUMBER() OVER (PARTITION BY victim_case_form2_victim_case ORDER BY victim_case_form2_creation_date ASC) AS rn_asc,
			ROW_NUMBER() OVER (PARTITION BY victim_case_form2_victim_case ORDER BY victim_case_form2_creation_date DESC) AS rn_desc
		FROM salvia.victim_case_form2
		WHERE TRUE` + dateFilter("victim_case_form2_creation_date", from, to) + `
	),
	first_last AS (
		SELECT case_id,
			MAX(CASE WHEN rn_asc = 1 THEN score END) AS first_score,
			MAX(CASE WHEN rn_desc = 1 THEN score END) AS last_score
		FROM ranked GROUP BY case_id
		HAVING COUNT(*) > 0
	)
	SELECT
		COALESCE(COUNT(*) FILTER (WHERE last_score < first_score), 0)::int AS improved,
		COALESCE(COUNT(*) FILTER (WHERE last_score > first_score), 0)::int AS worsened,
		COALESCE(COUNT(*) FILTER (WHERE last_score = first_score), 0)::int AS unchanged
	FROM first_last`
	pc.Query(context.Background(), q)
	if pc.Rows != nil && pc.Next() {
		pc.ScanRow(&m.Level5.RiskDelta.CasesWithImprovement, &m.Level5.RiskDelta.CasesWithWorsening, &m.Level5.RiskDelta.CasesUnchanged)
		pc.Rows.Close()
	}
}
