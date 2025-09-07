package src

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	topic  string
	qos    byte
	client mqtt.Client
}

// Configuração, conexão e retorna uma instância de um cliente de um broker MQTT.
func NewMQTTClient(clientID string) (*MQTTClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		fmt.Printf("Conexão perdida com o broker: %v\n", err)
	})

	client := mqtt.NewClient(opts)
	fmt.Printf("Conectando ao broker MQTT")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("falha ao conectar ao broker: %w", token.Error())
	}

	fmt.Println("Conectado ao broker MQTT com sucesso.")
	return &MQTTClient{client: client, topic: "kvstore/replication", qos: byte(1)}, nil
}

// Publica uma mensagem em um tópico
func (m *MQTTClient) Publish(payload string) {
	token := m.client.Publish(m.topic, m.qos, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Erro ao publicar mensagem no tópico '%s': %v\n", m.topic, token.Error())
	}
}

// Se inscreve em um tópico para receber mensagens.
// O 'callback' é a função que será executada quando uma mensagem chegar.
func (m *MQTTClient) Subscribe(callback mqtt.MessageHandler) {
	if token := m.client.Subscribe(m.topic, m.qos, callback); token.Wait() && token.Error() != nil {
		log.Fatalf("Falha ao se inscrever no tópico '%s': %v", m.topic, token.Error())
	}
	fmt.Printf("Inscrito com sucesso no tópico: %s\n", m.topic)
}

// Encerra a conexão com o broker de forma limpa.
func (m *MQTTClient) Disconnect() {
	fmt.Println("Desconectando do broker MQTT...")
	m.client.Disconnect(250) // Espera 250ms para completar
	fmt.Println("Desconectado.")
}
