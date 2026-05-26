// Package models — case_full_info.go
// Modelo de respuesta para el componente <case-info> que muestra toda la información del caso.
package models

// CaseFullInfo agrupa toda la información disponible de un caso
// combinando victim_case + form1 + form2 con prioridad form2.
type CaseFullInfo struct {
	Encabezado      CaseInfoEncabezado      `json:"encabezado"`
	Victima         CaseInfoVictima         `json:"victima"`
	DatosPersonales CaseInfoDatosPersonales `json:"datosPersonales"`
	Etnicos         CaseInfoEtnicos         `json:"etnicos"`
	Contacto        CaseInfoContacto        `json:"contacto"`
	Ubicacion       CaseInfoUbicacion       `json:"ubicacion"`
	Hechos          CaseInfoHechos          `json:"hechos"`
	Agresor         CaseInfoAgresor         `json:"agresor"`
	Riesgo          CaseInfoRiesgo          `json:"riesgo"`
}

type CaseInfoEncabezado struct {
	Funcionarios      string `json:"funcionarios"`
	FechaCreacion     string `json:"fechaCreacion"`
	FechaModificacion string `json:"fechaModificacion"`
	Estado            string `json:"estado"`
}

type CaseInfoVictima struct {
	Nombres           string `json:"nombres"`
	Apellidos         string `json:"apellidos"`
	NombreIdentitario string `json:"nombreIdentitario"`
	Edad              *int64 `json:"edad"`
	FechaNacimiento   string `json:"fechaNacimiento"`
	TipoDocumento     string `json:"tipoDocumento"`
	NumeroDocumento   string `json:"numeroDocumento"`
	Nacionalidad      string `json:"nacionalidad"`
	OtraNacionalidad  string `json:"otraNacionalidad"`
	Municipio         string `json:"municipio"`
	CorreoElectronico string `json:"correoElectronico"`
	DireccionResidencia string `json:"direccionResidencia"`
	AjusteRazonable   string `json:"ajusteRazonable"`
}

type CaseInfoDatosPersonales struct {
	CondicionMigratoria    string `json:"condicionMigratoria"`
	Direccion              string `json:"direccion"`
	Telefono               string `json:"telefono"`
	Genero                 string `json:"genero"`
	IdentidadGenero        string `json:"identidadGenero"`
	OtraIdentidadGenero    string `json:"otraIdentidadGenero"`
	OrientacionSexual      string `json:"orientacionSexual"`
	OtraOrientacionSexual  string `json:"otraOrientacionSexual"`
	Procedencia            string `json:"procedencia"`
	Ocupacion              string `json:"ocupacion"`
	OtraOcupacion          string `json:"otraOcupacion"`
}

type CaseInfoEtnicos struct {
	GrupoEtnico    string `json:"grupoEtnico"`
	OtroGrupoEtnico string `json:"otroGrupoEtnico"`
	Afrodescendiente string `json:"afrodescendiente"`
	Indigena        string `json:"indigena"`
	LenguaIndigena  string `json:"lenguaIndigena"`
	Campesino       string `json:"campesino"`
	VictimaConflicto string `json:"victimaConflicto"`
}

type CaseInfoContacto struct {
	NombreContacto   string `json:"nombreContacto"`
	TelefonoContacto string `json:"telefonoContacto"`
	Parentesco       string `json:"parentesco"`
	PersonasCargo    string `json:"personasCargo"`
	NumeroHijos      *int   `json:"numeroHijos"`
	EdadesHijos      string `json:"edadesHijos"`
	EstadoCivil      string `json:"estadoCivil"`
	OtroEstadoCivil  string `json:"otroEstadoCivil"`
	Discapacidad     string `json:"discapacidad"`
	ApoyoEspecial    string `json:"apoyoEspecial"`
}

type CaseInfoUbicacion struct {
	Departamento string `json:"departamento"`
	Ciudad       string `json:"ciudad"`
	Municipio    string `json:"municipio"`
}

type CaseInfoHechos struct {
	Descripcion          string `json:"descripcion"`
	Ocurrencia           string `json:"ocurrencia"`
	Horario              string `json:"horario"`
	DiaSemana            string `json:"diaSemana"`
	FechaHechos          string `json:"fechaHechos"`
	ViolenciaExperimentada string `json:"violenciaExperimentada"`
	OtroTipoViolencia    string `json:"otroTipoViolencia"`
	AmbitoViolencia      string `json:"ambitoViolencia"`
	EscenarioViolencia   string `json:"escenarioViolencia"`
	RiesgoFeminicida     string `json:"riesgoFeminicida"`
	DireccionHechos      string `json:"direccionHechos"`
}

type CaseInfoAgresor struct {
	TipoAgresor       string `json:"tipoAgresor"`
	Relacion          string `json:"relacion"`
	Nombre            string `json:"nombre"`
	TipoDocumento     string `json:"tipoDocumento"`
	NumeroDocumento   string `json:"numeroDocumento"`
	Direccion         string `json:"direccion"`
	Telefono          string `json:"telefono"`
	NumAgresores      string `json:"numAgresores"`
	Proximidad        string `json:"proximidad"`
	GeneroAgresor     string `json:"generoAgresor"`
}

type CaseInfoRiesgo struct {
	NivelRiesgo              *int   `json:"nivelRiesgo"`
	NivelRiesgoTexto         string `json:"nivelRiesgoTexto"`
	AmenazasMuerte           string `json:"amenazasMuerte"`
	AgresorTieneArmas        string `json:"agresorTieneArmas"`
	ViolenciaPrevia          string `json:"violenciaPrevia"`
	ViolenciaFisicaIncremento string `json:"violenciaFisicaIncremento"`
	SeparacionUltimoAnio     string `json:"separacionUltimoAnio"`
	AmenazoConArma           string `json:"amenazoConArma"`
	AmenazoHijos             string `json:"amenazoHijos"`
	CelosoViolento           string `json:"celosoViolento"`
	CreeCapazMatar           string `json:"creeCapazMatar"`
	RiesgoInminente          string `json:"riesgoInminente"`
	DenunciaPrevia           string `json:"denunciaPrevia"`
	SiDenuncioAntes          string `json:"siDenuncioAntes"`
}
