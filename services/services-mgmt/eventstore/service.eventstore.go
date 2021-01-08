package smgeventstore

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/direction"
	goeh "github.com/hetacode/go-eh"
)

// ServiceEventStore - implementation of EventStore base on EventStoreDB
type ServiceEventStore struct {
	EventStoreClient *client.Client
}

// GetSnapshot of state
func (s *ServiceEventStore) GetSnapshot(key string, stateType interface{}) (state interface{}, lastEventPosition int64) {
	streamName := key + "-snapshot"
	events, err := s.EventStoreClient.ReadStreamEvents(context.Background(), direction.Backwards, streamName, 0, 1, true)
	if err != nil {
		log.Printf("GetSnapshot err: %s", err)
		return nil, 0
	}
	if len(events) > 0 {
		data := events[0].Data
		t := reflect.TypeOf(stateType).Elem()
		i := reflect.New(t).Interface()
		if err := json.Unmarshal(data, &i); err != nil {
			log.Printf("cannot resolve snapshot data for '%s' type", t.Name())
			return nil, 0
		}
		state := i.(State)
		return i, int64(state.GetVersion())
	}
	return nil, 0
}

func (s *ServiceEventStore) GetEvents(key string, position int64) []goeh.Event {
	return nil
}
