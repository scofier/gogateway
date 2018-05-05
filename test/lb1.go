package main

// Package goproxy is a LoadBalancer based on httputil.ReverseProxy.
//
// ExtractNameVersion and LoadBalance can be overridden in order to customize
// the behavior.

import (

"fmt"
"math/rand"
"net"

)


// loadBalance is a basic loadBalancer which randomly
// tries to connect to one of the endpoints and try again
// in case of failure.
func loadBalance(network, serviceName, serviceVersion string, reg Registry) (net.Conn, error) {
	newSlice, err := reg.Lookup(serviceName, serviceVersion)
	if err != nil {
		return nil, err
	}
	endpoints := make([]string, len(newSlice))
	copy(endpoints, newSlice)
	for {
		// No more endpoint, stop
		if len(endpoints) == 0 {
			break
		}
		// Select a random endpoint
		i := rand.Int() % len(endpoints)
		endpoint := endpoints[i]
		// Try to connect
		conn, err := net.Dial(network, endpoint)
		if err != nil {
			reg.Failure(serviceName, serviceVersion, endpoint, err)
			// Failure: remove the endpoint from the current list and try again.
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		endpoints = nil
		// Success: return the connection.
		return conn, nil
	}
	// No available endpoint.
	return nil, fmt.Errorf("No endpoint available for %s/%s", serviceName, serviceVersion)
}

