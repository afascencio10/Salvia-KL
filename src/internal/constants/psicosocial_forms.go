// Package constants agrupa IDs "bien conocidos" generados por los scripts de seed,
// para evitar tenerlos repetidos como strings mágicos en repositorios/servicios/controllers.
package constants

// IDs de los 4 formularios psicosociales, generados por cmd/seed/seed_psicosocial.sql
// (ejecutado el 2026-07-14 contra la base de Supabase). Si el seed se vuelve a ejecutar
// contra otra base, estos valores deben actualizarse con los UUIDs que arroje el
// RAISE NOTICE final del script.
const (
	FormIDPrimerContacto      = "439b57e6-07ea-4da4-9721-8ed28c6ca43f" // Primer Contacto Psicosocial
	FormIDPrimeraAtencion     = "501ab3d7-8382-4447-96a6-f463152693c6" // Primera Atención Psicosocial
	FormIDAtencionPsicosocial = "7a7b61b3-7bc3-4088-a74d-0974e54a3563" // Atención Psicosocial (antes "Seguimiento")
	FormIDCierre              = "c31026f8-7ce2-49b9-8990-04b44ed4513c" // Cierre Psicosocial
)

// PsicosocialFormKey identifica de forma legible cuál de los 4 formularios está activo.
type PsicosocialFormKey string

const (
	PsicosocialFormPrimerContacto      PsicosocialFormKey = "PRIMER_CONTACTO"
	PsicosocialFormPrimeraAtencion     PsicosocialFormKey = "PRIMERA_ATENCION"
	PsicosocialFormAtencionPsicosocial PsicosocialFormKey = "ATENCION_PSICOSOCIAL"
	PsicosocialFormCierre              PsicosocialFormKey = "CIERRE"
)

// PsicosocialFormIDByKey mapea la clave legible al UUID real del formulario.
var PsicosocialFormIDByKey = map[PsicosocialFormKey]string{
	PsicosocialFormPrimerContacto:      FormIDPrimerContacto,
	PsicosocialFormPrimeraAtencion:     FormIDPrimeraAtencion,
	PsicosocialFormAtencionPsicosocial: FormIDAtencionPsicosocial,
	PsicosocialFormCierre:              FormIDCierre,
}

// PsicosocialFormKeyByID es el mapeo inverso: UUID → clave legible.
var PsicosocialFormKeyByID = map[string]PsicosocialFormKey{
	FormIDPrimerContacto:      PsicosocialFormPrimerContacto,
	FormIDPrimeraAtencion:     PsicosocialFormPrimeraAtencion,
	FormIDAtencionPsicosocial: PsicosocialFormAtencionPsicosocial,
	FormIDCierre:              PsicosocialFormCierre,
}
