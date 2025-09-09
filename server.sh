#!/bin/bash

# Caminho para o executável do servidor
SERVER_BINARY="server/bin/server"

# Verifica se o executável existe
if [ ! -f "$SERVER_BINARY" ]; then
    echo "ERRO: O executável do servidor não foi encontrado em '$SERVER_BINARY'."
    echo "Por favor, execute o script './compile.sh' primeiro."
    exit 1
fi

if ! command -v docker-compose &> /dev/null
then
    echo "ERRO: O comando 'docker-compose' não foi encontrado."
    echo "Por favor, instale o Docker e o Docker Compose para continuar."
    exit 1
fi

echo "--- Iniciando o broker MQTT via Docker Compose ---"

# Sobe o docker container que estará rodando o Broker MQTT
docker compose up -d

# Executa o servidor, passando todos os argumentos recebidos pelo script
# Exemplo de uso: ./server.sh -porta 50051 -node-id "node-A"
echo "Iniciando o servidor com os seguintes argumentos: $@"
"$SERVER_BINARY" "$@"