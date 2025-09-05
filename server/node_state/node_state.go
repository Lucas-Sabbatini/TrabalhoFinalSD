package nodestate

import (
	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
	"github.com/google/uuid"
)

type StoreEntry struct {
	key      string
	versions []pb.Version
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

func (node_state *NodeState) process_put(new_version_value string, is_replication_source bool) {

}

func (node_state *NodeState) process_get(key string) StoreEntry {
	entry, ok := node_state.Store[key]
	if ok {
		return entry
	}

	return StoreEntry{key: key}
}
