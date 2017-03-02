package main

import (
	"net/http"

	"encoding/json"

	"io/ioutil"

	astro "github.com/astromechio/astrohub/astrolib"
	"github.com/gorilla/mux"
)

func QueueRequestHandler(queue astro.AQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceName, ok := mux.Vars(r)["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		originalJSON, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		astroReq := astro.NewReq(serviceName, originalJSON)

		astroRes, resErr := queue.QueueRequest(astroReq)
		if resErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rawRes, rawErr := astroRes.Response()
		if rawErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(rawRes)
	}
}

func JobRequestHandler(queue astro.AQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceName, ok := mux.Vars(r)["name"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		astroReq, reqErr := queue.GetRequest(serviceName)
		if reqErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reqJSON, _ := json.Marshal(astroReq)

		w.WriteHeader(http.StatusOK)
		w.Write(reqJSON)
	}
}

func JobResponseHandler(resMap *astro.ResponseMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, ok := mux.Vars(r)["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		resBody, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sendErr := resMap.SendResponse(reqID, resBody)
		if sendErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
