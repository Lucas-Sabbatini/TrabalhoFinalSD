package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore" // ajuste para o caminho real do seu módulo

	"google.golang.org/grpc"
)

func main() {
	// conecta ao servidor
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("falha ao conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewKvStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// faz um PUT
	putResp, err := client.Put(ctx, &pb.PutRequest{Key: "foo", Value: "bar"})
	if err != nil {
		log.Fatalf("erro no Put: %v", err)
	}
	fmt.Println("Resposta Put:", putResp)

	// faz um GET
	getResp, err := client.Get(ctx, &pb.GetRequest{Key: "foo"})
	if err != nil {
		log.Fatalf("erro no Get: %v", err)
	}
	fmt.Println("Resposta Get:")
	for _, v := range getResp.Versions {
		fmt.Printf("  Valor=%s, Node=%s\n", v.Value, v.WriterNodeId)
	}
}
