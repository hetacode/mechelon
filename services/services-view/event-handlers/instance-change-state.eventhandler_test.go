// Readme
// To run these test should be launched docker-compose before in order to get a running mongodb instance
// Important! Tests are invoking setup() method on the start - it means the database is drop each time

package svveventhandlers

import (
	"context"
	"os"
	"testing"
	"time"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	svvdb "github.com/hetacode/mechelon/services/services-view/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupInstanceChangeState() *svvcontainer.Container {
	err := godotenv.Load("../../../.env.dev")
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_SERVER")))
	if err != nil {
		panic(err)
	}
	db := client.Database(os.Getenv("SVC_SERVICES_VIEW_DB_NAME"))
	err = db.Drop(context.Background())
	if err != nil {
		panic(err)
	}
	collection := db.Collection(os.Getenv("SVC_SERVICES_VIEW_DB_PROJECTS_COLLECTION"))

	container := &svvcontainer.Container{
		ProjectsMongoDBCollection: collection,
	}

	return container
}

// TestInstancesListShouldBeEmptyWhenInstanceIsAddedAndNextRemovedByEvents test
func TestChangeInstanceStatesFlowActiveIdleInactive(t *testing.T) {
	c := setupInstanceChangeState()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	rh := &InstanceAddedToServiceEventHandler{
		Container: c,
	}
	ih := &InstanceGotIdleEventHandler{
		Container: c,
	}
	iah := &InstanceGotInactiveEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}
	re := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}
	idlee := &eventsservicesmgmt.InstanceGotIdleEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}
	inactivee := &eventsservicesmgmt.InstanceGotInactiveEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	// Execute event
	h.Handle(e)

	// After attach instance status should be active
	rh.Handle(re)
	project := getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.Active))

	// Change to idle
	ih.Handle(idlee)
	project = getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.Idle))

	// Change to inactive
	iah.Handle(inactivee)
	project = getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.InActive))
}

// TestInstancesListShouldBeEmptyWhenInstanceIsAddedAndNextRemovedByEvents test
func TestChangeInstanceStatesFlowActiveIdleInactiveAndActiveAgain(t *testing.T) {
	c := setupInstanceChangeState()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	rh := &InstanceAddedToServiceEventHandler{
		Container: c,
	}
	ih := &InstanceGotIdleEventHandler{
		Container: c,
	}
	iah := &InstanceGotInactiveEventHandler{
		Container: c,
	}
	ah := &InstanceActivatedEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}
	re := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}
	idlee := &eventsservicesmgmt.InstanceGotIdleEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}
	inactivee := &eventsservicesmgmt.InstanceGotInactiveEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}
	activee := &eventsservicesmgmt.InstanceActivatedEvent{
		EventData:    &goeh.EventData{ID: "1234"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	// Execute event
	h.Handle(e)

	// After attach instance status should be active
	rh.Handle(re)
	project := getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.Active))

	// Change to idle
	ih.Handle(idlee)
	project = getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.Idle))

	// Change to inactive
	iah.Handle(inactivee)
	project = getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.InActive))

	// Change to active
	ah.Handle(activee)
	project = getProjectService(c, t)
	assert.Equal(t, project.Instances[0].Status, string(smgeventstore.Active))
}

func getProjectService(c *svvcontainer.Container, t *testing.T) *svvdb.ServiceEntity {
	res := c.ProjectsMongoDBCollection.FindOne(context.Background(), bson.M{"project_name": "test-project", "service_name": "test-project-service"})
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project *svvdb.ServiceEntity
	if err := res.Decode(&project); err != nil {
		t.Error(err)
	}
	return project
}
