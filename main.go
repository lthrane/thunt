package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Mission represents a treasure hunt that can be undertaken
type Mission struct {
	ID          uuid.UUID `json:id`
	Name        string    `json:name`
	Description string    `json:description`
}

var mutex sync.RWMutex
var missionByID map[uuid.UUID]*Mission = map[uuid.UUID]*Mission{}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/mission", listMissions).Methods("GET")
	router.HandleFunc("/mission", putMission).Methods("PUT")
	router.HandleFunc("/mission/{id}", getMission).Methods("GET")
	router.HandleFunc("/mission/{id}", deleteMission).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func putMission(resp http.ResponseWriter, rqst *http.Request) {
	var mission *Mission = new(Mission)
	if decode(rqst, mission) != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// defensive code against UUID collisions
	for {
		var err error
		mission.ID, err = uuid.NewRandom()
		if err != nil {
			Error.Printf("unable to create UUID: %v", err)
			resp.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		if missionByID[mission.ID] == nil {
			break
		}
	}

	missionByID[mission.ID] = mission

	encode(resp, mission.ID)
}

func getMission(resp http.ResponseWriter, rqst *http.Request) {
	encodedID := mux.Vars(rqst)["id"]
	id, err := uuid.Parse(encodedID)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex.RLock()
	defer mutex.RUnlock()

	mission := missionByID[id]
	if mission == nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	encode(resp, mission)
}

func deleteMission(resp http.ResponseWriter, rqst *http.Request) {
	encodedID := mux.Vars(rqst)["id"]
	id, err := uuid.Parse(encodedID)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	mission := missionByID[id]
	if mission == nil {
		resp.WriteHeader(http.StatusNotFound)
		return
	}

	delete(missionByID, mission.ID)

	err = encode(resp, mission)
	if err != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}

func listMissions(resp http.ResponseWriter, rqst *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	if encode(resp, missionByID) != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		return
	}
}

func encode(resp http.ResponseWriter, v interface{}) error {
	encoder := json.NewEncoder(resp)
	encoder.SetEscapeHTML(true)
	err := encoder.Encode(v)
	if err != nil {
		Error.Printf("Encode: %v", err)
		return err
	}
	return nil
}

func decode(rqst *http.Request, v interface{}) error {
	decoder := json.NewDecoder(rqst.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}
