package src

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
	"github.com/google/uuid"
)

type StoreEntry struct {
	Key      string
	Versions []*pb.Version
}

type NodeState struct {
	Node_id string
	Store   map[string]StoreEntry
}

func NewNodeState() *NodeState {
	return &NodeState{
		Node_id: uuid.New().String(),
		Store:   make(map[string]StoreEntry),
	}

}

// SerializeStoreEntry converte uma struct StoreEntry para uma representação JSON.
func SerializeStoreEntry(v *StoreEntry) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar StoreEntry para JSON: %w", err)
	}
	return string(jsonData), nil
}

func SerializeVersion(v *pb.Version) (string, error) {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao serializar Version para JSON: %w", err)
	}
	return string(jsonData), nil
}

func DeserializeVersion(jsonData string) (*pb.Version, error) {
	var v pb.Version
	err := json.Unmarshal([]byte(jsonData), &v)
	if err != nil {
		return nil, fmt.Errorf("erro ao desserializar JSON para Version: %w", err)
	}
	return &v, nil
}

func (node_state *NodeState) Process_put(key string, new_version_value string, is_replication_source bool, mqttClient *MQTTClient) {
	var incoming_version *pb.Version
	var current_vector_clock *pb.VectorClock
	current_store_entry := node_state.Process_get(key)
	current_timestamp := uint64(time.Now().UnixNano())

	if len(current_store_entry.Versions) == 0 {
		current_vector_clock = CreateVectorClock(node_state.Node_id, current_timestamp)
	} else {
		current_vector_clock = MergeStoreEntry(current_store_entry)
	}

	if !is_replication_source {
		incoming_version = &pb.Version{
			Timestamp:    current_timestamp,
			WriterNodeId: node_state.Node_id,
			VectorClock:  current_vector_clock,
			Value:        new_version_value,
		}
		found := false
		for _, entry := range incoming_version.VectorClock.Entries {
			if entry.NodeId == node_state.Node_id {
				entry.Counter++
				found = true
				break
			}
		}
		if !found {
			incoming_version.VectorClock.Entries = append(incoming_version.VectorClock.Entries, &pb.VectorClockEntry{
				NodeId:  node_state.Node_id,
				Counter: 1,
			})
		}
	} else {
		var err error
		incoming_version, err = DeserializeVersion(new_version_value)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	new_versions := []*pb.Version{incoming_version}

	for _, current_version := range current_store_entry.Versions {
		resultado := CompareVectorClocks(incoming_version.VectorClock, current_version.VectorClock)

		if resultado == Concorrente {
			new_versions = append(new_versions, current_version)
		}
	}

	current_store_entry.Versions = new_versions
	node_state.Store[key] = current_store_entry

	if !is_replication_source {
		go func() {
			serialize, err := SerializeStoreEntry(&current_store_entry)
			if err != nil {
				log.Fatalf("%v", err)
			}
			mqttClient.Publish(serialize)
		}()
	}
}

func (node_state *NodeState) Process_get(key string) StoreEntry {
	entry, ok := node_state.Store[key]
	if ok {
		return entry
	}

	return StoreEntry{Key: key, Versions: []*pb.Version{}}
}
