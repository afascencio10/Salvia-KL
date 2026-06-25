━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
🟢 EVENTO: Cuando presiona Crear caso o Ver
   Tipo: User Interaction
   Código: E-04
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Disparado por: clic en enlace de `menuToolsContactTable` en columna Acciones


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ACCIÓN (A) — Ver reporte
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

URL: `/salvia/primer_contacto/{victimContact.icode}/`

  → VictimContactGET con id
  → Template según rol:
       sv → get_victim_contact_sv_v2 (form2) o sv_v1 (form1 legacy)
       ro → get_victim_contact o get_victim_contact_v1

  Solo lectura del detalle del reporte.


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ACCIÓN (B) — Crear caso
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

URL: `/salvia/casos/{victimContact.icode}/nuevo`

  → VictimCasePOST_GET
  → Permiso: set_victim_case (sv ✓, ro ✓)
  → GetVictimContactByICode(id)
  → LoadFromVictimContactForm2(...) pre-carga set_victim_case.html
  → Usuario completa y guarda → SetVictimCase crea victim_case

  Tras crear el caso, el reporte deja de aparecer en el listado
  (LEFT JOIN victim_case IS NULL).


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  CAMBIO EN Menu.go (habilitar Crear caso en sv)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Antes (sv): solo icono Ver
Después (sv): Ver + Crear caso (igual que op/ro/et)

HTML de la fila (sin cambio — usa v-for menuToolsContactTable):

```html
<a :href="`${item.path}/${victimContact.icode}/${item.action}`">
  ${item.label}
</a>
```

  item.action vacío  → Ver   → /salvia/primer_contacto/{icode}/
  item.action "nuevo" → Crear → /salvia/casos/{icode}/nuevo


━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  FASE 2 OPCIONAL — Detalle sv
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agregar en `get_victim_contact_sv_v2.html` botón "Crear caso"
(mismo patrón que get_victim_contact.html) para no depender solo
de la acción en la tabla.

  → FIN EJECUCIÓN ✓
