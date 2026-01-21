package server

import (
	"net/http"
	"time"
)

const (
	healthEndpoint      = "http://localhost:8080/metrics/current"
	healthCheckTimeout  = 60 * time.Second
	healthCheckInterval = 2 * time.Second
	maxStartRetries     = 5
)

// checkHealthEndpoint makes a single HTTP GET request to the health endpoint
func checkHealthEndpoint() bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(healthEndpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// WaitForHealthy polls the health endpoint until it responds or timeout is reached
func WaitForHealthy() bool {
	deadline := time.Now().Add(healthCheckTimeout)
	for time.Now().Before(deadline) {
		if checkHealthEndpoint() {
			return true
		}
		time.Sleep(healthCheckInterval)
	}
	return false
}
