// Package constants — IDs fijos de la sección OSIGHN (Sep 2026).
//
// Serie c1xxxxxx = form_section OSIGHN (order = 2 en cada form)
// Serie c2xxxxxx = pregunta gatillo en Contacto S1
// Serie c3xxxxxx = preguntas dentro de la sección OSIGHN
//
// Sufijo de form en c1/c2:
//
//	1 = Primer Contacto
//	2 = Primera Atención
//	3 = Atención Psicosocial (SEG)
//	4 = Cierre
//
// Preguntas c3: últimos 3 dígitos = FQQ (F=form 1–4, QQ=pregunta).
// QQ 01–02 (tipo/número de documento) se retiraron el 2026-09-17; quedan 03–10.
// Ver migrate_psicosocial_osighn_sep2026.sql y
// migrate_psicosocial_osighn_drop_doc_sep2026.sql.
package constants

const (
	SectionPCOSIGHN  = "c1000001-0000-4000-8000-000000000001"
	SectionPAOSIGHN  = "c1000002-0000-4000-8000-000000000002"
	SectionSEGOSIGHN = "c1000003-0000-4000-8000-000000000003"
	SectionCIEOSIGHN = "c1000004-0000-4000-8000-000000000004"
)

const (
	QPCTieneOSIGHN  = "c2000001-0000-4000-8000-000000000001"
	QPATieneOSIGHN  = "c2000002-0000-4000-8000-000000000002"
	QSEGTieneOSIGHN = "c2000003-0000-4000-8000-000000000003"
	QCIETieneOSIGHN = "c2000004-0000-4000-8000-000000000004"
)

// Triggers de visibilidad de la pregunta gatillo (ya existentes en Contacto).
const (
	QPCContinuarPrimeraAtencion = "aa46351a-f6c0-4666-af3e-e0f1f7b6dc8c"
	QPAEsAtencionOSoloContacto  = "3117fd01-6a31-4595-9888-dcdcc96d84e7"
	QSEGEsAtencionOSoloContacto = "16de7674-4349-4acc-9d0d-656d1d5e5f10"
	QCIEEsAtencionOSoloContacto = "be185336-e025-4768-a1a2-3ab9e2281450"
)

// IDs retirados 2026-09-17 (tipo y número de documento). Las filas quedan
// con deleted_at; no reinsertar.
const (
	QPCOSIGHNTipoDocumento    = "c3000001-0000-4000-8000-000000000101"
	QPCOSIGHNNumeroDocumento  = "c3000001-0000-4000-8000-000000000102"
	QPAOSIGHNTipoDocumento    = "c3000001-0000-4000-8000-000000000201"
	QPAOSIGHNNumeroDocumento  = "c3000001-0000-4000-8000-000000000202"
	QSEGOSIGHNTipoDocumento   = "c3000001-0000-4000-8000-000000000301"
	QSEGOSIGHNNumeroDocumento = "c3000001-0000-4000-8000-000000000302"
	QCIEOSIGHNTipoDocumento   = "c3000001-0000-4000-8000-000000000401"
	QCIEOSIGHNNumeroDocumento = "c3000001-0000-4000-8000-000000000402"
)

// Primer Contacto — sección OSIGHN (Q03–Q10 vivas, order 1–8).
const (
	QPCOSIGHNInfoNombres       = "c3000001-0000-4000-8000-000000000103"
	QPCOSIGHNNombreIdentitario = "c3000001-0000-4000-8000-000000000104"
	QPCOSIGHNNombreJuridico    = "c3000001-0000-4000-8000-000000000105"
	QPCOSIGHNPronombres        = "c3000001-0000-4000-8000-000000000106"
	QPCOSIGHNIdentidadGenero   = "c3000001-0000-4000-8000-000000000107"
	QPCOSIGHNIdentidadCual     = "c3000001-0000-4000-8000-000000000108"
	QPCOSIGHNSexoAsignado      = "c3000001-0000-4000-8000-000000000109"
	QPCOSIGHNOrientacionSexual = "c3000001-0000-4000-8000-000000000110"
)

// Primera Atención — sección OSIGHN.
const (
	QPAOSIGHNInfoNombres       = "c3000001-0000-4000-8000-000000000203"
	QPAOSIGHNNombreIdentitario = "c3000001-0000-4000-8000-000000000204"
	QPAOSIGHNNombreJuridico    = "c3000001-0000-4000-8000-000000000205"
	QPAOSIGHNPronombres        = "c3000001-0000-4000-8000-000000000206"
	QPAOSIGHNIdentidadGenero   = "c3000001-0000-4000-8000-000000000207"
	QPAOSIGHNIdentidadCual     = "c3000001-0000-4000-8000-000000000208"
	QPAOSIGHNSexoAsignado      = "c3000001-0000-4000-8000-000000000209"
	QPAOSIGHNOrientacionSexual = "c3000001-0000-4000-8000-000000000210"
)

// Atención Psicosocial (SEG) — sección OSIGHN.
const (
	QSEGOSIGHNInfoNombres       = "c3000001-0000-4000-8000-000000000303"
	QSEGOSIGHNNombreIdentitario = "c3000001-0000-4000-8000-000000000304"
	QSEGOSIGHNNombreJuridico    = "c3000001-0000-4000-8000-000000000305"
	QSEGOSIGHNPronombres        = "c3000001-0000-4000-8000-000000000306"
	QSEGOSIGHNIdentidadGenero   = "c3000001-0000-4000-8000-000000000307"
	QSEGOSIGHNIdentidadCual     = "c3000001-0000-4000-8000-000000000308"
	QSEGOSIGHNSexoAsignado      = "c3000001-0000-4000-8000-000000000309"
	QSEGOSIGHNOrientacionSexual = "c3000001-0000-4000-8000-000000000310"
)

// Cierre — sección OSIGHN.
const (
	QCIEOSIGHNInfoNombres       = "c3000001-0000-4000-8000-000000000403"
	QCIEOSIGHNNombreIdentitario = "c3000001-0000-4000-8000-000000000404"
	QCIEOSIGHNNombreJuridico    = "c3000001-0000-4000-8000-000000000405"
	QCIEOSIGHNPronombres        = "c3000001-0000-4000-8000-000000000406"
	QCIEOSIGHNIdentidadGenero   = "c3000001-0000-4000-8000-000000000407"
	QCIEOSIGHNIdentidadCual     = "c3000001-0000-4000-8000-000000000408"
	QCIEOSIGHNSexoAsignado      = "c3000001-0000-4000-8000-000000000409"
	QCIEOSIGHNOrientacionSexual = "c3000001-0000-4000-8000-000000000410"
)
