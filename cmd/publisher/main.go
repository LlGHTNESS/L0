package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "my-client2"

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Couldn't connect to NATS Streaming: %v", err)
	}
	defer sc.Close()

	files, err := filepath.Glob("../../misc/*.json")
	if err != nil {
		log.Fatalf("Error when searching for files: %v", err)
	}

	for _, file := range files {
		processJSONFile(file, sc)
	}
}

func processJSONFile(filename string, sc stan.Conn) {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Error reading from a file %s: %v", filename, err)
		return
	}

	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Printf("Error when deserializing JSON from a file %s: %v", filename, err)
		return
	}

	jsonData, err = json.Marshal(data)
	if err != nil {
		log.Printf("Error when serializing JSON from a file %s: %v", filename, err)
		return
	}

	subject := "main"
	if err := sc.Publish(subject, jsonData); err != nil {
		log.Printf("Error sending a message from a file %s: %v", filename, err)
		return
	}

	log.Printf("The message from the %s file has been successfully sent to NATS Streaming!", filename)
}
