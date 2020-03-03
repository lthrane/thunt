package main

import (
	"time"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
)

// Hunt represents a treasure hunt that can be undertaken
type Hunt struct {
	ID        int64     `json:"ID" datastore:"-"`
	MissionID int64     `json:"missionID" datastore:"missionID"`
	Timestamp time.Time `json:"name" datastore:"timestamp"`
}

// HuntHandler is ...
type HuntHandler struct {
	EntityHandler
}

// NewHuntHandler is ...
func NewHuntHandler(client *datastore.Client, router *mux.Router) *EntityHandler {
	handler := &EntityHandler{
		client:        client,
		kind:          "Hunt",
		entityFactory: func() Entity { return new(Hunt) },
		listFactory:   func() interface{} { return new([]Hunt) },
		entityAt: func(anonymous interface{}, index int) Entity {
			return &(*anonymous.(*[]Hunt))[index]
		},
	}

	router.HandleFunc("/hunt", handler.List).Methods("GET")
	router.HandleFunc("/hunt", handler.Put).Methods("PUT")
	router.HandleFunc("/hunt/{id}", handler.Get).Methods("GET")
	router.HandleFunc("/hunt/{id}", handler.Delete).Methods("DELETE")

	return handler
}

func (m *Hunt) setID(id int64) {
	m.ID = id
}
