package gtwhandlers

import (
	"net/http"

	gtwcontainer "github.com/hetacode/mechelon/services/gateway/container"
)

// ClientsHandlers struct is keeping and linking all needed handlers required by clients part of api
// Client part of api is consuming by external services in order to expose their state of activity
type ClientsHandlers struct {
	container *gtwcontainer.Container
}

// NewClientsHandlers constructor
func NewClientsHandlers(c *gtwcontainer.Container, h http.Handler) {
	hc := &ClientsHandlers{
		container: c,
	}

	panic(hc)
}
