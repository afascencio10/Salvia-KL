━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona un badge de "Oficios por estado" (En gestión / Radicado / Respondido)
   Tipo: User Interaction
   Función: abrirModalOficios(bucket)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

INPUT: {
  oficios: array ya cargado en E01                        → estado del componente
  bucket: 'En gestión' | 'Radicado' | 'Respondido'        → click del usuario sobre el badge
}

PASO 1 — Filtrar los oficios según el bucket seleccionado

  SEGÚN bucket:
    CASO 'En gestión':
      → filasFiltradas = oficios.filter(o => o.state !== 'radicado' && o.state !== 'respondido')
      → CONTINÚA
    CASO 'Radicado':
      → filasFiltradas = oficios.filter(o => o.state === 'radicado')
      → CONTINÚA
    CASO 'Respondido':
      → filasFiltradas = oficios.filter(o => o.state === 'respondido')
      → CONTINÚA

PASO 2 — Construir el objeto modal

  modal = {
    tipo: 'oficios',
    titulo: `Oficios — ${bucket}`,
    filas: filasFiltradas,
  }

PASO 3 — Renderizar el modal

  → Misma tabla que flow-E07 (Entidad | Estado | Fecha), acotada al bucket
  → Ver flow-E09 para el cierre

→ FIN EJECUCIÓN ✓
