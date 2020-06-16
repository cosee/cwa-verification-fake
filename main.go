package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	envDelayMillis = "CWA_FAKE_DELAY_MILLIS"
	envIp = "CWA_FAKE_IP"
	envPort = "CWA_FAKE_PORT"
)

type responseBody struct {
	Tan string
}

func (resp responseBody) isValid(validtans []string) bool {
	for _, t := range validtans {
		if t == resp.Tan {
			return true
		}
	}
	return false
}

func loadDelayInMillis() int64 {
	delayAsString := os.Getenv(envDelayMillis)
	if delayAsString != "" {
		delayMillis, err := strconv.ParseInt(delayAsString, 10, 64)
		if err != nil {
			log.Fatalf("Cannot parse '%s' of '%s' as int.", envDelayMillis, delayAsString)
		}
		return delayMillis
	}
	return 0
}

func main() {
	validTans := []string{"edc07f08-a1aa-11ea-bb37-0242ac130002"}

	serverAddress := os.Getenv(envIp)
	if serverAddress == "" {
		serverAddress = "0.0.0.0"
	}

	port := os.Getenv(envPort)
	if port == "" {
		port = "8004"
	}

	delayMillis := loadDelayInMillis()

	http.HandleFunc("/version/v1/tan/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			var resp responseBody

			if delayMillis > 0 {
				time.Sleep(time.Duration(delayMillis) * time.Millisecond)
			}

			if err := decoder.Decode(&resp); err != nil {
				http.Error(w, err.Error(), 404)
			}
			log.Printf("%+v", resp)
			if resp.isValid(validTans) {
				w.WriteHeader(200)
			} else {
				http.Error(w, "wrong tan", 404)
			}
		} else {
			http.Error(w, "Only POST is supported", 404)
		}
	})

	log.Printf("listening on %s, port %s, with delay in millis: %v", serverAddress, port, delayMillis)
	log.Fatal(http.ListenAndServe(serverAddress+":"+port, nil)) //if serverAddress is "", listen all
}
