package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
	src "github.com/Lucas-Sabbatini/TrabalhoFinalSD/server/src"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedKvStoreServer
	nodeState  *src.NodeState
	mqttClient *src.MQTTClient
}

// Implementação do método Put (gRPC)
func (s *Server) Put(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	fmt.Printf("[gRPC PUT] Recebido: key=%s value=%s\n", req.GetKey(), req.GetValue())

	// Chama a lógica de processamento principal.
	// `is_replication_source` é `false` porque esta é uma requisição original de um cliente gRPC.
	s.nodeState.Process_put(req.GetKey(), req.GetValue(), false, s.mqttClient)

	return &pb.PutResponse{
		Success: true,
	}, nil
}

// Implementação do método Get (gRPC)
func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	fmt.Printf("[gRPC GET] Recebido: key=%s\n", req.GetKey())

	// Usa a função process_get para encontrar a entrada no store.
	storeEntry := s.nodeState.Process_get(req.GetKey())

	// Retorna todas as versões ativas (concorrentes) para a chave.
	return &pb.GetResponse{
		Versions: storeEntry.Versions,
	}, nil
}

func main() {
	portPtr := flag.Int("porta", 50051, "Porta em que o servidor web irá ouvir as conexões")
	flag.Parse()
	addr := fmt.Sprintf(":%d", *portPtr)

	nodeState := src.NewNodeState()
	fmt.Printf("Nó iniciado com ID: %s\n", nodeState.Node_id)

	mqttClient, err := src.NewMQTTClient(nodeState.Node_id)
	if err != nil {
		log.Fatalf("Falha ao inicializar o cliente MQTT: %v", err)
	}
	defer mqttClient.Disconnect()

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("[MQTT] Mensagem recebida no tópico %s\n", msg.Topic())

		var repMsg src.StoreEntry
		if err := json.Unmarshal(msg.Payload(), &repMsg); err != nil {
			log.Printf("Erro ao desserializar mensagem de replicação: %v \n", err)
			return
		}

		var achou bool
		for _, version := range repMsg.Versions {
			achou = false
			for _, stored_version := range nodeState.Store[repMsg.Key].Versions {
				if version.Timestamp == stored_version.Timestamp {
					achou = true
				}
			}
			if !achou {
				serialized_version, err := src.SerializeVersion(version)
				if err != nil {
					fmt.Printf("Erro ao serializar uma versão: %v", err)
				}
				nodeState.Process_put(repMsg.Key, serialized_version, true, mqttClient)
			}
		}
	}

	mqttClient.Subscribe(messageHandler)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("falha ao escutar: %v", err)
	}

	grpcServer := grpc.NewServer()

	srv := &Server{
		nodeState:  nodeState,
		mqttClient: mqttClient,
	}
	pb.RegisterKvStoreServer(grpcServer, srv)

	go func() {
		fmt.Printf("Servidor gRPC ouvindo em %s", addr)
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
