package smgworkers

import (
	"fmt"
	"log"
	"time"

	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
)

const (
	// ValidActivePeriod a period while service instance can be active
	ValidActivePeriod int64 = 30
	// ValidIdlePeriod a period while service instance can be idle
	ValidIdlePeriod int64 = 10
)

// InstanceActivityWorker struct
// Worker which check instances of service linked by aggregator state key in every setting period of time
type InstanceActivityWorker struct {
	repository     *smgeventstore.ServiceStateRepository
	projectName    string
	serviceName    string
	period         int64 // seconds
	createdAt      int64
	updatedAt      int64
	isRunning      bool
	disableChannel chan<- string
}

// NewInstanceActivityWorker create new instance
func NewInstanceActivityWorker(repository *smgeventstore.ServiceStateRepository, projectName, serviceName string, period int64, disableChannel chan string) *InstanceActivityWorker {
	worker := &InstanceActivityWorker{
		repository:     repository,
		projectName:    projectName,
		serviceName:    serviceName,
		period:         period,
		isRunning:      true,
		createdAt:      time.Now().Unix(),
		disableChannel: disableChannel,
	}
	worker.updatedAt = worker.createdAt

	return worker
}

// Start worker job
func (w *InstanceActivityWorker) Start() {
	for {
		if !w.isRunning {
			return
		}

		time.Sleep(time.Second)
		if time.Now().Unix()-w.updatedAt < w.period {
			continue
		}
		w.updatedAt = time.Now().Unix()

		aggr := w.repository.GetAggregator(w.projectName, w.serviceName)
		if aggr == nil {
			log.Printf("InstanceActivityWorker | cannot find an aggregator for %s service for %s project", w.serviceName, w.projectName)

			w.disableChannel <- fmt.Sprintf("%s-%s", w.projectName, w.serviceName)
			continue
		}

		instancesCount := len(aggr.State.Instances)
		inActiveInstancesCount := 0
		for _, instance := range aggr.State.Instances {
			switch instance.State {
			case smgeventstore.Active:
				if time.Now().Unix()-instance.UpdatedAt > ValidActivePeriod {
					aggr.SetInstanceAsIdle(instance.Name)
					log.Printf("InstanceActivityWorker | '%s' instance of '%s' service for project '%s' is IDLE", instance.Name, w.serviceName, w.projectName)
				}
			case smgeventstore.Idle:
				if time.Now().Unix()-instance.UpdatedAt > ValidIdlePeriod {
					aggr.SetInstanceAsInactive(instance.Name)
					log.Printf("InstanceActivityWorker | '%s' instance of '%s' service for project '%s' is IN_ACTIVE", instance.Name, w.serviceName, w.projectName)
				}
			case smgeventstore.InActive:
				inActiveInstancesCount++
			}
		}

		// No active instances - worker should be stopped
		if instancesCount == inActiveInstancesCount {
			w.disableChannel <- fmt.Sprintf("%s-%s", w.projectName, w.serviceName)
			continue
		}

		events := aggr.GetPendingEvents()
		if err := w.repository.SaveEvents(w.projectName, w.serviceName, events); err != nil {
			log.Printf("InstanceActivityWorker | cannot save events for %s service for %s project", w.serviceName, w.projectName)
		}
	}
}

// Stop worker job
func (w *InstanceActivityWorker) Stop() {
	w.isRunning = false
}
