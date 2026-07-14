---
okf_version: "1.0"
type: Coding_Standard
title: "Code Review Checklist"
description: "Lista de verificación estructurada que todo reviewer debe seguir al revisar un Pull Request."
owner: "@platform-team"
status: active
tags: [standards, review, quality, pr]
dependencies:
  - DoD: 5-standards/definition-of-done.md
  - Coding: 5-standards/coding-standards.md
last_updated: "2026-06-23"
---

# Code Review Checklist

Este checklist debe usarse como template en cada PR. El reviewer copia la sección relevante y marca los items en el comentario de aprobación.

## 1. Spec Compliance (Spec-Driven)

- [ ] ¿El PR tiene un documento OKF asociado que describe el cambio?
- [ ] ¿El código implementa exactamente lo que el spec describe? (ni más, ni menos)
- [ ] ¿Las `dependencies:` del spec fueron actualizadas si se agregaron nuevas relaciones?
- [ ] Si es un endpoint, ¿el contrato de Request/Response coincide con el spec?

## 2. Correctness (Funcionalidad)

- [ ] ¿El código resuelve el problema descrito en el issue/spec?
- [ ] ¿Se manejan todos los edge cases mencionados en el spec?
- [ ] ¿Los errores se propagan correctamente (no se silencian)?
- [ ] ¿La lógica de negocio está en la capa de servicio (no en el controller/componente)?

## 3. Security

- [ ] ¿No hay secretos hardcodeados (API keys, passwords, tokens)?
- [ ] ¿Los inputs del usuario están validados/sanitizados?
- [ ] ¿Los endpoints respetan la [RBAC matrix](/.knowledge/6-security/rbac-matrix.md)?
- [ ] ¿No se expone información sensible en logs o mensajes de error?
- [ ] ¿Se usa parametrized queries (no string concatenation) para DB?
- [ ] Si maneja autenticación, ¿sigue el ADR de autenticación del proyecto (ver `1-architecture/adrs/`)?

## 4. Performance

- [ ] ¿No hay queries N+1 (consultas dentro de loops)?
- [ ] ¿Las consultas a DB tienen los índices necesarios?
- [ ] ¿Se usa paginación para listas potencialmente largas?
- [ ] ¿No hay operaciones blocking en el event loop (Node.js)?
- [ ] ¿Los assets estáticos (imágenes, fonts) están optimizados?

## 5. Testing

- [ ] ¿Los tests cubren el happy path Y los edge cases?
- [ ] ¿Los tests son independientes entre sí (no comparten estado)?
- [ ] ¿Los nombres de tests son descriptivos (`should X when Y`)?
- [ ] ¿La cobertura del PR no reduce la cobertura global?
- [ ] ¿Se usan factories/mocks en vez de datos hardcodeados?

## 6. Code Quality

- [ ] ¿El naming sigue los [coding-standards](/.knowledge/5-standards/coding-standards.md)?
- [ ] ¿No hay código duplicado que debería extraerse a un util/helper?
- [ ] ¿Los imports siguen el orden estándar?
- [ ] ¿No hay `console.log` residual (debe usar `logger`)?
- [ ] ¿No hay código comentado (dead code)?
- [ ] ¿Las funciones tienen menos de ~50 líneas?

## 7. Documentation

- [ ] ¿Las funciones públicas tienen JSDoc/docstring?
- [ ] ¿Los TODOs tienen issue vinculado (`// TODO(#123):`)?
- [ ] ¿Los breaking changes están documentados con `> [!WARNING]` en el spec?

## 8. Formato de Aprobación

Al aprobar un PR, el reviewer debe incluir el siguiente bloque:

```text
## ✅ Review Approved

### Checklist Verified:
- [x] Spec Compliance
- [x] Correctness
- [x] Security
- [x] Performance
- [x] Testing
- [x] Code Quality
- [x] Documentation

### Notes:
[Comentarios opcionales del reviewer]
```

> [!IMPORTANT]
> Un PR **no puede ser mergeado** si alguna sección del checklist tiene items sin marcar y sin justificación de `N/A`.
