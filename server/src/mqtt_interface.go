package src

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	ReplicationTopic = "kvstore/replication"
	DefaultQoS       = 1 
)

type MQTTClient struct {
	topic  string
	qos    byte
	client mqtt.Client
}

func NewMQTTClientWithBroker(nodeID string, brokerURL string) (*MQTTClient, error) {
    opts := mqtt.NewClientOptions().
        AddBroker(brokerURL).
        SetClientID(nodeID).
        SetCleanSession(true)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        return nil, token.Error()
    }

    return &MQTTClient{client: client}, nil
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
	if token := m.client.Subscribe(ReplicationTopic, DefaultQoS, callback); token.Wait() && token.Error() != nil {
		log.Fatalf("Falha ao se inscrever no tópico '%s': %v", ReplicationTopic, token.Error())
	}
	fmt.Printf("Inscrito com sucesso no tópico: %s\n", ReplicationTopic)
}

// Encerra a conexão com o broker de forma limpa.
func (m *MQTTClient) Disconnect() {
	fmt.Println("Desconectando do broker MQTT...")
	m.client.Disconnect(250) // Espera 250ms para completar
	fmt.Println("Desconectado.")
}
