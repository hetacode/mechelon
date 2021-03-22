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
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	svvdb "github.com/hetacode/mechelon/services/services-view/db"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupInstanceAddedToService() *svvcontainer.Container {
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

// TestShouldAttachNewInstanceOfServiceToProject test
func TestShouldAttachNewInstanceOfServiceToProject(t *testing.T) {
	c := setupInstanceAddedToService()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	ah := &InstanceAddedToServiceEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}

	ae := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	// Execute events
	h.Handle(e)
	ah.Handle(ae)

	// Checks
	res, err := c.ProjectsMongoDBCollection.Find(context.Background(), bson.M{"project_name": "test-project"})
	if err != nil {
		t.Error(err)
	}
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project []svvdb.ServiceEntity
	if err := res.All(context.Background(), &project); err != nil {
		t.Error(err)
	}

	assert.Len(t, project, 1)

	assert.Equal(t, project[0].ProjectName, "test-project")
	assert.Equal(t, project[0].ServiceName, "test-project-service")
	assert.Len(t, project[0].Instances, 1)
	assert.Equal(t, project[0].Instances[0].Name, "inst-1")
}

// TestShouldAttachTwoInstancesOfServiceToProject test
func TestShouldAttachTwoInstancesOfServiceToProject(t *testing.T) {
	c := setupInstanceAddedToService()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	ah := &InstanceAddedToServiceEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}

	ae1 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	ae2 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-2",
	}

	// Execute events
	h.Handle(e)
	ah.Handle(ae1)
	ah.Handle(ae2)

	// Checks
	res := c.ProjectsMongoDBCollection.FindOne(context.Background(), bson.M{"project_name": "test-project", "service_name": "test-project-service"})
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project *svvdb.ServiceEntity
	res.Decode(&project)

	assert.Equal(t, project.ProjectName, "test-project")
	assert.Equal(t, project.ServiceName, "test-project-service")
	assert.Len(t, project.Instances, 2)
	assert.Equal(t, project.Instances[0].Name, "inst-1")
	assert.Equal(t, project.Instances[1].Name, "inst-2")
}

// TestShouldAttachOneInstanceToTheServiceWhenInstanceAddedToServiceEventCallEventHandlerTwice test
func TestShouldAttachOneInstanceToTheServiceWhenInstanceAddedToServiceEventCallEventHandlerTwice(t *testing.T) {
	c := setupInstanceAddedToService()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	ah := &InstanceAddedToServiceEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}

	ae1 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	ae2 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service",
		InstanceName: "inst-1",
	}

	// Execute events
	h.Handle(e)
	ah.Handle(ae1)
	ah.Handle(ae2)

	// Checks
	res := c.ProjectsMongoDBCollection.FindOne(context.Background(), bson.M{"project_name": "test-project", "service_name": "test-project-service"})
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project *svvdb.ServiceEntity
	res.Decode(&project)

	assert.Equal(t, project.ProjectName, "test-project")
	assert.Equal(t, project.ServiceName, "test-project-service")
	assert.Len(t, project.Instances, 1)
	assert.Equal(t, project.Instances[0].Name, "inst-1")
}

func TestShouldAttachInstanceToTheFirstOfProjectService(t *testing.T) {
	c := setupInstanceAddedToService()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	ah := &InstanceAddedToServiceEventHandler{
		Container: c,
	}

	// Create events
	e1 := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service-1",
	}
	e2 := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service-2",
	}
	ae := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service-1",
		InstanceName: "inst-1",
	}

	// Execute events
	h.Handle(e1)
	h.Handle(e2)
	ah.Handle(ae)

	// Checks
	res, err := c.ProjectsMongoDBCollection.Find(context.Background(), bson.M{"project_name": "test-project"})
	if err != nil {
		t.Error(err)
	}
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project []svvdb.ServiceEntity
	if err := res.All(context.Background(), &project); err != nil {
		t.Error(err)
	}

	assert.Len(t, project, 2)
	assert.Equal(t, project[0].ProjectName, "test-project")
	assert.Equal(t, project[1].ProjectName, "test-project")
	assert.Equal(t, project[0].ServiceName, "test-project-service-1")
	assert.Equal(t, project[1].ServiceName, "test-project-service-2")
	assert.Len(t, project[0].Instances, 1)
	assert.Len(t, project[1].Instances, 0)
	assert.Equal(t, project[0].Instances[0].Name, "inst-1")
}

func TestShouldAttachInstancesToTheBothOfProjectServices(t *testing.T) {
	c := setupInstanceAddedToService()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	ah := &InstanceAddedToServiceEventHandler{
		Container: c,
	}

	// Create events
	e1 := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service-1",
	}
	e2 := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service-2",
	}

	// First instance for first service
	ae1s1 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service-1",
		InstanceName: "inst-1",
	}

	// First instance for second service
	ae1s2 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service-2",
		InstanceName: "inst-1",
	}

	// Second instance for second service
	ae2s2 := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: "1235"},
		ProjectName:  "test-project",
		ServiceName:  "test-project-service-2",
		InstanceName: "inst-2",
	}

	// Execute events
	h.Handle(e1)
	h.Handle(e2)
	ah.Handle(ae1s2)
	ah.Handle(ae1s1)
	ah.Handle(ae2s2)

	// Checks
	res, err := c.ProjectsMongoDBCollection.Find(context.Background(), bson.M{"project_name": "test-project"})
	if err != nil {
		t.Error(err)
	}
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project []svvdb.ServiceEntity
	if err := res.All(context.Background(), &project); err != nil {
		t.Error(err)
	}

	assert.Len(t, project, 2)
	assert.Equal(t, project[0].ProjectName, "test-project")
	assert.Equal(t, project[1].ProjectName, "test-project")
	assert.Equal(t, project[0].ServiceName, "test-project-service-1")
	assert.Equal(t, project[1].ServiceName, "test-project-service-2")
	assert.Len(t, project[0].Instances, 1)
	assert.Len(t, project[1].Instances, 2)
	assert.Equal(t, project[0].Instances[0].Name, "inst-1")
	assert.Equal(t, project[1].Instances[0].Name, "inst-1")
	assert.Equal(t, project[1].Instances[1].Name, "inst-2")
}
