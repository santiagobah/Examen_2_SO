# ---- Etapa de Construcción ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copiar archivos de módulos y descargar dependencias primero para aprovechar el caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación
# CGO_ENABLED=0 para compilación estática sin CGO
# -ldflags="-w -s" para reducir el tamaño del binario (opcional)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go-api-server .

# ---- Etapa Final ----
FROM alpine:latest

WORKDIR /app

# Copiar el binario compilado de la etapa de construcción
COPY --from=builder /go-api-server .

# Copiar cualquier otro archivo necesario (ej. templates, archivos de configuración)
# COPY --from=builder /app/templates ./templates

# Exponer el puerto que usa la API Go
EXPOSE 8080

# Comando para ejecutar la aplicación
# Usamos /app/go-api-server porque ese es el nombre que le dimos en la etapa de build
CMD ["/app/go-api-server"]
