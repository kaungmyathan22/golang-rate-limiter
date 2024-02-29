package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status: "failed",
				Body:   "The api is at capacity, try again later.",
			}
			writer.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(writer).Encode(&message)
			return
		} else {
			next(writer, request)
		}
	})
	// next(w, r)
}
