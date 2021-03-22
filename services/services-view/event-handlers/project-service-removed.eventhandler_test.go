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

func setupProjectServiceRemoved() *svvcontainer.Container {
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
func TestProjectServiceShouldNotExistsAfterProjectServiceRemovedEventExecute(t *testing.T) {
	c := setupProjectServiceRemoved()
	defer c.ProjectsMongoDBCollection.Database()

	// Init handler
	h := &ProjectServiceCreatedEventHandler{
		Container: c,
	}
	rh := &ProjectServiceRemovedEventHandler{
		Container: c,
	}

	// Create event
	e := &eventsservicesmgmt.ProjectServiceCreatedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}
	re := &eventsservicesmgmt.ProjectServiceRemovedEvent{
		EventData:   &goeh.EventData{ID: "1234"},
		ProjectName: "test-project",
		ServiceName: "test-project-service",
	}
	// Execute event
	h.Handle(e)
	rh.Handle(re)

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

	assert.Len(t, project, 0)
}
