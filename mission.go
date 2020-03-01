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

// Mission represents a treasure hunt that can be undertaken
type Mission struct {
	ID          int64  `json:"id" datastore:"id"`
	Name        string `json:"name" datastore:"name"`
	Description string `json:"description" datastore:"description"`
}

func register(router *mux.Router) {
	router.HandleFunc("/mission", listMissions).Methods("GET")
	router.HandleFunc("/mission", putMission).Methods("PUT")
	router.HandleFunc("/mission/{id}", getMission).Methods("GET")
	router.HandleFunc("/mission/{id}", deleteMission).Methods("DELETE")
}

func putMission(resp http.ResponseWriter, rqst *http.Request) {
	Info.Printf("putMission\n")

	var mission *Mission = new(Mission)
	if decode(rqst, mission) != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IncompleteKey("Mission", nil)
	key, err := client.Put(rqst.Context(), key, mission)
	if err != nil {
		log.Fatalf("datastore.Put: %v", err)
	}
	mission.ID = key.ID

	encode(resp, mission.ID)
}

func getMission(resp http.ResponseWriter, rqst *http.Request) {
	Info.Printf("getMission\n")

	ID, err := strconv.ParseInt(mux.Vars(rqst)["id"], 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mission := new(Mission)
	key := datastore.IDKey("Mission", ID, nil)
	err = client.Get(rqst.Context(), key, mission)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	encode(resp, mission)
}

func deleteMission(resp http.ResponseWriter, rqst *http.Request) {
	Info.Printf("deleteMission\n")

	ID, err := strconv.ParseInt(mux.Vars(rqst)["id"], 10, 64)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IDKey("Mission", ID, nil)
	err = client.Delete(rqst.Context(), key)
}

func listMissions(resp http.ResponseWriter, rqst *http.Request) {
	Info.Printf("listMissions\n")

	var missions []*Mission

	query := datastore.NewQuery("Mission")
	keys, err := client.GetAll(rqst.Context(), query, &missions)
	if err != nil {
		Error.Printf("GetAll: %v", err)
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := encodeRaw(resp, "["); err != nil {
		Error.Printf("listMissions: %v", err)
		resp.WriteHeader(http.StatusServiceUnavailable)
	}

	empty := true
	for i, key := range keys {
		if !empty {
			if err := encodeRaw(resp, ","); err != nil {
				Error.Printf("listMissions: %v", err)
				resp.WriteHeader(http.StatusServiceUnavailable)
			}
		}

		missions[i].ID = key.ID
		if encode(resp, missions[i]) != nil {
			Error.Printf("can't emit element: %v", missions[i])
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
