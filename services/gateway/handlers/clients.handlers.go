package gtwhandlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	gtwcontainer "github.com/hetacode/mechelon/services/gateway/container"
)

// ClientsHandlers struct is keeping and linking all needed handlers required by clients part of api
// Client part of api is consuming by external services in order to expose their state of activity
type ClientsHandlers struct {
	container *gtwcontainer.Container
}

// NewClientsHandlers constructor
func NewClientsHandlers(c *gtwcontainer.Container, h *mux.Router) {
	hc := &ClientsHandlers{
		container: c,
	}

	h.HandleFunc("/register", hc.RegisterServiceHandler)
}

// RegisterServiceBody request
type RegisterServiceBody struct {
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// RegisterServiceHandler - handle endpoint to register new project service with given instance name
// This endpoint is even using to register  new another instance of existing service
func (h *ClientsHandlers) RegisterServiceHandler(w http.ResponseWriter, r *http.Request) {
	bb, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "body error", http.StatusBadRequest)
		return
	}

	var body *RegisterServiceBody
	if err := json.Unmarshal(bb, &body); err != nil {
		http.Error(w, "body parse", http.StatusBadRequest)
		return
	}

	id, _ := uuid.NewV4()
	command := &eventsservicesmgmt.RegisterServiceCommand{
		EventData:    &goeh.EventData{ID: id.String()},
		ProjectName:  body.ProjectName,
		ServiceName:  body.ServiceName,
		InstanceName: body.InstanceName,
	}

	if err := h.container.CommandsBusProducer.Publish(command); err != nil {
		http.Error(w, "send command: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
