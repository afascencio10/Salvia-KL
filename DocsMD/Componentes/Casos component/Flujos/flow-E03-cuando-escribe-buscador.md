━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando el usuario escribe en el buscador
   Tipo: User Interaction
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: input en el campo SearchInput de FilterBar
               El evento se ejecuta en cada pulsación de tecla (o con debounce)

INPUT: {
  searchText:   texto escrito por el usuario   → valor del <input type="text">
}


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FRONTEND — casos-component.js  (filtrado del lado del cliente)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━


PASO 1 — Actualizar el estado de búsqueda

  searchText = searchText.trim()


PASO 2 — Verificar si hay texto de búsqueda

  SI searchText === "":
    → filteredCases = cases    // restaurar lista completa del filtro activo
    → FIN EJECUCIÓN ✓

  SI searchText !== "":
    → CONTINÚA PASO 3


PASO 3 — Filtrar cases del lado del cliente

  term = searchText.toLowerCase()

  filteredCases = cases.filter(caso => {

    matchICode = caso.i_code.toLowerCase().includes(term)
    matchPhone = caso.phone.toLowerCase().includes(term)
                 // caso.phone → número de teléfono de la víctima

    return matchICode || matchPhone
  })


PASO 4 — Actualizar la vista

  SI filteredCases.length === 0:
    → Mostrar EmptyState "No se encontraron casos"

  SI filteredCases.length > 0:
    → Renderizar tabla con filteredCases

  → FIN EJECUCIÓN ✓


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⚠️  GAPS — Información pendiente
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

| Variable / decisión                                                                  | Paso afectado |
|--------------------------------------------------------------------------------------|---------------|
| Campo exacto del teléfono de la víctima en el objeto de caso (¿phone? ¿contactPhone?) | PASO 3      |
| ¿El teléfono viene de victim_case, victim_case_form1 u otra tabla?                   | PASO 3        |
| ¿Se aplica debounce al input? ¿Cuántos ms? (recomendado: 300ms)                      | PASO 1        |
| ¿La búsqueda es case-sensitive o normaliza acentos / caracteres especiales?          | PASO 3        |
