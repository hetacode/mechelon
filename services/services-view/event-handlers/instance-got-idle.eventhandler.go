package svveventhandlers

import (
	"context"
	"log"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	vssdb "github.com/hetacode/mechelon/services/services-view/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InstanceAddedToServiceEventHandler struct
type InstanceGotIdleEventHandler struct {
	Container *svvcontainer.Container
}

// Handle ProjectServiceCreatedEvent
func (e *InstanceGotIdleEventHandler) Handle(event goeh.Event) {
	log.Printf("InstanceGotIdleEventHandler start")
	ev := event.(*eventsservicesmgmt.InstanceGotIdleEvent)
	// Check if project and service exists
	filter := bson.M{
		"project_name": ev.ProjectName,
		"service_name": ev.ServiceName,
	}

	var project *vssdb.ServiceEntity
	res := e.Container.ProjectsMongoDBCollection.FindOne(context.Background(), filter)
	if res.Err() == mongo.ErrNoDocuments {
		log.Printf("\033[31mInstanceGotIdleEventHandler cannot find project '%s' with service '%s'\033[0m", ev.ProjectName, ev.ServiceName)
		return
	} else if res.Err() != nil {
		log.Printf("\033[31mInstanceGotIdleEventHandler find project '%s' with service '%s' err: %s\033[0m", ev.ProjectName, ev.ServiceName, res.Err())
		return
	}

	err := res.Decode(&project)
	if err != nil {
		log.Printf("\033[31mInstanceGotIdleEventHandler decode project entity err: %s\033[0m", err)
		return
	}

	for idx, i := range project.Instances {
		if i.Name == ev.InstanceName {
			project.Instances[idx].Status = string(smgeventstore.Idle)
			project.Instances[idx].UpdatedAt = ev.UpdateAt

			_, err = e.Container.ProjectsMongoDBCollection.UpdateOne(context.Background(), filter, bson.M{"$set": &project})
			if err != nil {
				log.Printf("\033[31mInstanceGotIdleEventHandler instance '%s' for service '%s' in project '%s' update err: %s\033[0m", ev.InstanceName, ev.ServiceName, ev.ProjectName, err)

			}
			return
		}
	}
	log.Printf("\033[31mInstanceGotIdleEventHandler instance '%s' for service '%s' in project '%s' doesn't exists\033[0m", ev.InstanceName, ev.ServiceName, ev.ProjectName)

	log.Printf("InstanceGotIdleEventHandler end")
}
