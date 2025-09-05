package nodestate

import (
	pb "github.com/Lucas-Sabbatini/TrabalhoFinalSD/pkg/kvstore"
)

type ComparisonResult int

const (
	Concorrente ComparisonResult = iota
	Vc1Antes
	Vc2Antes
)

func vcToMap(vc *pb.VectorClock) map[string]uint64 {
	m := make(map[string]uint64)
	if vc == nil {
		return m
	}
	for _, entry := range vc.GetEntries() {
		m[entry.GetNodeId()] = entry.GetCounter()
	}
	return m
}

// MergeVectorClocks mescla dois Vector Clocks, resultando em um novo
// Vector Clock com o valor máximo para cada NodeId.
func MergeVectorClocks(vc1, vc2 *pb.VectorClock) *pb.VectorClock {
	map1 := vcToMap(vc1)
	map2 := vcToMap(vc2)

	mergedMap := make(map[string]uint64)
	for k, v := range map1 {
		mergedMap[k] = v
	}

	for nodeId, counter2 := range map2 {
		counter1, exists := mergedMap[nodeId]
		if !exists || counter2 > counter1 {
			mergedMap[nodeId] = counter2
		}
	}

	entries := make([]*pb.VectorClockEntry, 0, len(mergedMap))
	for nodeId, counter := range mergedMap {
		entries = append(entries, &pb.VectorClockEntry{
			NodeId:  nodeId,
			Counter: counter,
		})
	}

	return &pb.VectorClock{Entries: entries}
}

// CompareVectorClocks compara dois Vector Clocks para determinar sua ordem causal.
func CompareVectorClocks(vc1, vc2 *pb.VectorClock) ComparisonResult {
	map1 := vcToMap(vc1)
	map2 := vcToMap(vc2)

	vc1HasGreater := false
	vc2HasGreater := false

	allKeys := make(map[string]struct{})
	for k := range map1 {
		allKeys[k] = struct{}{}
	}
	for k := range map2 {
		allKeys[k] = struct{}{}
	}

	for key := range allKeys {
		c1 := map1[key]
		c2 := map2[key]

		if c1 > c2 {
			vc1HasGreater = true
		}
		if c2 > c1 {
			vc2HasGreater = true
		}
	}

	if vc1HasGreater && !vc2HasGreater {
		return Vc2Antes
	}
	if !vc1HasGreater && vc2HasGreater {
		return Vc1Antes
	}
	if vc1HasGreater && vc2HasGreater {
		return Concorrente
	}

	return Concorrente
}
