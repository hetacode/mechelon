package smgeventstore

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"reflect"

	"github.com/EventStore/EventStore-Client-Go/client"
	"github.com/EventStore/EventStore-Client-Go/direction"
	"github.com/EventStore/EventStore-Client-Go/errors"
	"github.com/EventStore/EventStore-Client-Go/messages"
	"github.com/EventStore/EventStore-Client-Go/streamrevision"
	"github.com/gofrs/uuid"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/arch"
)

// ServiceEventStore - implementation of EventStore base on EventStoreDB
type ServiceEventStore struct {
	EventStoreClient *client.Client
	EventsMapper     *goeh.EventsMapper
}

// NewServiceEventStore instance
func NewServiceEventStore(em *goeh.EventsMapper) arch.EventStore {
	c, e := client.NewClient(&client.Configuration{Address: os.Getenv("EVENTSTOREDB_SERVER"), DisableTLS: true})
	if e != nil {
		panic(e)
	}

	if e = c.Connect(); e != nil {
		panic(e)
	}

	es := &ServiceEventStore{
		EventStoreClient: c,
		EventsMapper:     em,
	}

	return es
}

// GetSnapshot of state
func (s *ServiceEventStore) GetSnapshot(key string, stateType interface{}) (state interface{}) {
	streamName := key + "-snapshot"
	events, err := s.EventStoreClient.ReadStreamEvents(context.Background(), direction.Backwards, streamName, streamrevision.StreamRevisionEnd, 1, true)
	if err != nil && err != errors.ErrStreamNotFound {
		log.Printf("GetSnapshot err: %s", err)
		return nil
	} else if err != nil {
		return nil
	}
	if len(events) > 0 {
		data := events[0].Data
		t := reflect.TypeOf(stateType).Elem()
		i := reflect.New(t).Interface()
		if err := json.Unmarshal(data, &i); err != nil {
			log.Printf("cannot resolve snapshot data for '%s' type", t.Name())
			return nil
		}
		return i
	}
	return nil
}

// SaveSnapshot - save new snapshot to special stream
func (s *ServiceEventStore) SaveSnapshot(key string, state interface{}) error {
	streamName := key + "-snapshot"
	oid, _ := uuid.NewV4()
	bytesState, err := json.Marshal(state)

	if err != nil {
		log.Printf("SaveSnapshot json marshal err: %s", err)
		return err
	}

	snapEv := messages.ProposedEvent{
		EventID:      oid,
		EventType:    "Snapshot",
		ContentType:  "application/json",
		Data:         bytesState,
		UserMetadata: make([]byte, 0),
	}
	_, err = s.EventStoreClient.AppendToStream(context.Background(), streamName, streamrevision.StreamRevisionAny, []messages.ProposedEvent{snapEv})

	if err != nil {
		log.Printf("SaveSnapshot AppendToStream err: %s", err)
		return err
	}

	return nil
}

// GetEvents all for givent key and position
func (s *ServiceEventStore) GetEvents(key string, position int64) []goeh.Event {
	streamName := key + "-events"
	result := make([]goeh.Event, 0)
	for {
		events, err := s.EventStoreClient.ReadStreamEvents(context.Background(), direction.Forwards, streamName, uint64(position+1), 20, true)
		if err != nil {
			log.Printf("GetEvents err: %s", err)
			return result
		}
		for _, ev := range events {
			data := ev.Data
			mappedEvent, err := s.EventsMapper.Resolve(string(data))
			if err != nil {
				log.Printf("cannot resoleve event data err: %s", err)
				return nil
			}

			ee := mappedEvent.(arch.ExtendedEvent)
			ee.SetVersion(ev.EventNumber)

			result = append(result, mappedEvent)
		}
		if len(events) < 20 {
			break
		} else {
			position = int64(events[:1][0].EventNumber)
		}
	}
	return result
}

// SaveNewEvents - save new events to store
func (s *ServiceEventStore) SaveNewEvents(key string, events []goeh.Event) error {
	streamName := key + "-events"

	storeEvents := make([]messages.ProposedEvent, 0)
	for _, ev := range events {
		ev.SavePayload(ev)

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
	_, err := s.EventStoreClient.AppendToStream(context.Background(), streamName, streamrevision.StreamRevisionAny, storeEvents)
	if err != nil {
		log.Printf("PushNewEvents AppendToStream err: %s", err)
		return err
	}

	return nil
}
