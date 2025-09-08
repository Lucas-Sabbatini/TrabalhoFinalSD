package src

import (
	"flag"
	"fmt"

	"github.com/google/uuid"
)

// RuntimeConfig representa a configuração obtida por linha de comando.
type RuntimeConfig struct {
	NodeID         string
	ListenAddr     string
	MQTTBrokerAddr string
	MQTTBrokerPort int
}

var (
	flagNodeID         = flag.String("node-id", "", "ID único do nó (ex.: node_A). Se vazio, um UUID será gerado.")
	flagListenAddr     = flag.String("listen-addr", "127.0.0.1:50051", "Endereço IP e porta onde o servidor gRPC irá escutar (ex.: 127.0.0.1:50051).")
	flagMQTTBrokerAddr = flag.String("mqtt-broker-addr", "127.0.0.1", "Endereço IP do broker MQTT (ex.: 127.0.0.1).")
	flagMQTTBrokerPort = flag.Int("mqtt-broker-port", 1883, "Porta do broker MQTT (ex.: 1883).")
)

// ParseRuntimeFlags faz o parse dos argumentos de linha de comando e retorna um objeto RuntimeConfig.
func ParseRuntimeFlags() RuntimeConfig {
	if !flag.Parsed() {
		flag.Parse()
	}

	id := *flagNodeID
	if id == "" {
		id = uuid.NewString() // gera automaticamente um UUID caso o usuário não forneça
	}

	return RuntimeConfig{
		NodeID:         id,
		ListenAddr:     *flagListenAddr,
		MQTTBrokerAddr: *flagMQTTBrokerAddr,
		MQTTBrokerPort: *flagMQTTBrokerPort,
	}
}

// MQTTBrokerURL monta a URL completa do broker MQTT no formato usado pela biblioteca Paho.
func (c RuntimeConfig) MQTTBrokerURL() string {
	return fmt.Sprintf("tcp://%s:%d", c.MQTTBrokerAddr, c.MQTTBrokerPort)
}
