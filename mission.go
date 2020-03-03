package main

import (
	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

// Mission represents a treasure hunt that can be undertaken
type Mission struct {
	ID          int64  `json:"id" datastore:"id"`
	Name        string `json:"name" datastore:"name"`
	Description string `json:"description" datastore:"description"`
}

// MissionHandler is ...
type MissionHandler struct {
	EntityHandler
}

// NewMissionHandler is ...
func NewMissionHandler(client *datastore.Client, router *mux.Router) *EntityHandler {
	handler := &EntityHandler{
		client:        client,
		kind:          "Mission",
		entityFactory: func() Entity { return new(Mission) },
		listFactory:   func() interface{} { return new([]Mission) },
		entityAt: func(anonymous interface{}, index int) Entity {
			return &(*anonymous.(*[]Mission))[index]
		},
	}

	router.HandleFunc("/mission", handler.List).Methods("GET")
	router.HandleFunc("/mission", handler.Put).Methods("PUT")
	router.HandleFunc("/mission/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/mission/{id}", handler.Delete).Methods("DELETE")

	return handler
}

func (m *Mission) setID(id int64) {
	m.ID = id
}
