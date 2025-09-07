package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore" // ajuste para o caminho real do seu módulo

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // Usado para a conexão gRPC
)

// printUsage exibe a ajuda principal e sai.
func printUsage() {
	fmt.Fprintf(os.Stderr, "Uso: %s <comando> [argumentos]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Comandos disponíveis: get, put")
	fmt.Fprintln(os.Stderr, "Use \"<comando> -h\" para mais ajuda sobre um comando específico.")
	os.Exit(1)
}

func main() {
	// --- Configuração dos Subcomandos ---
	// Cria um FlagSet para o comando 'get'
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getAddr := getCmd.String("addr", "localhost:50051", "Endereço do servidor (host:porta)")
	getKey := getCmd.String("key", "", "A chave a ser buscada (obrigatório)")

	// Cria um FlagSet para o comando 'put'
	putCmd := flag.NewFlagSet("put", flag.ExitOnError)
	putAddr := putCmd.String("addr", "localhost:50051", "Endereço do servidor (host:porta)")
	putKey := putCmd.String("key", "", "A chave a ser inserida (obrigatório)")
	putValue := putCmd.String("value", "", "O valor a ser inserido (obrigatório)")

	// Verifica se um subcomando foi fornecido
	if len(os.Args) < 2 {
		printUsage()
	}

	// --- Lógica de Roteamento de Subcomandos ---
	var addr, key, value string
	var err error

	switch os.Args[1] {
	case "get":
		getCmd.Parse(os.Args[2:])
		if *getKey == "" {
			fmt.Fprintln(os.Stderr, "Erro: a flag -key é obrigatória para o comando get.")
			getCmd.Usage()
			os.Exit(1)
		}
		addr = *getAddr
		key = *getKey

	case "put":
		putCmd.Parse(os.Args[2:])
		if *putKey == "" || *putValue == "" {
			fmt.Fprintln(os.Stderr, "Erro: as flags -key e -value são obrigatórias para o comando put.")
			putCmd.Usage()
			os.Exit(1)
		}
		addr = *putAddr
		key = *putKey
		value = *putValue

	default:
		fmt.Fprintf(os.Stderr, "Comando desconhecido: %s\n", os.Args[1])
		printUsage()
	}

	// --- Conexão e Execução gRPC ---
	// Conecta ao servidor usando o endereço fornecido
	// Nota: grpc.WithInsecure() está obsoleto. Usamos a nova forma.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Falha ao conectar no servidor em %s: %v", addr, err)
	}
	defer conn.Close()

	client := pb.NewKvStoreClient(conn)

	// Contexto com timeout para as chamadas
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Executa a ação baseada no comando
	switch os.Args[1] {
	case "get":
		fmt.Printf("Buscando a chave '%s' no servidor %s...\n", key, addr)
		getResp, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("Erro na chamada Get: %v", err)
		}
		if len(getResp.Versions) == 0 {
			fmt.Println("Nenhum valor encontrado para esta chave.")
		} else {
			fmt.Println("Resposta Get:")
			for i, v := range getResp.Versions {
				fmt.Printf("  Versão %d: Valor=%s, Nó Escritor=%s, Timestamp=%d\n", i+1, v.Value, v.WriterNodeId, v.Timestamp)
			}
		}

	case "put":
		fmt.Printf("Enviando Put para o servidor %s: Chave='%s', Valor='%s'\n", addr, key, value)
		putResp, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("Erro na chamada Put: %v", err)
		}
		fmt.Println("Resposta Put:", putResp.Success)
	}
}
