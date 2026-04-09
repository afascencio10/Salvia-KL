// Package docs — Placeholder generado manualmente.
// Regenerar con: swag init -g main.go --parseDependency --parseInternal --output docs
//
// Este archivo será SOBREESCRITO por swag init. No editar manualmente.
package docs

import "github.com/swaggo/swag"

func init() {
	swag.Register(swag.Name, &swag.Spec{
		Version:          "2.0",
		Host:             "localhost",
		BasePath:         "/api/v1",
		Schemes:          []string{"https"},
		Title:            "SALVIA API",
		Description:      "API del Sistema de Atención y Vigilancia de Violencias contra la mujer (SALVIA) — Fase 2.",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  `{"swagger":"2.0","info":{"title":"SALVIA API","description":"API SALVIA Fase 2","version":"2.0"},"host":"localhost","basePath":"/api/v1","paths":{}}`,
	})
}
