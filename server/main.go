package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"

	"google.golang.org/grpc"
)

// Server implementa a interface KvStoreServer gerada em kv_store_grpc.pb.go
type Server struct {
	pb.UnimplementedKvStoreServer
	// aqui você pode manter seu estado, ex: map[string]string ou NodeState
	store map[string]string
}

// Implementação do método Put (gRPC)
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	// exemplo simples: guarda no map
	s.store[req.Key] = req.Value

	fmt.Printf("[PUT] key=%s value=%s\n", req.Key, req.Value)

	return &pb.PutResponse{
		Success:      true,
		ErrorMessage: "",
	}, nil
}

// Implementação do método Get (gRPC)
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	val, ok := s.store[req.Key]
	if !ok {
		return &pb.GetResponse{
			Versions:     []*pb.Version{}, // vazio se não encontrou
			ErrorMessage: "key not found",
		}, nil
	}

	// aqui simplificamos: só retornamos uma versão sem VectorClock
	version := &pb.Version{
		Value:        val,
		VectorClock:  &pb.VectorClock{Entries: []*pb.VectorClockEntry{}},
		Timestamp:    0,
		WriterNodeId: "node_1",
	}

	return &pb.GetResponse{
		Versions: []*pb.Version{version},
	}, nil
}

func main() {
	// cria listener TCP
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := &Server{store: make(map[string]string)}

	// registra o serviço no gRPC
	pb.RegisterKvStoreServer(grpcServer, srv)

	fmt.Println("Servidor gRPC ouvindo em :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("falha ao iniciar servidor: %v", err)
	}
}
