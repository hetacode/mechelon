package smgworkers

import (
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
	repository  *smgeventstore.ServiceStateRepository
	projectName string
	serviceName string
	period      int64 // seconds
	createdAt   int64
	updatedAt   int64
}

// NewInstanceActivityWorker create new instance
func NewInstanceActivityWorker(projectName, serviceName string, period int64) *InstanceActivityWorker {
	worker := &InstanceActivityWorker{
		projectName: projectName,
		serviceName: serviceName,
		period:      period,
		createdAt:   time.Now().Unix(),
	}
	worker.updatedAt = worker.createdAt

	return worker
}

// Start worker
func (w *InstanceActivityWorker) Start() {
	for {
		if time.Now().Unix()-w.updatedAt < w.period {
			continue
		}

		aggr := w.repository.GetAggregator(w.projectName, w.serviceName)
		if aggr == nil {
			log.Printf("InstanceActivityWorker | cannot find an aggregator for %s service for %s project", w.serviceName, w.projectName)
			// TODO: disable worker
			return
		}

		for _, instance := range aggr.State.Instances {
			switch instance.State {
			case smgeventstore.Active:
				if time.Now().Unix()-instance.UpdatedAt > ValidActivePeriod {
					aggr.SetInstanceAsIdle(instance.Name)
				}
			case smgeventstore.Idle:
				if time.Now().Unix()-instance.UpdatedAt > ValidIdlePeriod {
					aggr.SetInstanceAsInactive(instance.Name)
				}
			case smgeventstore.InActive:
				// TODO: disable worker?
			}
		}
		events := aggr.GetPendingEvents()
		if err := w.repository.SaveEvents(w.projectName, w.serviceName, events); err != nil {
			log.Printf("InstanceActivityWorker | cannot save events for %s service for %s project", w.serviceName, w.projectName)
		}
	}
}
