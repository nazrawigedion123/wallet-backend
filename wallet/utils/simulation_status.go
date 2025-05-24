package utils

import (
	"sync"
	"time"
)

type SimulationStatus struct {
	Entity     string    `json:"entity"`
	StartedAt  time.Time `json:"started_at"`
	Count      int       `json:"count"`
	Completed  int       `json:"completed"`
	Running    bool      `json:"running"`
}

var (
	statuses = make(map[string]*SimulationStatus)
	mu       sync.Mutex
)

func StartSimulation(entity string, total int) {
	mu.Lock()
	defer mu.Unlock()

	statuses[entity] = &SimulationStatus{
		Entity:    entity,
		StartedAt: time.Now(),
		Count:     total,
		Completed: 0,
		Running:   true,
	}
}

func IncrementProgress(entity string) {
	mu.Lock()
	defer mu.Unlock()

	if stat, exists := statuses[entity]; exists {
		stat.Completed++
	}
}

func EndSimulation(entity string) {
	mu.Lock()
	defer mu.Unlock()

	if stat, exists := statuses[entity]; exists {
		stat.Running = false
	}
}

func GetStatuses() []*SimulationStatus {
	mu.Lock()
	defer mu.Unlock()

	var list []*SimulationStatus
	for _, s := range statuses {
		list = append(list, s)
	}
	return list
}
