package main

import (
	"log"
	"net/http"

	"loadbalancer/loadbalancer"
)

func main() {
	backendURLs := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	lb := loadbalancer.NewLoadBalancer(backendURLs)

	server := &http.Server{
		Addr:    ":8080",
		Handler: lb,
	}

	log.Println("Load Balancer started at :8080")
	log.Fatal(server.ListenAndServe())
}
