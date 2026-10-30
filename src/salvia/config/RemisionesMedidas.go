package salvia_config

const (
	VistaMedidasSupervisor          = "supervisor"
	VistaMedidasEquipoTecnico       = "equipo-tecnico"
	VistaMedidasSupervisionConvenio = "supervision-convenio"
	VistaMedidasCruzRoja            = "cruz-roja"
	VistaMedidasAgente              = "agente"
	VistaMedidasRemitente           = "remitente"
)

// VistaRemisionesMedidas traduce el rol de sesión a la vista de la bandeja.
// sm y sv reasignan. am es la bandeja propia. ro ve lo que él envió.
// ok es false si el rol no participa.
func VistaRemisionesMedidas(role string) (vista string, puedeReasignar bool, ok bool) {
	switch role {
	case "em":
		return VistaMedidasEquipoTecnico, false, true
	case "sc":
		return VistaMedidasSupervisionConvenio, false, true
	case "cr":
		return VistaMedidasCruzRoja, false, true
	case "sm", "sv":
		return VistaMedidasSupervisor, true, true
	case "am":
		return VistaMedidasAgente, false, true
	case "ro":
		return VistaMedidasRemitente, false, true
	default:
		return "", false, false
	}
}
