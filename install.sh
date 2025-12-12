#!/bin/bash

set -e

# Colores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Instalando orgmcron...${NC}"

# Detectar arquitectura y OS
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Arquitectura no soportada: $ARCH${NC}"
        exit 1
        ;;
esac

# Obtener la última versión desde GitHub
echo -e "${YELLOW}Obteniendo última versión...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/osmargm1202/orgmcron/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "latest")

if [ -z "$LATEST_VERSION" ] || [ "$LATEST_VERSION" = "latest" ]; then
    echo -e "${YELLOW}No se pudo obtener la versión, usando 'latest'${NC}"
    LATEST_VERSION="latest"
fi

# URL del binario
BINARY_URL="https://github.com/osmargm1202/orgmcron/releases/${LATEST_VERSION}/download/orgmcron-${OS}-${ARCH}"

# Directorio de destino
INSTALL_DIR="$HOME/.local/bin"
BINARY_PATH="$INSTALL_DIR/orgmcron"

# Crear directorio si no existe
mkdir -p "$INSTALL_DIR"

# Descargar binario
echo -e "${YELLOW}Descargando binario desde GitHub...${NC}"
if ! curl -L -f -o "$BINARY_PATH" "$BINARY_URL"; then
    echo -e "${RED}Error: No se pudo descargar el binario desde $BINARY_URL${NC}"
    echo -e "${YELLOW}Intentando construir desde el código fuente...${NC}"
    
    # Si falla la descarga, intentar construir desde el código
    if command -v go &> /dev/null; then
        echo -e "${YELLOW}Compilando desde el código fuente...${NC}"
        go build -o "$BINARY_PATH" .
    else
        echo -e "${RED}Error: Go no está instalado y no se pudo descargar el binario${NC}"
        exit 1
    fi
fi

# Hacer ejecutable
chmod +x "$BINARY_PATH"

# Verificar que está en el PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}Advertencia: $INSTALL_DIR no está en tu PATH${NC}"
    echo -e "${YELLOW}Agrega esta línea a tu ~/.bashrc, ~/.zshrc o ~/.config/fish/config.fish:${NC}"
    echo -e "${GREEN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
fi

echo -e "${GREEN}✓ orgmcron instalado exitosamente en $BINARY_PATH${NC}"
echo ""
echo -e "${GREEN}Próximos pasos:${NC}"
echo -e "  1. Configura la pingkey: ${YELLOW}orgmcron config pingkey <tu-key>${NC}"
echo -e "  2. Agrega un job: ${YELLOW}orgmcron add${NC}"
echo -e "  3. Instala el servicio: ${YELLOW}orgmcron install${NC}"
echo -e "  4. Inicia el servicio: ${YELLOW}systemctl --user start orgmcron${NC}"

