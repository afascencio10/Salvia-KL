package service

import (
	"bitsflow/common/utils"
	"bitsflow/internal/models"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Columnas de rutas_atencion_territorial.csv, del paquete "BD consolidada mapeo de
// actores" que entrega SALVIA (su diccionario_de_datos.csv describe cada una).
const (
	colDepartamentoCarpeta = "departamento_carpeta"
	colDepartamento        = "departamento"
	colMunicipio           = "municipio_municipios_de_accion"
	colTipoInstitucion     = "nombre_de_la_institucion_entidad"
	colOtraEntidad         = "otra_entidad_organizacion_cual"
	colSector              = "sector_especifico_corresponde_con_relacion_a"
	colDireccion           = "direccion"
	colTelefono            = "contacto_entidad_telefono_fijo_celular"
	colEmail               = "e_mail_del_punto_focal_de"
	colHorario             = "horarios_de_atencion"
)

var directoryCSVColumns = []string{
	colDepartamentoCarpeta, colDepartamento, colMunicipio, colTipoInstitucion, colOtraEntidad,
	colSector, colDireccion, colTelefono, colEmail, colHorario,
}

// directoryNoAddress — address es NOT NULL, pero ~560 filas del archivo no traen dirección.
const directoryNoAddress = "No registra"

// minLocationScore — con utils.ScoreCityNameMatch, >= 500 exige coincidencia exacta, de
// prefijo, de contenido o de todas las palabras. Evita que una fila con columnas corridas
// ("Comisaria de familia" en la columna de municipio) caiga en cualquier municipio que
// comparta la palabra "DE".
const minLocationScore = 500

// directoryEmptyValues — textos (normalizados) que el archivo usa para "sin dato".
var directoryEmptyValues = map[string]bool{
	"": true, "NA": true, "N A": true, "NO": true, "0": true, "NO APLICA": true,
	"SIN INFORMACION": true, "NINGUNA": true, "NINGUNO": true, "NO REGISTRA": true, "NO TIENE": true,
	// Valores de la columna "Nivel" que quedaron en "Otra entidad ¿Cuál?" en algunos departamentos.
	"ORDEN MUNICIPAL": true, "ORDEN DEPARTAMENTAL": true, "ORDEN NACIONAL": true,
}

type DirectoryUpsertRowError struct {
	Row          int    `json:"row"`
	Name         string `json:"name,omitempty"`
	Municipality string `json:"municipality,omitempty"`
	Error        string `json:"error"`
}

type DirectoryUpsertResult struct {
	Total   int                       `json:"total"`
	Created int                       `json:"created"`
	Updated int                       `json:"updated"`
	Failed  int                       `json:"failed"`
	Errors  []DirectoryUpsertRowError `json:"errors"`
}

// UpsertFromCSV procesa el archivo fila por fila, sin transacción global: una fila que
// falla se reporta en Errors y no afecta a las demás.
//
// Criterio de coincidencia (acordado con negocio): nombre efectivo + municipio + dirección,
// normalizados (sin tildes, mayúsculas, sin puntuación). Nombre efectivo = "Otra
// entidad/organización ¿Cuál?" si trae un nombre real; si no, "Nombre de la
// institución/entidad" (que en el archivo es en realidad el tipo: "Comisaría de Familia").
// La dirección entra en la llave porque un mismo municipio tiene varias sedes con el mismo
// nombre (p. ej. 11 "Comisaría de Familia" en Cali).
func (s *directoryService) UpsertFromCSV(ctx context.Context, r io.Reader) (*DirectoryUpsertResult, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("%w: no se pudo leer el encabezado del CSV: %v", ErrDirectoryInvalidInput, err)
	}
	col := make(map[string]int, len(header))
	for i, h := range header {
		col[strings.TrimSpace(strings.TrimPrefix(h, "\ufeff"))] = i
	}
	for _, name := range directoryCSVColumns {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("%w: al CSV le falta la columna '%s'", ErrDirectoryInvalidInput, name)
		}
	}

	u, err := s.newDirectoryUpsert(ctx)
	if err != nil {
		return nil, err
	}

	result := &DirectoryUpsertResult{Errors: []DirectoryUpsertRowError{}}
	// row = número de fila al abrir el CSV en Excel (la 1 es el encabezado).
	for row := 2; ; row++ {
		rec, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		var parseErr *csv.ParseError
		if err != nil && !errors.As(err, &parseErr) {
			return nil, fmt.Errorf("error leyendo el archivo en la fila %d: %w", row, err)
		}
		result.Total++
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DirectoryUpsertRowError{Row: row, Error: "fila mal formada: " + err.Error()})
			continue
		}

		get := func(name string) string {
			if i := col[name]; i < len(rec) {
				return rec[i]
			}
			return ""
		}
		created, name, err := u.upsertRow(ctx, row, get)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, DirectoryUpsertRowError{
				Row:          row,
				Name:         name,
				Municipality: cleanDirectoryValue(get(colMunicipio)),
				Error:        err.Error(),
			})
			continue
		}
		if created {
			result.Created++
		} else {
			result.Updated++
		}
	}
	return result, nil
}

type directoryUpsert struct {
	s          *directoryService
	deptIDs    []uint64
	deptNorms  []string
	cities     map[uint64][]models.CityLocationOption
	cityNorms  map[uint64][]string
	byKey      map[string]*models.Directory
	seenInFile map[string]int
}

func (s *directoryService) newDirectoryUpsert(ctx context.Context) (*directoryUpsert, error) {
	departments, err := s.locRepo.GetDepartments(ctx)
	if err != nil {
		return nil, err
	}
	existing, err := s.dirRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	u := &directoryUpsert{
		s:          s,
		cities:     map[uint64][]models.CityLocationOption{},
		cityNorms:  map[uint64][]string{},
		byKey:      make(map[string]*models.Directory, len(existing)),
		seenInFile: map[string]int{},
	}
	for _, d := range departments {
		id, err := strconv.ParseUint(d.Value, 10, 64)
		if err != nil {
			continue
		}
		u.deptIDs = append(u.deptIDs, id)
		u.deptNorms = append(u.deptNorms, utils.NormalizeLocationName(d.Label))
	}
	for i := range existing {
		d := &existing[i]
		u.byKey[directoryKey(d.Name, d.CityID, d.Address)] = d
	}
	return u, nil
}

func (u *directoryUpsert) upsertRow(ctx context.Context, row int, get func(string) string) (created bool, name string, err error) {
	name = effectiveDirectoryName(get(colTipoInstitucion), get(colOtraEntidad))
	if name == "" {
		return false, "", errors.New("la fila no tiene nombre de institución/entidad")
	}
	municipio := cleanDirectoryValue(get(colMunicipio))
	if municipio == "" {
		return false, name, errors.New("la fila no tiene municipio")
	}

	// La columna "Departamento" es la que diligenció cada equipo (distingue p. ej. Bogotá
	// dentro de la carpeta de Cundinamarca); si no resuelve, se usa la carpeta de origen.
	depID, ok := u.resolveDepartment(get(colDepartamento), get(colDepartamentoCarpeta))
	if !ok {
		return false, name, fmt.Errorf("departamento '%s' no encontrado", cleanDirectoryValue(get(colDepartamentoCarpeta)))
	}
	city, err := u.resolveCity(ctx, depID, municipio)
	if err != nil {
		return false, name, err
	}

	address := cleanDirectoryValue(get(colDireccion))
	if address == "" {
		address = directoryNoAddress
	}
	key := directoryKey(name, city.ICode, address)
	if prev, dup := u.seenInFile[key]; dup {
		return false, name, fmt.Errorf("duplicada de la fila %d (mismo nombre, municipio y dirección)", prev)
	}
	u.seenInFile[key] = row

	dirType := matchDirectoryType(get(colTipoInstitucion))
	if dirType == models.DirectoryTypeOtra {
		dirType = matchDirectoryType(name)
	}
	sector := matchDirectorySector(get(colSector))
	phone := optionalDirectoryValue(get(colTelefono))
	email := optionalDirectoryValue(get(colEmail))
	hours := optionalDirectoryValue(get(colHorario))

	if d, found := u.byKey[key]; found {
		// Un dato vacío en el archivo no borra lo ya registrado: parte de las entidades
		// las cargó SOG antes con campos que este archivo no trae. Tampoco se toca
		// Active: habilitar/deshabilitar sigue siendo manual desde la pantalla.
		d.Name, d.Address, d.DepartmentID = name, address, &depID
		if phone != nil {
			d.Phone = phone
		}
		if email != nil {
			d.Email = email
		}
		if hours != nil {
			d.OpeningHours = hours
		}
		if sector != nil {
			d.Sector = sector
		}
		if dirType != models.DirectoryTypeOtra || d.Type == "" {
			d.Type = dirType
		}
		return false, name, u.s.dirRepo.Update(ctx, d)
	}

	d := &models.Directory{
		CityID:       city.ICode,
		DepartmentID: &depID,
		Name:         name,
		Address:      address,
		Phone:        phone,
		Email:        email,
		OpeningHours: hours,
		Type:         dirType,
		Sector:       sector,
		Active:       true,
	}
	if err := u.s.dirRepo.Create(ctx, d); err != nil {
		return false, name, err
	}
	u.byKey[key] = d
	return true, name, nil
}

func (u *directoryUpsert) resolveDepartment(values ...string) (uint64, bool) {
	for _, v := range values {
		if i, _ := bestLocationMatch(v, u.deptNorms); i >= 0 {
			return u.deptIDs[i], true
		}
	}
	return 0, false
}

func (u *directoryUpsert) resolveCity(ctx context.Context, depID uint64, municipio string) (*models.CityLocationOption, error) {
	cities, ok := u.cities[depID]
	if !ok {
		var err error
		if cities, err = u.s.locRepo.GetCities(ctx, &depID); err != nil {
			return nil, err
		}
		norms := make([]string, len(cities))
		for i, c := range cities {
			norms[i] = utils.NormalizeLocationName(c.Label)
		}
		u.cities[depID], u.cityNorms[depID] = cities, norms
	}

	i, score := bestLocationMatch(municipio, u.cityNorms[depID])
	if score < 1000 && looksLikeSeveralMunicipalities(municipio) {
		return nil, fmt.Errorf("la fila trae varios municipios ('%s'); debe venir uno por fila", municipio)
	}
	if i < 0 {
		return nil, fmt.Errorf("municipio '%s' no encontrado en el departamento", municipio)
	}
	return &cities[i], nil
}

// bestLocationMatch devuelve el índice de la mejor etiqueta (ya normalizada) para input y su
// puntaje, o -1 si ninguna alcanza minLocationScore. Empate → la etiqueta más corta, igual
// que LocationRepository.FindBestCityMatchByName.
func bestLocationMatch(input string, normalizedLabels []string) (int, int) {
	in := utils.NormalizeLocationName(input)
	best, bestScore := -1, 0
	for i, label := range normalizedLabels {
		score := utils.ScoreCityNameMatch(in, label)
		if score == 0 {
			continue
		}
		if score > bestScore || (score == bestScore && len(label) < len(normalizedLabels[best])) {
			best, bestScore = i, score
		}
	}
	if bestScore < minLocationScore {
		return -1, bestScore
	}
	return best, bestScore
}

func looksLikeSeveralMunicipalities(raw string) bool {
	return strings.ContainsAny(raw, ",;/\n") || strings.Contains(" "+utils.NormalizeLocationName(raw)+" ", " Y ")
}

func directoryKey(name, cityICode, address string) string {
	return utils.NormalizeLocationName(name) + "|" + cityICode + "|" + utils.NormalizeLocationName(address)
}

// effectiveDirectoryName — "Otra entidad ¿Cuál?" si trae un nombre real; si no, el tipo.
// Un tipo "Otra" sin especificar cuál no es un nombre: se devuelve "" (fila rechazada).
func effectiveDirectoryName(tipo, otra string) string {
	name := cleanDirectoryValue(otra)
	if name == "" {
		name = cleanDirectoryValue(tipo)
	}
	if utils.NormalizeLocationName(name) == "OTRA" {
		return ""
	}
	return name
}

// cleanDirectoryValue colapsa espacios/saltos de línea y vacía los "sin dato" del archivo.
func cleanDirectoryValue(v string) string {
	v = strings.Join(strings.Fields(v), " ")
	if directoryEmptyValues[utils.NormalizeLocationName(v)] {
		return ""
	}
	return v
}

func optionalDirectoryValue(v string) *string {
	if c := cleanDirectoryValue(v); c != "" {
		return &c
	}
	return nil
}

// directoryRule — todas las palabras deben aparecer en el texto normalizado. Se evalúan en
// orden y gana la primera: las específicas van antes ("Inspección de Policía" antes que
// "Policía", "Centro zonal ICBF" antes que "ICBF").
type directoryRule struct {
	words []string
	code  string
}

var directoryTypeRules = []directoryRule{
	{[]string{"INSPECCION"}, models.DirectoryTypeInspeccionesPolicia},
	{[]string{"CENTRO", "ZONAL"}, models.DirectoryTypeCentroZonalICBF},
	{[]string{"COMISARIA"}, models.DirectoryTypeComisariaFamilia},
	{[]string{"FISCALIA"}, models.DirectoryTypeFiscalia},
	{[]string{"URI"}, models.DirectoryTypeFiscalia},
	{[]string{"SECRETARIA", "SALUD"}, models.DirectoryTypeSecretariaSalud},
	{[]string{"SECRETARIA", "EDUCACION"}, models.DirectoryTypeSecretariaEducacion},
	{[]string{"GENERO"}, models.DirectoryTypeSecretariaGenero},
	{[]string{"MUJER"}, models.DirectoryTypeSecretariaGenero},
	{[]string{"CAIVAS"}, models.DirectoryTypeCAIVAS},
	{[]string{"CAVIF"}, models.DirectoryTypeCAVIF},
	{[]string{"ICBF"}, models.DirectoryTypeICBF},
	{[]string{"BIENESTAR", "FAMILIAR"}, models.DirectoryTypeICBF},
	{[]string{"POLICIA"}, models.DirectoryTypePoliciaNacional},
	{[]string{"PERSONERIA"}, models.DirectoryTypePersoneria},
	{[]string{"DEFENSORIA"}, models.DirectoryTypeDefensoriaPueblo},
	{[]string{"PROCURADURIA"}, models.DirectoryTypeProcuraduria},
	{[]string{"JUSTICIA", "PROPIA"}, models.DirectoryTypeJusticiaPropia},
	{[]string{"CASA", "JUSTICIA"}, models.DirectoryTypeCasaJusticia},
	{[]string{"JUZGADO"}, models.DirectoryTypeJuzgados},
	{[]string{"JUEZ"}, models.DirectoryTypeJuzgados},
	{[]string{"HOSPITAL"}, models.DirectoryTypeHospital},
	{[]string{"CENTRO", "SALUD"}, models.DirectoryTypeCentroSalud},
	{[]string{"PUESTO", "SALUD"}, models.DirectoryTypeCentroSalud},
	{[]string{"IPS"}, models.DirectoryTypeIPSEPS},
	{[]string{"EPS"}, models.DirectoryTypeIPSEPS},
	{[]string{"CLINICA"}, models.DirectoryTypeClinicaPrivada},
	{[]string{"INSTITUCION", "EDUCATIVA"}, models.DirectoryTypeInstitucionesEducativas},
	{[]string{"COLEGIO"}, models.DirectoryTypeInstitucionesEducativas},
	{[]string{"INTEGRATE"}, models.DirectoryTypeCentroIntegrate},
	{[]string{"LGBT"}, models.DirectoryTypeCasaLGBTIQ},
	{[]string{"LGBTI"}, models.DirectoryTypeCasaLGBTIQ},
	{[]string{"CABILDO"}, models.DirectoryTypeCabildoIndigena},
	{[]string{"UARIV"}, models.DirectoryTypeUARIV},
	{[]string{"UNP"}, models.DirectoryTypeUNP},
	{[]string{"GOBERNACION"}, models.DirectoryTypeGobernacion},
	{[]string{"ALCALDIA"}, models.DirectoryTypeAlcaldia},
	{[]string{"LINEA"}, models.DirectoryTypeLineasEmergencia},
	{[]string{"URGENCIA"}, models.DirectoryTypeRedUrgencias},
}

var directorySectorRules = []directoryRule{
	{[]string{"JUSTICIA", "PROPIA"}, models.DirectorySectorJusticiaPropia},
	{[]string{"SALUD"}, models.DirectorySectorSalud},
	{[]string{"JUSTICIA"}, models.DirectorySectorJusticia},
	{[]string{"JUTICIA"}, models.DirectorySectorJusticia},
	{[]string{"JUSTICICA"}, models.DirectorySectorJusticia},
	{[]string{"PROTECCION"}, models.DirectorySectorProteccion},
	{[]string{"MINISTERIO", "PUBLICO"}, models.DirectorySectorMinisterioPublico},
	{[]string{"SERVICIO"}, models.DirectorySectorServiciosEspecializados},
	{[]string{"EDUCACION"}, models.DirectorySectorEducacion},
	{[]string{"ATENCION", "HUMANITARIA"}, models.DirectorySectorAtencionHumanitaria},
}

// matchDirectoryType clasifica el texto libre del archivo; sin coincidencia → "Otra".
// ponytail: reglas por palabra clave, no un clasificador; si aparecen tipos nuevos en el
// archivo, se agrega una regla aquí y su código en models.DirectoryTypeLabels.
func matchDirectoryType(raw string) string {
	if code := matchDirectoryRule(raw, directoryTypeRules); code != "" {
		return code
	}
	return models.DirectoryTypeOtra
}

// matchDirectorySector devuelve nil si el texto no corresponde a un sector conocido (el
// campo es opcional; hay filas con la dirección o el municipio corridos a esta columna).
func matchDirectorySector(raw string) *string {
	if code := matchDirectoryRule(raw, directorySectorRules); code != "" {
		return &code
	}
	return nil
}

func matchDirectoryRule(raw string, rules []directoryRule) string {
	tokens := strings.Fields(utils.NormalizeLocationName(raw))
	for _, rule := range rules {
		if allDirectoryWordsPresent(tokens, rule.words) {
			return rule.code
		}
	}
	return ""
}

func allDirectoryWordsPresent(tokens, words []string) bool {
	for _, w := range words {
		found := false
		for _, t := range tokens {
			// Siglas cortas (URI, IPS, UNP, ICBF...) exigen palabra exacta para no confundir
			// "URIBIA" con URI; las largas aceptan plural ("JUZGADOS", "INSPECCIONES").
			if t == w || (len(w) >= 5 && strings.HasPrefix(t, w)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
