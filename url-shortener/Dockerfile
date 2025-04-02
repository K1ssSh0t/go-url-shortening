# Etapa de construcción
FROM golang:1.23.1-alpine AS builder

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de módulos Go
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Construir la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/main.go

# Etapa final
FROM alpine:latest

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario construido
COPY --from=builder /app/url-shortener .

# Exponer puerto
EXPOSE 8080

# Agregar variable de entorno para configurar la conexión a MongoDB
ENV MONGO_URI="mongodb://mongo:27017"

# Comando para ejecutar
CMD ["./url-shortener"]