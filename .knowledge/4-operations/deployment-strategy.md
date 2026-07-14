---
okf_version: "1.0"
type: Deployment_Strategy
title: "Estrategia de Despliegue"
description: "Definición de la estrategia de despliegue Blue/Green, rollback automático, canary releases y feature flags para la plataforma."
owner: "@sre-squad"
status: active
tags: [ops, deployment, blue-green, rollback, canary]
last_updated: "2026-06-23"
---

# Estrategia de Despliegue

## 1. Modelo de Despliegue: Blue/Green

```text
                   Load Balancer
                       │
              ┌────────┴────────┐
              │                 │
         ┌────▼────┐       ┌───▼─────┐
         │  BLUE   │       │  GREEN  │
         │ (live)  │       │ (idle)  │
         │  v1.2.0 │       │  v1.3.0 │
         └─────────┘       └─────────┘
              │                 │
         ┌────▼────┐       ┌───▼─────┐
         │ DB (rw) │       │ DB (rw) │ ← Misma DB, ambos leen/escriben
         └─────────┘       └─────────┘

Deploy: Green recibe nueva versión → Smoke tests → Switch LB → Blue se vuelve idle
```

### Proceso Paso a Paso

| Paso | Acción | Automatizado | Duración |
| :---: | :--- | :---: | :--- |
| 1 | Desplegar nueva versión en entorno **idle** (Green) | ✅ CI/CD | ~3 min |
| 2 | Ejecutar smoke tests en Green (healthcheck + E2E críticos) | ✅ CI/CD | ~2 min |
| 3 | Si smoke tests pasan → Switch Load Balancer a Green | ✅ CI/CD | ~10 seg |
| 4 | Monitorear métricas por 5 minutos (error rate, latencia) | ⚠️ Semi-auto | 5 min |
| 5 | Si métricas OK → Blue se convierte en idle | ✅ Auto | Inmediato |
| 6 | Si métricas KO → Rollback: Switch LB de vuelta a Blue | ✅ Auto | ~10 seg |

## 2. Rollback Automático

### Triggers de Rollback

| Condición | Threshold | Acción |
| :--- | :--- | :--- |
| Error rate (5xx) | > 5% por 2 minutos | Rollback automático |
| Latencia P99 | > 5 segundos por 3 minutos | Rollback automático |
| Healthcheck fallido | 3 checks consecutivos | Rollback automático |
| Alerta manual | DevOps ejecuta rollback | Rollback manual (< 30 seg) |

### Comando de Rollback Manual
```text
# Revertir al entorno anterior
harness rollback --pipeline golden-pipeline --to-previous

# O directamente via kubectl/infraestructura
kubectl rollout undo deployment/api-gateway -n production
```

## 3. Canary Releases (Para cambios de alto riesgo)

Para cambios que afectan flujos críticos (auth, pagos), usamos canary:

```text
Fase 1: 5% del tráfico → nueva versión  (10 min, monitorear)
Fase 2: 25% del tráfico → nueva versión (15 min, monitorear)
Fase 3: 50% del tráfico → nueva versión (15 min, monitorear)
Fase 4: 100% del tráfico → rollout completo
```

> [!WARNING]
> Canary es **más lento** que Blue/Green (~1 hora vs ~10 min). Solo usar para cambios que afectan: autenticación, pagos, migraciones de datos.

## 4. Feature Flags

Para funcionalidades que se despliegan pero no se activan inmediatamente:

| Flag | Tipo | Default | Descripción |
| :--- | :--- | :---: | :--- |
| `auth.oauth_google` | Boolean | `false` | Habilita login con Google OAuth |
| `auth.oauth_apple` | Boolean | `false` | Habilita login con Apple Sign-In |
| `auth.mfa_enabled` | Boolean | `false` | Habilita MFA para roles admin |
| `auth.biometric_login` | Boolean | `false` | Login biométrico (V3, out of scope) |

### Evaluación de Feature Flags
```text
// En código
if (featureFlags.isEnabled('auth.oauth_google', { user, tenant })) {
  renderGoogleLoginButton();
}
```

## 5. Ventanas de Despliegue

| Entorno | Ventana Permitida | Aprobación Requerida |
| :--- | :--- | :--- |
| Development | 24/7 (auto-deploy on push) | Ninguna |
| Staging | 24/7 (auto-deploy on merge to `main`) | Pipeline verde |
| Production | Lunes-Jueves, 10:00-16:00 (hora local) | Pipeline verde + 1 approval manual |

> [!CAUTION]
> **No** desplegar a producción los viernes, fines de semana, ni después de las 16:00. Los incidentes fuera de horario son más difíciles de resolver.

## 6. Checklist Pre-Deploy a Producción

- [ ] Pipeline CI/CD completamente verde (validate → build → test → security)
- [ ] Smoke tests pasan en staging
- [ ] No hay alertas activas en producción
- [ ] Si hay migración de DB → probada en staging + backup verificado
- [ ] Runbook actualizado si el deploy cambia operaciones
- [ ] Comunicación al equipo en `#deploys`
