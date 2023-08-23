package main

import (
	"fmt"
	"os"
)

func main(){
	os.Setenv("API_KEY","l7u502p8v46ba3ppgvj5y2aad50lb9")

	easyBrokerClient := NewClient("https://api.stagingeb.com")

	// Test ListProperties for success
	err := easyBrokerClient.ListProperties()
	if err != nil {
		fmt.Println(fmt.Errorf("Unexpected error: %v", err).Error())
	}
}