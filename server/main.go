package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
	nodestate "github.com/Lucas-Sabbatini/TrabalhoFinalSD/server/node_state"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"google.golang.org/grpc"
)

// Server implementa a interface KvStoreServer gerada em kv_store_grpc.pb.go
type Server struct {
	pb.UnimplementedKvStoreServer
	// aqui você pode manter seu estado, ex: map[string]string ou NodeState
	store      map[string]string
	mqttClient *nodestate.MQTTClient
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

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Mensagem recebida no tópico %s: %s\n", msg.Topic(), string(msg.Payload()))
	// Aqui você pode adicionar a lógica para processar a mensagem recebida
}

func main() {
	nodeState := nodestate.NewNodeState()

	mqttClient, err := nodestate.NewMQTTClient(nodeState.Node_id)
	if err != nil {
		log.Fatalf("Falha ao inicializar o cliente MQTT: %v", err)
	}
	defer mqttClient.Disconnect()
	mqttClient.Subscribe(messageHandler)

	// Configura e inicia o servidor gRPC
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv := &Server{
		store:      make(map[string]string),
		mqttClient: mqttClient,
	}
	pb.RegisterKvStoreServer(grpcServer, srv)

	// Inicia o servidor gRPC em uma goroutine para não bloquear a main
	go func() {
		fmt.Println("Servidor gRPC ouvindo em :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("falha ao iniciar servidor gRPC: %v", err)
		}
	}()

	fmt.Println("Servidor em execução. Pressione Ctrl+C para sair.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\nEncerrando servidor...")
	grpcServer.GracefulStop()
	fmt.Println("Servidor gRPC encerrado.")
}
