package models

import "time"

// ReportListItem representa un reporte en la respuesta de GET /api/v1/victim-contacts/list.
// Ver DocsMD/Screens/Listado Reportes/listado-reportes-interface.md.
type ReportListItem struct {
	ID                 int64     `json:"id"`
	ICode              string    `json:"i_code"`
	Names              string    `json:"names"`
	LastNames          string    `json:"lastNames"`
	CreationDate       time.Time `json:"creationDate"`
	ReportType         string    `json:"reportType"` // Víctima | Tercero | Familiar | Hombre | ""
	TotalFallidos      int       `json:"totalFallidos"`
	DiasCompletos      int       `json:"diasCompletos"`
	ReportStatus       string    `json:"reportStatus"` // sin_asignar | asignado | cerrado_valido | cerrado_invalido
	AssignedAgentICode string    `json:"assignedAgentICode"`
	AssignedAgentName  string    `json:"assignedAgentName"`
	Team               string    `json:"team"`
}

// ReportAttemptItem representa un intento de contacto en el detalle de un reporte.
type ReportAttemptItem struct {
	Date      time.Time `json:"date"`
	Result    string    `json:"result"` // fallido | exitoso
	AgentName string    `json:"agentName"`
}

// ReportDetail respuesta de GET /api/v1/victim-contacts/:icode.
type ReportDetail struct {
	ICode                string              `json:"icode"`
	NombreMostrado       string              `json:"nombreMostrado"`
	ReportType           string              `json:"reportType"`
	ReporterNames        string              `json:"reporterNames"`
	ReporterLastNames    string              `json:"reporterLastNames"`
	Parentesco           string              `json:"parentesco"`
	Telefono             string              `json:"telefono"`
	ReporterEmail        string              `json:"reporterEmail"`
	AjusteRazonable      string              `json:"ajusteRazonable"`
	DescripcionHechos    string              `json:"descripcionHechos"`
	Attempts             []ReportAttemptItem `json:"intentos"`
	IntentosHoy          int                 `json:"intentosHoy"`
	DiasCompletos        int                 `json:"diasCompletos"`
	CloseReasonCode      string              `json:"closeReasonCode"`
	CloseReasonLabel     string              `json:"closeReasonLabel"`
	CloseCategory        string              `json:"closeCategory"`
	CloseObservation     string              `json:"closeObservation"`
	CanManage            bool                `json:"canManage"`
}
