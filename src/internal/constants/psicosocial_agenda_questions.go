// Package constants — IDs fijos de preguntas de agenda condicional (Aug 2026).
//
// Serie a1xxxxxx = ¿Agendar nueva sesión? (single Sí/No)
// Serie a2xxxxxx = Hora próxima atención (time)
//
// Sufijo N en aN00000X:
//
//	1 = Primer Contacto S1
//	2 = Primer Contacto S4
//	3 = Primera Atención S1
//	4 = Primera Atención S4
//	5 = Atención Psicosocial (SEG) S1
//	6 = Atención Psicosocial (SEG) S4
//	7 = Cierre S1
//	8 = Cierre S4
//
// Las fechas existentes ("Fecha nueva" / "Fecha próxima atención") ya tienen UUIDs
// hardcodeados en extractAgendaAnswers / extractFechaProximaAtencion — no se redefinen aquí.
package constants

const (
	// Primer Contacto — S1 (visible si Continuar Primera Atención = No)
	QPCAgendarNuevaSesionS1  = "a1000001-0000-4000-8000-000000000001"
	QPCHoraProximaAtencionS1 = "a2000001-0000-4000-8000-000000000001"
	// Primer Contacto — S4 (visible si Continuar = Sí)
	QPCAgendarNuevaSesionS4  = "a1000002-0000-4000-8000-000000000002"
	QPCHoraProximaAtencionS4 = "a2000002-0000-4000-8000-000000000002"

	// Primera Atención — S1 (Solo Contacto) / S4 (Atención)
	QPAAgendarNuevaSesionS1  = "a1000003-0000-4000-8000-000000000003"
	QPAHoraProximaAtencionS1 = "a2000003-0000-4000-8000-000000000003"
	QPAAgendarNuevaSesionS4  = "a1000004-0000-4000-8000-000000000004"
	QPAHoraProximaAtencionS4 = "a2000004-0000-4000-8000-000000000004"

	// Atención Psicosocial (SEG) — S1 / S4
	QSEGAgendarNuevaSesionS1  = "a1000005-0000-4000-8000-000000000005"
	QSEGHoraProximaAtencionS1 = "a2000005-0000-4000-8000-000000000005"
	QSEGAgendarNuevaSesionS4  = "a1000006-0000-4000-8000-000000000006"
	QSEGHoraProximaAtencionS4 = "a2000006-0000-4000-8000-000000000006"

	// Cierre — S1 / S4
	QCIEAgendarNuevaSesionS1  = "a1000007-0000-4000-8000-000000000007"
	QCIEHoraProximaAtencionS1 = "a2000007-0000-4000-8000-000000000007"
	QCIEAgendarNuevaSesionS4  = "a1000008-0000-4000-8000-000000000008"
	QCIEHoraProximaAtencionS4 = "a2000008-0000-4000-8000-000000000008"
)

// Fechas ya existentes (seed original) — documentadas aquí para el seed de agenda.
const (
	QPCFechaProximaS1  = "58ce2d34-24d2-4e73-bf95-26a2c608f8e6"
	QPCFechaProximaS4  = "fd2fb664-de83-4069-a161-6348dfef48bf"
	QPAFechaNuevaS1    = "64d63b79-edee-464b-be56-1104efd31a46"
	QPAFechaProximaS4  = "d68c7334-74bb-47b1-a2ee-f1d3da04627b"
	QSEGFechaNuevaS1   = "72ce49d2-f853-4f4a-9f1a-f795f4d514c3"
	QSEGFechaProximaS4 = "7040a37d-f346-4bf7-9b80-6286b9da62c5"
	QCIEFechaNuevaS1   = "1b4d09f0-e5ca-4b4e-9d74-428478a113c6"
	QCIEFechaProximaS4 = "3ce6ff5a-f139-4417-9255-155207e9a970"
)
