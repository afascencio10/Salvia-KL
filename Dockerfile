# Etapa 1: Build de la aplicación
FROM golang:1.22-alpine AS builder

# Instalar dependencias necesarias para compilar
RUN apk add --no-cache git ca-certificates tzdata

# Definir el directorio de trabajo donde ocurrirá la compilación
WORKDIR /app

# Copiar todo el contenido de la carpeta src (que incluye main.go, módulos locales y recursos a embeber)
COPY src/ .

# Descargar módulos (el go.work lo resolverá automáticamente)
RUN go mod download

# Compilar estáticamente el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -installsuffix cgo -o salvia-app main.go

# Etapa 2: Imagen final ultra ligera
FROM alpine:latest

# Instalar los certificados y zona horaria (útil para la BD y timestamps)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copiar únicamente el binario generado en la capa anterior
COPY --from=builder /app/salvia-app .

# Copiar los certificados de prueba por si el código hace comprobación local
# (Aunque con el uso de PORT, RunTLS se ignora en producción)
COPY src/certs/ /app/certs/

# Definir el puerto explícitamente (Render usa uno dinámico y pisa este)
ENV PORT=10000

# Exponer el puerto
EXPOSE ${PORT}

# Comando de ejecución
CMD ["./salvia-app"]
