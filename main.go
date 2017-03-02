package main

import (
	"net/http"

	"fmt"

	"log"

	astro "github.com/astromechio/astrohub/astrolib"
	"github.com/gorilla/mux"
)

func main() {
	queue := astro.SharedRequestQueue()
	resMap := astro.SharedResponseMap()

	mux := mux.NewRouter()

	mux.Methods("POST").Path("/service/{name}").HandlerFunc(LogWrapper(QueueRequestHandler(queue)))
	mux.Methods("GET").Path("/service/{name}/jobs").HandlerFunc(LogWrapper(JobRequestHandler(queue)))
	mux.Methods("POST").Path("/jobs/{id}/response").HandlerFunc(LogWrapper(JobResponseHandler(resMap)))

	errChan := make(chan error)
	go startServer(3000, mux, errChan)

	for err := range errChan {
		log.Println(err)
	}
}

func startServer(port int, mux http.Handler, errChan chan (error)) {
	log.Println("Starting server on", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		errChan <- err
	}
}

func LogWrapper(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.RequestURI

		log.Printf("%s %s", method, path)

		h(w, r)
	}
}
