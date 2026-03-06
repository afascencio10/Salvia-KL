package salvia_daos

import (
	"bitsflow/common/utils"
	"time"
)

var (
	VictimCaseReportEntityName string = "VictimCaseReport"
	VictimCaseReportJSONName   string = "report"

	//Atributos relacionados con las validaciones ------------------------------

	//Campos que vienen como string del JSON. El booleano indica si son strings en el modelo o no (como en el caso de una fecha)
	VictimCaseReportFieldDefinitions map[string]utils.FieldDefinition = map[string]utils.FieldDefinition{
		"VictimCaseReportStartDate":    {Name: "VictimCaseReportStartDate", DBName: "", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseReportEndDate":      {Name: "VictimCaseReportEndDate", DBName: "", Alias: "", ModelType: "date", MinSize: 0, MaxSize: 0, Required: true},
		"VictimCaseReportDepartment":   {Name: "VictimCaseReportDepartment", DBName: "", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 36, Required: false},
		"VictimCaseReportCity":         {Name: "VictimCaseReportCity", DBName: "", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 36, Required: false},
		"VictimCaseReportTown":         {Name: "VictimCaseReportTown", DBName: "", Alias: "", ModelType: "string", MinSize: 2, MaxSize: 36, Required: false},
		"VictimCaseReportViolenceType": {Name: "VictimCaseReportViolenceType", DBName: "", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: false},
		"VictimCaseReportCaseStatus":   {Name: "VictimCaseReportCaseStatus", DBName: "", Alias: "", ModelType: "string", MinSize: 1, MaxSize: 2, Required: false},
	}
)

type VictimCaseReportDTO struct {
	VictimCaseReportStartDate    time.Time `json:"startDate"`
	VictimCaseReportEndDate      time.Time `json:"endDate"`
	VictimCaseReportDepartment   string    `json:"department"`
	VictimCaseReportCity         string    `json:"city"`
	VictimCaseReportTown         string    `json:"town"`
	VictimCaseReportViolenceType string    `json:"type"`
	VictimCaseReportCaseStatus   string    `json:"status"`
}
