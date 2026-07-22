#!/bin/bash
# Levanta el backend Go en modo HTTP local (sin TLS) para pruebas de QA/Playwright.
cd "$(dirname "$0")/src"
export PORT=9090
exec go run -a .
