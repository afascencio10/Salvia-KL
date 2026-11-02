package models

// ReportDetailExt es la respuesta de GET /api/v1/victim-contacts/:icode: el
// ReportDetail de siempre (report_list_item.go, archivo del Listado) más los
// datos que agrega WP-1, para no tocar ese archivo.
type ReportDetailExt struct {
	ReportDetail
	// Radicado «XXXX-XXXX» del reporte V2 (T3); "" en reportes anteriores y en los de V1.
	Radicado string `json:"radicado"`
	// Código del tipo de reporte V2: v | r | fm | ho. "" en V1 (no tiene).
	ReportTypeCode string `json:"reportTypeCode"`
}
