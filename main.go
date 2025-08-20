package main

import (
	"fmt"
	"go-solar-client/endpoints"
)

func main() {
	apiUrl := "http://10.67.67.25:4545"
	health, err := endpoints.CheckHealth(apiUrl)
	if err != nil {
		fmt.Println("Error checking health:", err)
		return
	}
	fmt.Printf("Service: %s, Status: %s\n", health.Service, health.Status)
}
