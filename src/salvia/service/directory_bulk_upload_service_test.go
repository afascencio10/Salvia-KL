package service

import (
	"bitsflow/internal/models"
	"bitsflow/internal/repository"
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeUpsertLocRepo struct {
	repository.LocationRepository
	depts  []models.LocationOption
	cities map[uint64][]models.CityLocationOption
}

func (f *fakeUpsertLocRepo) GetDepartments(context.Context) ([]models.LocationOption, error) {
	return f.depts, nil
}

func (f *fakeUpsertLocRepo) GetCities(_ context.Context, depID *uint64) ([]models.CityLocationOption, error) {
	return f.cities[*depID], nil
}

type fakeUpsertDirRepo struct {
	repository.DirectoryRepository
	existing []models.Directory
	created  []models.Directory
	updated  []models.Directory
}

func (f *fakeUpsertDirRepo) FindAll(context.Context) ([]models.Directory, error) { return f.existing, nil }

func (f *fakeUpsertDirRepo) Create(_ context.Context, d *models.Directory) error {
	f.created = append(f.created, *d)
	return nil
}

func (f *fakeUpsertDirRepo) Update(_ context.Context, d *models.Directory) error {
	f.updated = append(f.updated, *d)
	return nil
}

const upsertTestHeader = "id,departamento_carpeta,departamento,municipio_municipios_de_accion," +
	"nombre_de_la_institucion_entidad,otra_entidad_organizacion_cual," +
	"sector_especifico_corresponde_con_relacion_a,direccion," +
	"contacto_entidad_telefono_fijo_celular,e_mail_del_punto_focal_de,horarios_de_atencion\n"

// Cada fila reproduce un caso real del consolidado "Mapeo de actores" de SALVIA.
const upsertTestRows = "" +
	// fila 2: nueva — nombre real en "Otra entidad", tipo/sector del desplegable.
	"1,Amazonas,Amazonas,Leticia,Hospital,ESE Hospital San Rafael,Salud,Av. Internacional,3001234567,h@x.co,24 horas\n" +
	// fila 3: ya existe (nombre + municipio + dirección) → actualiza; email vacío no borra el registrado.
	"2,Valle del Cauca,Valle del Cauca,Cali,Comisaría de Familia,NA,Protección,Calle 5 # 10-20,222,,\n" +
	// fila 4: mismo nombre y municipio, otra sede (otra dirección) → nueva.
	"3,Valle del Cauca,Valle del Cauca,Cali,Comisaría de Familia,N/A,Protección,Carrera 1 # 2-3,333,,\n" +
	// fila 5: repetida de la fila 2 con otras tildes/mayúsculas → duplicada.
	"4,Amazonas,Amazonas,LETICIA,Hospital,ese hospital san rafael,Salud,Av Internacional,,,\n" +
	// fila 6: varios municipios en una sola fila.
	"5,Cauca,Cauca,\"Popayán, Timbío\",Fiscalía,NA,Justicia,Calle 1,,,\n" +
	// fila 7: municipio que no existe en el departamento.
	"6,Amazonas,Amazonas,Tarapacá,Policía Nacional,NA,Protección,Calle 2,,,\n" +
	// fila 8: Bogotá viene en la carpeta de Cundinamarca; se resuelve por la columna Departamento.
	"7,Cundinamarca,BOGOTA DC,Bogotá,Secretaria Distrital de la Mujer,NA,Protección,Calle 26,,,\n" +
	// fila 9: sin departamento diligenciado (usa la carpeta), tipo "Otra" clasificado por el nombre, sin dirección.
	"8,Amazonas,,Puerto Nariño,Otra,Comisaria de familia Puerto Nariño,Protección,N/A,,,\n" +
	// fila 10: tipo "Otra" sin decir cuál → sin nombre.
	"9,Amazonas,Amazonas,Leticia,Otra,NA,Protección,Calle 3,,,\n" +
	// fila 11: columnas corridas (municipio trae el tipo) → no debe caer en "Guadalajara de Buga" por la palabra DE.
	"10,Valle del Cauca,Valle del Cauca,Comisaria de familia,Protección,NA,Salud,Calle 4,,,\n"

func newUpsertTestService() (*directoryService, *fakeUpsertDirRepo) {
	loc := &fakeUpsertLocRepo{
		depts: []models.LocationOption{
			{Label: "AMAZONAS", Value: "91"},
			{Label: "BOGOTÁ, D.C.", Value: "11"},
			{Label: "CUNDINAMARCA", Value: "25"},
			{Label: "VALLE DEL CAUCA", Value: "76"},
			{Label: "CAUCA", Value: "19"},
		},
		cities: map[uint64][]models.CityLocationOption{
			91: {{Label: "LETICIA", ICode: "c-leticia"}, {Label: "PUERTO NARIÑO", ICode: "c-pnarino"}},
			11: {{Label: "BOGOTÁ, D.C.", ICode: "c-bogota"}},
			76: {{Label: "CALI", ICode: "c-cali"}, {Label: "GUADALAJARA DE BUGA", ICode: "c-buga"}},
			19: {{Label: "POPAYÁN", ICode: "c-popayan"}, {Label: "TIMBÍO", ICode: "c-timbio"}},
		},
	}
	oldPhone, oldEmail := "111", "viejo@x.co"
	dir := &fakeUpsertDirRepo{existing: []models.Directory{{
		ID: "d1", CityID: "c-cali", Name: "Comisaria de familia", Address: "Calle 5 # 10-20",
		Phone: &oldPhone, Email: &oldEmail, Type: models.DirectoryTypeComisariaFamilia, Active: true,
	}}}
	return &directoryService{dirRepo: dir, locRepo: loc}, dir
}

func TestUpsertFromCSV(t *testing.T) {
	svc, dir := newUpsertTestService()
	res, err := svc.UpsertFromCSV(context.Background(), strings.NewReader(upsertTestHeader+upsertTestRows))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if res.Total != 10 || res.Created != 4 || res.Updated != 1 || res.Failed != 5 {
		t.Fatalf("total/created/updated/failed = %d/%d/%d/%d, esperaba 10/4/1/5; errores: %+v",
			res.Total, res.Created, res.Updated, res.Failed, res.Errors)
	}

	wantErrors := map[int]string{
		5:  "duplicada de la fila 2",
		6:  "varios municipios",
		7:  "municipio 'Tarapacá' no encontrado",
		10: "no tiene nombre",
		11: "municipio 'Comisaria de familia' no encontrado",
	}
	for _, e := range res.Errors {
		want, ok := wantErrors[e.Row]
		if !ok || !strings.Contains(e.Error, want) {
			t.Errorf("fila %d: error %q, esperaba que contuviera %q", e.Row, e.Error, want)
		}
	}

	up := dir.updated[0]
	if up.ID != "d1" || *up.Phone != "222" || *up.Email != "viejo@x.co" || up.Sector == nil || *up.Sector != models.DirectorySectorProteccion {
		t.Errorf("actualización incorrecta: id=%s phone=%v email=%v sector=%v", up.ID, *up.Phone, *up.Email, up.Sector)
	}

	byName := map[string]models.Directory{}
	for _, d := range dir.created {
		byName[d.Name+"|"+d.Address] = d
	}
	checks := []struct{ key, city, dirType string }{
		{"ESE Hospital San Rafael|Av. Internacional", "c-leticia", models.DirectoryTypeHospital},
		{"Comisaría de Familia|Carrera 1 # 2-3", "c-cali", models.DirectoryTypeComisariaFamilia},
		{"Secretaria Distrital de la Mujer|Calle 26", "c-bogota", models.DirectoryTypeSecretariaGenero},
		{"Comisaria de familia Puerto Nariño|" + directoryNoAddress, "c-pnarino", models.DirectoryTypeComisariaFamilia},
	}
	for _, c := range checks {
		d, ok := byName[c.key]
		if !ok {
			t.Errorf("no se creó %q; creados: %+v", c.key, dir.created)
			continue
		}
		if d.CityID != c.city || d.Type != c.dirType || !d.Active {
			t.Errorf("%q: city=%s type=%s active=%v, esperaba city=%s type=%s active=true", c.key, d.CityID, d.Type, d.Active, c.city, c.dirType)
		}
	}
}

func TestUpsertFromCSV_MissingColumn(t *testing.T) {
	svc, _ := newUpsertTestService()
	_, err := svc.UpsertFromCSV(context.Background(), strings.NewReader("departamento,direccion\nAmazonas,Calle 1\n"))
	if !errors.Is(err, ErrDirectoryInvalidInput) {
		t.Fatalf("esperaba ErrDirectoryInvalidInput por columnas faltantes, obtuve %v", err)
	}
}

func TestMatchDirectoryType(t *testing.T) {
	cases := map[string]string{
		"Inspección de Policia":                    models.DirectoryTypeInspeccionesPolicia,
		"Centro zonal ICBF":                        models.DirectoryTypeCentroZonalICBF,
		"IPS/EPS":                                  models.DirectoryTypeIPSEPS,
		"Juez promiscuo":                           models.DirectoryTypeJuzgados,
		"Juzgados promiscuos municipales":          models.DirectoryTypeJuzgados,
		"Cabildos indigenas ":                      models.DirectoryTypeCabildoIndigena,
		"Casa LGBTIQ+":                             models.DirectoryTypeCasaLGBTIQ,
		"URI":                                      models.DirectoryTypeFiscalia,
		"Alcaldía de Uribia":                       models.DirectoryTypeAlcaldia, // URIBIA no es la sigla URI
		"Mesa Permanente de Concertación Indigena": models.DirectoryTypeOtra,
	}
	for raw, want := range cases {
		if got := matchDirectoryType(raw); got != want {
			t.Errorf("matchDirectoryType(%q) = %q, esperaba %q", raw, got, want)
		}
	}
}
