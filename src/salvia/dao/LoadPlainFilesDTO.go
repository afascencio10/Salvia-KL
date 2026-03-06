package salvia_daos

var (
	LoadPlainFilesEntityName string = "LoadPlainFiles"
	LoadPlainFilesJSONName   string = "files"
)

type LoadPlainFiles_CustomServiceDTO struct {
	LoadPlainFilesDocType          string `json:"docType"`
	LoadPlainFilesDocNumber        string `json:"docNumber"`
	LoadPlainFilesEntityBranchCode string `json:"branchCode"`
	LoadPlainFilesVictimCaseCode   string `json:"caseCode"`
	LoadPlainFilesMomentCode       string `json:"momentCode"`
}
