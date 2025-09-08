#!/bin/bash

CLIENTE="./client-test/bin/client"

SERVIDOR1="127.0.0.1:50052"
SERVIDOR2="127.0.0.1:50051"

CHAVE="Paulo-Coelho"
VALOR1="Python"
VALOR2="Golang"

echo "Iniciando teste de concorrência..."
echo "------------------------------------"

echo "Enviando para $SERVIDOR1: Chave=$CHAVE, Valor=$VALOR1"
$CLIENTE put -addr $SERVIDOR1 -key $CHAVE -value $VALOR1 &
PID1=$! 

echo "Enviando para $SERVIDOR2: Chave=$CHAVE, Valor=$VALOR2"
$CLIENTE put -addr $SERVIDOR2 -key $CHAVE -value $VALOR2 &
PID2=$! 

echo "Aguardando a conclusão dos PUTs (PID: $PID1, $PID2)..."
wait $PID1
wait $PID2

echo "------------------------------------"
echo "PUTs concluídos. Realizando GET para verificar as versões..."

sleep 1


$CLIENTE get -addr $SERVIDOR1 -key $CHAVE

echo "------------------------------------"
echo "Teste finalizado."