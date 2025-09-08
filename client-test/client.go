package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"google.golang.org/grpc"
)

var (
	broker   = "tcp://localhost:1883"
	clientID = "go-client-for-kvstore"
)

// subscribeToReplicationTopic configura o cliente para assinar o tópico de replicação.
func subscribeToReplicationTopic(mqttClient mqtt.Client) {
	topic := "kvstore/replication"
	qos := byte(1)

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("\n--- Mensagem de replicação recebida! ---\n")
		fmt.Printf("Tópico: %s\n", msg.Topic())
		fmt.Printf("Payload: %s\n", msg.Payload())

		var entry pb.StoreEntry
		if err := json.Unmarshal(msg.Payload(), &entry); err != nil {
			log.Printf("Erro ao desserializar a mensagem JSON: %v", err)
			return
		}

		fmt.Printf("Dados desserializados: Chave=%s, Valor=%s, Versão=%+v\n", entry.Key, entry.Value, entry.Version)
		fmt.Println("-------------------------------------------")
	}

	token := mqttClient.Subscribe(topic, qos, messageHandler)
	token.Wait()

	if token.Error() != nil {
		log.Fatalf("Falha ao assinar o tópico '%s': %v", topic, token.Error())
	}

	fmt.Printf("Assinado com sucesso o tópico '%s'.\n", topic)
}

func main() {
	// --- Conexão gRPC ---
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("falha ao conectar gRPC: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewKvStoreClient(conn)

	// --- Conexão MQTT ---
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("falha ao conectar ao broker MQTT: %v", token.Error())
	}
	fmt.Println("Conectado com sucesso ao broker MQTT!")

	// --- Inicia a assinatura em segundo plano ---
	go subscribeToReplicationTopic(mqttClient)

	// --- Realiza as operações PUT e GET ---
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Operação PUT
	key := "foo"
	value := "bar"
	fmt.Printf("\nChamando PUT com Chave='%s', Valor='%s'...\n", key, value)
	putResp, err := grpcClient.Put(ctx, &pb.PutRequest{Key: key, Value: value})
	if err != nil {
		log.Fatalf("Erro no Put: %v", err)
	}
	fmt.Println("Resposta Put:", putResp)

	// Publica a alteração via MQTT se o PUT foi bem-sucedido
	if putResp.Success {
		storeEntry := &pb.StoreEntry{
			Key:     key,
			Value:   value,
			Version: putResp.Version,
		}

		payload, err := json.Marshal(storeEntry)
		if err != nil {
			log.Fatalf("Falha ao serializar para JSON: %v", err)
		}

		topic := "kvstore/replication"
		qos := byte(1)
		token := mqttClient.Publish(topic, qos, false, payload)
		token.Wait()

		if token.Error() != nil {
			log.Fatalf("Falha ao publicar mensagem MQTT: %v", token.Error())
		}
		fmt.Printf("Publicado com sucesso no tópico '%s'!\n", topic)
	}

	// Operação GET
	fmt.Printf("\nChamando GET com Chave='%s'...\n", key)
	getResp, err := grpcClient.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("Erro no Get: %v", err)
	}
	fmt.Println("Resposta Get:")
	for _, v := range getResp.Versions {
		fmt.Printf("   Valor=%s, Node=%s\n", v.Value, v.WriterNodeId)
	}

	// Espera por um sinal de interrupção para fechar o programa
	fmt.Println("\nOperações completas. Aguardando mensagens de replicação... Pressione Ctrl+C para sair.")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c // Bloqueia até receber o sinal
	fmt.Println("Encerrando o cliente...")
}
