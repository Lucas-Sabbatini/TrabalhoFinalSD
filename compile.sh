#!/bin/bash

# Mensagem de início
echo "--- Iniciando a compilação do projeto ---"

# Garante que todas as dependências estão baixadas e o go.mod está limpo
echo "1. Sincronizando dependências..."
go mod tidy

# Compila o executável do servidor
echo "2. Compilando o servidor..."
go build -o server/bin/server ./server
if [ $? -ne 0 ]; then
    echo "ERRO: Falha ao compilar o servidor."
    exit 1
fi

# Compila o executável do cliente
echo "3. Compilando o cliente de teste..."
go build -o client-test/bin/client ./client-test
if [ $? -ne 0 ]; then
    echo "ERRO: Falha ao compilar o cliente."
    exit 1
fi

echo "--- Compilação concluída com sucesso! ---"
echo "Executáveis criados em 'server/bin/server' e 'client-test/bin/client'."