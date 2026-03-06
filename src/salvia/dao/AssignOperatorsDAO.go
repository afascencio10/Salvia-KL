package salvia_daos

var (
	AssignOperatorsEntityName string = "AssignOperators"
	AssignOperatorsJSONName   string = "files"
)

type AssignOperators_CustomServiceDTO struct {
	AssignOperatorsDocType          string `json:"docType"`
	AssignOperatorsDocNumber        string `json:"docNumber"`
	AssignOperatorsEntityBranchCode string `json:"branchCode"`
	AssignOperatorsVictimCaseCode   string `json:"caseCode"`
	AssignOperatorsMomentCode       string `json:"momentCode"`
}
