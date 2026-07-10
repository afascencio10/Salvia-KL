---
okf_version: "1.0"
type: Project
title: "Feature: Autenticación y Redirección de Sesión (Auth)"
description: "Flujo de autenticación, control de sesiones y redirección según el rol asignado al usuario."
owner: "@security-guild"
status: active
tags: [feature, auth, security]
dependencies:
  - Spec: 6-security/rbac-matrix.md
last_updated: "2026-07-02"
---

# Feature: Autenticación y Redirección de Sesión (Auth)

## Descripción

El sistema de autenticación de SALVIA permite a los usuarios acceder al aplicativo validando sus credenciales (login y contraseña encriptada con bcrypt) y verificando el captcha. Una vez autenticado exitosamente, el backend retorna un conjunto de reglas de navegación (`nav_rules`) y el menú asignado según el rol del usuario, permitiendo al frontend realizar la redirección correcta a su espacio de trabajo correspondiente.

## Documentos Relacionados

### Seguridad
- [Matriz RBAC](/.knowledge/6-security/rbac-matrix.md) — Matriz de control de acceso y catálogo de roles

### Backend — Endpoints
| Endpoint | Spec | Estado |
| :--- | :--- | :---: |
| `POST /seguridad/login` | [login-post.md](/.knowledge/3-features/Auth/backend/endpoints/login-post.md) | ✅ Activo |

## Historias de Usuario

| ID | Historia | Estado |
| :--- | :--- | :---: |
| US-AUTH-01 | Como usuario del sistema, quiero iniciar sesión de forma segura y ser redirigido automáticamente a mi pantalla de trabajo de acuerdo a mi rol. | ✅ Implementado |
