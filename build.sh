#!/bin/bash
set -e

# Limpia todas las imágenes no utilizadas
echo "Eliminando imágenes Docker no utilizadas..."
docker image prune -a -f

# Detener contenedores en ejecución basados en url-shortener:latest
containers=$(docker ps -q --filter ancestor=url-shortener:latest)
if [ -n "$containers" ]; then
    echo "Deteniendo contenedores en ejecución de url-shortener:latest..."
    docker stop $containers
fi

# Construye la imagen Docker sin cache
echo "Realizando un build limpio..."
docker build --no-cache -t url-shortener:latest .

echo "Build completado con éxito."
