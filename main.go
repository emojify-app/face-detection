package main

import (
	"fmt"
	"net/http"
	"os"

	"net/http/pprof"

	"github.com/emojify-app/face-detection/handlers"
	"github.com/emojify-app/face-detection/logging"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
)

var listenAddress = env.String("BIND_ADDRESS", false, "127.0.0.1", "Listen address for the server")
var listenPort = env.String("BIND_PORT", false, "9090", "Listen port for the server")
var statsDAddress = env.String("STATSD_ADDRESS", false, "localhost:8125", "Location of the statsd collector")
var logLevel = env.String("LOG_LEVEL", false, "info", "log level [info,debug,trace]")
var cascadeFolder = env.String("CASCADE_FOLDER", false, "./cascades", "location of the OpenCV cascades")
var profile = env.Bool("PROFILE", false, false, "enable pprof profiling")

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

	// check the cascades folder exists
	f, err := os.Open(*cascadeFolder)
	if err != nil {
		l.Log().Error("Invalid opencv cascades folder", "folder", *cascadeFolder, "error", err)
		os.Exit(1)
	}
	f.Close()

	h := handlers.NewHealth(l)
	fd := handlers.NewPost(*cascadeFolder)

	r := mux.NewRouter()
	r.Handle("/health", h).Methods("GET")
	r.Handle("/", fd).Methods("POST")

	if *profile {
		// Register pprof handlers
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

		r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		r.Handle("/debug/pprof/block", pprof.Handler("block"))
	}

	l.ServiceStart(*listenAddress, *listenPort)
	fmt.Println("Error starting server", "error", http.ListenAndServe(fmt.Sprintf("%s:%s", *listenAddress, *listenPort), r))
}
