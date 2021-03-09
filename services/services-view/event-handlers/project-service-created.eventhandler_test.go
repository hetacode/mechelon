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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gotest.tools/assert"
)

func setup() *svvcontainer.Container {
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

// TestShouldCreateProjectWithServiceOnce test
func TestShouldCreateProjectWithServiceOnce(t *testing.T) {
	c := setup()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}

	// Execute event
	h.Handle(e)

	// Checks
	res := c.ProjectsMongoDBCollection.FindOne(context.Background(), bson.M{"name": "test-project"})
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project *svvdb.ProjectEntity
	res.Decode(&project)

	assert.Equal(t, project.Name, "test-project")
	assert.Equal(t, len(project.Services), 1)
	assert.Equal(t, project.Services[0].Name, "test-project-service")
}

// TestShouldCreateOneProjectWithTwoService test
func TestShouldCreateOneProjectWithTwoService(t *testing.T) {
	c := setup()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
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

	// Execute events
	h.Handle(e1)
	h.Handle(e2)

	// Checks
	res := c.ProjectsMongoDBCollection.FindOne(context.Background(), bson.M{"name": "test-project"})
	if res.Err() != nil {
		t.Error(res.Err())
	}

	var project *svvdb.ProjectEntity
	res.Decode(&project)

	assert.Equal(t, project.Name, "test-project")
	assert.Equal(t, len(project.Services), 2)
	assert.Equal(t, project.Services[0].Name, "test-project-service-1")
	assert.Equal(t, project.Services[1].Name, "test-project-service-2")
}
