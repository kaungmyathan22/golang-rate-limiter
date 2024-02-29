package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
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

func perClientLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip, _, err := net.SplitHostPort(request.RemoteAddr)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(2, 4),
			}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			message := Message{
				Status: "failed",
				Body:   "The api is at capacity, try again later.",
			}
			writer.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(writer).Encode(&message)
			return
		}
		mu.Unlock()
		next(writer, request)
	})
}

func main() {
	http.Handle("/ping", perClientLimiter(endpointHandler))
	log.Println("server is running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}

}
