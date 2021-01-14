package smgeventstore

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/gofrs/uuid"
	goeh "github.com/hetacode/go-eh"
)

// ServiceEventStore - implementation of EventStore base on EventStoreDB
type ServiceEventStore struct {
	EventStoreClient *client.Client
	EventsMapper     *goeh.EventsMapper
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

// GetEvents all for givent key and position
func (s *ServiceEventStore) GetEvents(key string, position int64) []goeh.Event {
	streamName := key + "-events"
	result := make([]goeh.Event, 0)
	for {
		events, err := s.EventStoreClient.ReadStreamEvents(context.Background(), direction.Forwards, streamName, uint64(position), 20, true)
		if err != nil {
			log.Printf("GetEvents err: %s", err)
			return nil
		}
		for _, ev := range events {
			data := ev.Data
			mappedEvent, err := s.EventsMapper.Resolve(string(data))
			if err != nil {
				log.Printf("cannot resoleve event data err: %s", err)
				return nil
			}
			result = append(result, mappedEvent)
		}
		if len(events) < 20 {
			break
		} else {
			position = int64(events[:1][0].Position.Commit)
		}
	}
	return result
}

func (s *ServiceEventStore) PushNewEvents(key string, events []goeh.Event) error {
	streamName := key + "-events"

	storeEvents := make([]messages.ProposedEvent, 0)
	for _, ev := range events {
		oid, _ := uuid.NewV4()
		esEv := messages.ProposedEvent{
			EventID:      oid,
			EventType:    ev.GetType(),
			ContentType:  "application/json",
			Data:         []byte(ev.GetPayload()),
			UserMetadata: make([]byte, 0),
		}
		storeEvents = append(storeEvents, esEv)
	}
	r, err := s.EventStoreClient.AppendToStream(context.Background(), streamName, streamrevision.StreamRevisionAny, storeEvents)
	if err != nil {
		log.Printf("PushNewEvents AppendToStream err: %s", err)
		return err
	}

	// TODO: save snapshot here?
}
