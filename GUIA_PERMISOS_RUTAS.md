# Guía de Permisos y Rutas — Módulo de Seguimiento

Cómo funciona el control de acceso en SALVIA y cómo agregar permisos a tus nuevas rutas.

---

## 1. Mapeo de roles del negocio vs roles en BD

El sistema usa códigos cortos de 2 letras almacenados en `security.role.role_code`:

| Rol de negocio | role_code | team (campo nuevo) |
|---|---|---|
| Superadmin | `ad` | (no aplica) |
| Supervisor General | `sv` | `RIESGO_BAJO` |
| Supervisor Integral | `sv` | `RIESGO_ALTO` |
| Agente General | `op` | `RIESGO_BAJO` |
| Agente Integral | `op` | `RIESGO_ALTO` |

Los roles `sv` y `op` se diferencian por el campo `general_user_team` en la tabla `security.general_user`. Este campo se carga en la sesión como `s.Team`.

---

## 2. Cómo registrar permisos para una nueva ruta

Abre `salvia/config/Menu.go` y agrega una entrada al mapa `PermissionsByRole`:

```go
// En salvia/config/Menu.go, dentro del mapa PermissionsByRole:

"mi_nueva_accion": {
    "ad": true,  // Superadmin
    "sv": true,  // Supervisores (General e Integral)
    "op": true,  // Agentes (General e Integral)
},
```

Si la acción es solo para agentes:
```go
"solo_agentes": {
    "op": true,
},
```

Si es solo para supervisores y admin:
```go
"solo_supervisores": {
    "ad": true,
    "sv": true,
},
```

---

## 3. Cómo aplicar el chequeo en un controlador Gin

### Caso simple: solo validar rol (sin diferenciar equipo)

```go
func MiHandler(c *gin.Context) {
    session := sessions.Default(c)
    s, err := utils.GetCommonSession(session.Get("userData").(string))
    if err != nil {
        c.JSON(401, gin.H{"error": "sesión inválida"})
        return
    }

    // Valida que el rol tenga permiso para esta acción
    if !utils.CheckPermission(salvia_config.PermissionsByRole, "mi_nueva_accion", s.CurrentRole, c) {
        return // CheckPermission ya envió el 401
    }

    // ... lógica del handler
}
```

### Caso con equipo: diferenciar General vs Integral

Usa `CheckPermissionWithTeam` cuando necesites que solo usuarios de un equipo específico accedan:

```go
func SeguimientosAreaHandler(c *gin.Context) {
    session := sessions.Default(c)
    s, err := utils.GetCommonSession(session.Get("userData").(string))
    if err != nil {
        c.JSON(401, gin.H{"error": "sesión inválida"})
        return
    }

    // Solo supervisores de RIESGO_ALTO pueden ver seguimientos de riesgo alto
    riskLevel := c.Query("risk_level") // "RIESGO_ALTO" o "RIESGO_BAJO"

    if !utils.CheckPermissionWithTeam(
        salvia_config.PermissionsByRole,
        "get_seguimientos_area",
        s.CurrentRole,
        s.Team,          // equipo del usuario en sesión
        riskLevel,       // equipo requerido por el recurso
        c,
    ) {
        return // ya envió 401 o 403
    }

    // ... lógica del handler
}
```

### Caso admin: acceso total sin importar equipo

`CheckPermissionWithTeam` ya maneja esto internamente — si el rol es `ad`, ignora la validación de equipo y deja pasar.

---

## 4. Ejemplo completo: proteger "Mis Seguimientos del Día"

```go
// En el controlador de seguimientos:

func (ctrl *FollowUpV2Controller) GetMisSeguimientosDia(c *gin.Context) {
    session := sessions.Default(c)
    s, err := utils.GetCommonSession(session.Get("userData").(string))
    if err != nil {
        c.JSON(401, gin.H{"error": "sesión inválida"})
        return
    }

    // Paso 1: Validar que el rol tenga permiso base
    if !utils.CheckPermission(salvia_config.PermissionsByRole, "get_mis_seguimientos_dia", s.CurrentRole, c) {
        return
    }

    // Paso 2: Filtrar seguimientos según el equipo del agente
    // Un agente de RIESGO_BAJO solo ve casos de riesgo bajo
    // Un agente de RIESGO_ALTO ve casos de riesgo alto
    items, err := ctrl.svc.GetByCaseIDAndTeam(c.Request.Context(), s.UserICode, s.Team)
    if err != nil {
        c.JSON(500, gin.H{"error": "error interno"})
        return
    }

    c.JSON(200, items)
}
```

---

## 5. Tabla de permisos del módulo de seguimiento

| Pantalla | Acción en Menu.go | ad | sv | op | Requiere team? |
|---|---|---|---|---|---|
| Detalle Caso | `get_seguimiento_detalle_caso` | si | si | si | No |
| Detalle Seguimiento | `get_seguimiento_detalle` | si | si | si | No |
| Formulario Dinámico | `get_seguimiento_formulario` | — | — | si | No (solo op) |
| Mis Seguimientos Día | `get_mis_seguimientos_dia` | — | — | si | Sí (filtra datos) |
| Seguimientos Área | `get_seguimientos_area` | si | si | — | Sí (filtra datos) |
| Generar Calendario | `generate_calendario_seguimiento` | si | si | si | No |

---

## 6. Migración SQL requerida

Ejecutar antes de usar el campo `team`:

```sql
ALTER TABLE security.general_user
    ADD COLUMN IF NOT EXISTS general_user_team VARCHAR(50);
```

Luego asignar el equipo a cada usuario desde el panel de administración o directamente en BD:

```sql
-- Ejemplo: asignar equipo a un operador
UPDATE security.general_user
SET general_user_team = 'RIESGO_ALTO'
WHERE general_user_login = 'agente_integral_01';
```

---

## 7. Resumen del flujo

```
Request HTTP
    → AuthMiddleware (valida sesión activa — ya existe, NO tocar)
    → CheckPermission (valida rol vs acción — mapa en Menu.go)
    → CheckPermissionWithTeam (opcional: valida equipo si aplica)
    → Handler (lógica de negocio, filtra datos por team si es necesario)
```
