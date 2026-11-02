package repository

import "testing"

func sp(s string) *string { return &s }

// WP-1 T7: regla única del estado derivado del reporte («Estado derivado del reporte»).
func TestReportStatusFrom_DerivedFromLinkedCase(t *testing.T) {
	cases := []struct {
		name        string
		agent       *string
		closeReason *string
		closeCat    *string
		legacy      string
		linked      *string
		want        string
	}{
		{"sin asignar", nil, nil, nil, "v", nil, "sin_asignar"},
		{"asignado", sp("a1"), nil, nil, "v", nil, "asignado"},
		{"cerrado válido por motivo ap", sp("a1"), sp("ap"), sp("valido"), "v", nil, "cerrado_valido"},
		{"cerrado inválido por motivo ni", nil, sp("ni"), sp("invalido"), "v", nil, "cerrado_invalido"},
		{"legado invalidado", nil, nil, nil, "i", nil, "cerrado_invalido"},
		{"caso en borrador", nil, nil, nil, "v", sp("bo"), "borrador_iniciado"},
		{"caso iniciado con agente asignado", sp("a1"), nil, nil, "v", sp("in"), "borrador_iniciado"},
		{"caso terminado", sp("a1"), nil, nil, "v", sp("te"), "cerrado_valido"},
		{"caso cerrado", nil, nil, nil, "v", sp("cd"), "cerrado_valido"},
		{"el motivo de cierre gana sobre el caso", nil, sp("co"), sp("invalido"), "v", sp("bo"), "cerrado_invalido"},
		{"el legado invalidado gana sobre el caso", nil, nil, nil, "i", sp("te"), "cerrado_invalido"},
		{"un caso de otro flujo no cambia el estado", sp("a1"), nil, nil, "v", sp("ra"), "asignado"},
	}
	for _, c := range cases {
		if got := reportStatusFrom(c.agent, c.closeReason, c.closeCat, c.legacy, c.linked); got != c.want {
			t.Errorf("%s: got %s, want %s", c.name, got, c.want)
		}
	}
}

func TestValidateConversion(t *testing.T) {
	ok := func() *reportConversionSource {
		return &reportConversionSource{ID: 1, ICode: "r1", ReportTypeCode: "v", HasContact: true}
	}
	cases := []struct {
		name   string
		mutate func(*reportConversionSource)
		linked []reportLinkedCase
		check  func(error) bool
	}{
		{"víctima con contacto", func(s *reportConversionSource) {}, nil, func(e error) bool { return e == nil }},
		{"tercero con contacto", func(s *reportConversionSource) { s.ReportTypeCode = "r" }, nil, func(e error) bool { return e == nil }},
		{"familiar", func(s *reportConversionSource) { s.ReportTypeCode = "fm" }, nil, func(e error) bool { return e == ErrReportNotConvertibleType }},
		{"hombre", func(s *reportConversionSource) { s.ReportTypeCode = "ho" }, nil, func(e error) bool { return e == ErrReportNotConvertibleType }},
		{"V1 sin código de tipo", func(s *reportConversionSource) { s.ReportTypeCode = "" }, nil, func(e error) bool { return e == ErrReportNotConvertibleType }},
		{"sin intento exitoso", func(s *reportConversionSource) { s.HasContact = false }, nil, func(e error) bool { return e == ErrReportNoContact }},
		{"reporte cerrado", func(s *reportConversionSource) { s.Closed = true }, nil, func(e error) bool { return e == ErrVictimContactAlreadyClosed }},
		{"con borrador bo", func(s *reportConversionSource) {}, []reportLinkedCase{{ICode: "c1", Status: "bo"}}, func(e error) bool {
			d, isDraft := e.(*ErrReportHasDraft)
			return isDraft && d.CaseICode == "c1"
		}},
		{"con borrador in (el más reciente)", func(s *reportConversionSource) {}, []reportLinkedCase{{ICode: "c2", Status: "in"}, {ICode: "c1", Status: "bo"}}, func(e error) bool {
			d, isDraft := e.(*ErrReportHasDraft)
			return isDraft && d.CaseICode == "c2"
		}},
		{"caso terminado = reporte cerrado", func(s *reportConversionSource) {}, []reportLinkedCase{{ICode: "c1", Status: "te"}}, func(e error) bool { return e == ErrVictimContactAlreadyClosed }},
		{"caso de otro flujo", func(s *reportConversionSource) {}, []reportLinkedCase{{ICode: "c1", Status: "ra"}}, func(e error) bool { return e == ErrReportHasLinkedCase }},
		{"cerrado gana sobre tipo", func(s *reportConversionSource) { s.Closed = true; s.ReportTypeCode = "fm" }, nil, func(e error) bool { return e == ErrVictimContactAlreadyClosed }},
		{"tipo gana sobre contacto", func(s *reportConversionSource) { s.ReportTypeCode = "fm"; s.HasContact = false }, nil, func(e error) bool { return e == ErrReportNotConvertibleType }},
	}
	for _, c := range cases {
		src := ok()
		c.mutate(src)
		if err := validateConversion(src, c.linked); !c.check(err) {
			t.Errorf("%s: error inesperado: %v", c.name, err)
		}
	}
}
