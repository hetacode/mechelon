package gtwhandlers

import (
	"github.com/gorilla/mux"
	gtwcontainer "github.com/hetacode/mechelon/services/gateway/container"
)

// FrontendHandlers struct is keeping and linking all needed handlers required by frontend application
type FrontendHandlers struct {
	container *gtwcontainer.Container
}

// NewFrontendHandlers constructor
func NewFrontendHandlers(c *gtwcontainer.Container, h *mux.Router) {
	hc := &FrontendHandlers{
		container: c,
	}

	panic(hc) // TODO: to implement
}
