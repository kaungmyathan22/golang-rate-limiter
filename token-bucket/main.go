package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	message := Message{
		Status: "success",
		Body:   "Hi!, You're reached the API.",
	}
	err := json.NewEncoder(writer).Encode(&message)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.Handle("/ping", rateLimiter(endpointHandler))
	log.Println("server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}
