package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/emojify-app/face-detection/handlers"
	"github.com/emojify-app/face-detection/logging"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var listenAddress = env.String("ADDRESS", false, "127.0.0.1", "Listen address for the server")
var listenPort = env.String("PORT", false, "9090", "Listen port for the server")
var statsDAddress = env.String("STATSD", false, "localhost:8125", "Location of the statsd collector")
var logLevel = env.String("LOG_LEVEL", false, "INFO", "Log level [INFO,DEBUG,TRACE]")

func main() {
	// Parse the config env vars
	err := env.Parse()
	if err != nil {
		fmt.Println(env.Help())
		os.Exit(1)
	}

	// setup the default logger
	l, err := logging.New("facedetection", *statsDAddress, *logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	h := handlers.NewHealth(l)

	r := mux.NewRouter()
	r.Handle("/health", h).Methods("GET")

	l.ServiceStart(*listenAddress, *listenPort)
	fmt.Println("Error starting server", "error", http.ListenAndServe(":9090", r))
}
