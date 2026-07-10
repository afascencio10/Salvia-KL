---
okf_version: "1.0"
type: SLA
title: "Service Level Objectives (SLOs)"
description: "Definición de SLIs, SLOs y SLAs para los servicios críticos de la plataforma con umbrales numéricos y acciones de escalación."
owner: "@sre-squad"
status: active
tags: [ops, sla, slo, sli, monitoring, reliability]
last_updated: "2026-06-23"
---

# Service Level Objectives (SLOs)

## Definiciones

| Concepto | Definición |
| :--- | :--- |
| **SLI** (Indicator) | Métrica medible del servicio |
| **SLO** (Objective) | Target numérico del SLI (interno) |
| **SLA** (Agreement) | Compromiso contractual con clientes (externo) |
| **Error Budget** | % de indisponibilidad permitida antes de pausar deploys |

## SLOs por Servicio

### Auth Service (Crítico)

| SLI | Métrica | SLO | SLA | Medición |
| :--- | :--- | :---: | :---: | :--- |
| Disponibilidad | % de requests exitosos (non-5xx) | 99.9% | 99.5% | Rolling 30 días |
| Latencia Login (P50) | Percentil 50 de `POST /login` | < 500ms | < 2s | Rolling 7 días |
| Latencia Login (P99) | Percentil 99 de `POST /login` | < 2s | < 5s | Rolling 7 días |
| Latencia Verify (P99) | Percentil 99 de `GET /verify` | < 50ms | < 200ms | Rolling 7 días |
| Error Rate | % de responses 5xx | < 0.1% | < 0.5% | Rolling 24 horas |

### API Gateway

| SLI | Métrica | SLO | SLA | Medición |
| :--- | :--- | :---: | :---: | :--- |
| Disponibilidad | Uptime del gateway | 99.95% | 99.9% | Rolling 30 días |
| Latencia (P99) | Tiempo de routing + middleware | < 100ms | < 500ms | Rolling 7 días |
| Throughput | Requests por segundo sostenidos | ≥ 200 RPS | ≥ 100 RPS | Peak hour |

### Base de Datos (Primaria)

| SLI | Métrica | SLO | Medición |
| :--- | :--- | :---: | :--- |
| Disponibilidad | Uptime del primary | 99.99% | Rolling 30 días |
| Latencia queries (P99) | Tiempo de ejecución de queries | < 100ms | Rolling 7 días |
| Conexiones activas | % de max_connections usado | < 80% | Tiempo real |

## Error Budget

### Cálculo

```text
SLO de Disponibilidad: 99.9%
Error Budget mensual: 100% - 99.9% = 0.1%

En un mes de 30 días = 43,200 minutos
Error Budget = 43,200 × 0.001 = 43.2 minutos de downtime permitidos

Si consumimos el error budget → PAUSAR deployments a producción
```

### Políticas de Error Budget

| Error Budget Restante | Estado | Acción |
| :---: | :--- | :--- |
| > 50% | 🟢 Saludable | Deploy normal, permitir experimentación |
| 25-50% | 🟡 Precaución | Solo deploys con test completo, no experimentar |
| < 25% | 🔴 Crítico | Solo deploys de hotfix, zero-risk changes |
| 0% | ⛔ Agotado | FREEZE de deploys hasta próximo período |

## Alerting

### Thresholds de Alerta

| Severidad | Condición | Canal | Respuesta Esperada |
| :--- | :--- | :--- | :--- |
| 🔴 P1 (Critical) | Disponibilidad < 99% por 5 min | PagerDuty + Slack `#incidents` | < 15 min |
| 🟠 P2 (High) | Latencia P99 > 5s por 10 min | Slack `#alerts` | < 30 min |
| 🟡 P3 (Medium) | Error rate > 1% por 15 min | Slack `#monitoring` | < 2 horas |
| 🔵 P4 (Low) | Error budget < 50% | Email semanal | Siguiente sprint |

## Dashboards

| Dashboard | Herramienta | Contenido |
| :--- | :--- | :--- |
| SLO Overview | Datadog/Grafana | SLIs actuales vs targets, error budget |
| Auth Performance | Datadog APM | Latencia, throughput, error rate del Auth Service |
| Infrastructure | Grafana | CPU, memoria, disco, conexiones DB |
| Business Metrics | Custom | Login success rate, tiempo promedio de autenticación |

> [!IMPORTANT]
> Los SLOs se revisan **trimestralmente**. Si un SLO se incumple consistentemente, se ajusta el target o se prioriza la mejora de reliability.
