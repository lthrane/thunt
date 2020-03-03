package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Entity is ...
type Entity interface {
	setID(id int64)
}

// EntityHandler is ...
type EntityHandler struct {
	client        *datastore.Client
	kind          string
	entityFactory func() Entity
	listFactory   func() interface{}
	entityAt      func(entities interface{}, index int) Entity
}

// Put is ...
func (eh *EntityHandler) Put(resp http.ResponseWriter, rqst *http.Request) {
	entity := eh.entityFactory()
	if decode(rqst, entity) != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IncompleteKey(eh.kind, nil)
	key, err := eh.client.Put(rqst.Context(), key, entity)
	if err != nil {
		log.Fatalf("datastore.Put: %v", err)
	}

	encode(resp, key.ID)
}

// Get is ...
func (eh *EntityHandler) Get(resp http.ResponseWriter, rqst *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(rqst)["id"], 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IDKey(eh.kind, ID, nil)
	entity := eh.entityFactory()
	err = eh.client.Get(rqst.Context(), key, entity)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	entity.setID(key.ID)

	encode(resp, entity)
}

// Delete is ...
func (eh *EntityHandler) Delete(resp http.ResponseWriter, rqst *http.Request) {
	ID, err := strconv.ParseInt(mux.Vars(rqst)["id"], 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IDKey(eh.kind, ID, nil)
	err = eh.client.Delete(rqst.Context(), key)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
	}
}

// List is ...
func (eh *EntityHandler) List(resp http.ResponseWriter, rqst *http.Request) {
	var querySink interface{} = eh.listFactory()

	query := datastore.NewQuery(eh.kind)
	keys, err := eh.client.GetAll(rqst.Context(), query, querySink)
	if err != nil {
		Error.Printf("GetAll: %v", err)
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := encodeRaw(resp, "["); err != nil {
		Error.Printf("encodeRaw: %v", err)
		resp.WriteHeader(http.StatusServiceUnavailable)
	}

	empty := true
	for index, key := range keys {
		entity := eh.entityAt(querySink, index)
		entity.setID(key.ID)

		if !empty {
			if err := encodeRaw(resp, ","); err != nil {
				Error.Printf("encodeRaw: %v", err)
				resp.WriteHeader(http.StatusServiceUnavailable)
			}
		}

		if err := encode(resp, entity); err != nil {
			Error.Printf("encode: %v", err)
			resp.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		empty = false
	}

	if err := encodeRaw(resp, "]"); err != nil {
		Error.Printf("listMissions: %v", err)
		resp.WriteHeader(http.StatusServiceUnavailable)
	}
}

func encode(resp http.ResponseWriter, v interface{}) error {
	encoder := json.NewEncoder(resp)
	encoder.SetEscapeHTML(true)
	err := encoder.Encode(v)
	if err != nil {
		return errors.Errorf("Encode: %v", err)
	}
	return nil
}

func encodeRaw(resp http.ResponseWriter, s string) error {
	nBytes, err := resp.Write([]byte(s))
	if nBytes != len(s) || err != nil {
		return errors.Errorf("can't emit %s", s)
	}
	return nil
}

func decode(rqst *http.Request, v interface{}) error {
	decoder := json.NewDecoder(rqst.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
