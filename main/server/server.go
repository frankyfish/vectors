package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// create constructor

type VectorServer struct {
	port int
}

type SemanticSearchRequest struct {
	Text string
}

func NewVectorServer(port int) *VectorServer {
	return &VectorServer{
		port: port,
	}
}

func (vs *VectorServer) PrepareAndStart(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", vs.port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServe())
}

func EnableCors(w http.ResponseWriter) {
	// w.Header().Set("Access-Control-Allow-Origin", "*") // uh oh
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
