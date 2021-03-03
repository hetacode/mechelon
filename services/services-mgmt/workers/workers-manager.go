package smgworkers

import (
	"fmt"
	"log"

	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
)

// RefreshWorkerPeriod in seconds
const RefreshWorkerPeriod int64 = 15

// WorkersManager - create, remove workers
type WorkersManager struct {
	serviceStateRepository *smgeventstore.ServiceStateRepository
	workers                map[string]*InstanceActivityWorker
	disableWorkerChannel   chan string
}

// NewWorkersManager new instance
func NewWorkersManager(serviceStateRepository *smgeventstore.ServiceStateRepository) *WorkersManager {
	inst := &WorkersManager{
		workers:                make(map[string]*InstanceActivityWorker),
		serviceStateRepository: serviceStateRepository,
		disableWorkerChannel:   make(chan string),
	}
	go func() {
		for toDisable := range inst.disableWorkerChannel {
			if w, ok := inst.workers[toDisable]; ok {
				w.Stop()
				delete(inst.workers, toDisable)
				log.Printf("WorkersManager | worker '%s' has been deleted", toDisable)
			}
		}
	}()

	return inst
}

// CreateWorker and run it
func (wm *WorkersManager) CreateWorker(projectName, serviceName string) {
	key := fmt.Sprintf("%s-%s", projectName, serviceName)
	w := NewInstanceActivityWorker(wm.serviceStateRepository, projectName, serviceName, RefreshWorkerPeriod, wm.disableWorkerChannel)
	wm.workers[key] = w
	go w.Start()

	log.Printf("WorkersManager | worker '%s' is running", key)
}
