package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"

	"github.com/tapglue/multiverse/signals/pb"
)

const (
	envProject = "GCLOUD_PROJECT"
	envTopic   = "PUBSUB_TOPIC"
)

func main() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, mustGetenv(envProject))
	if err != nil {
		log.Fatal(err)
	}

	topic, err := client.CreateTopic(ctx, mustGetenv(envTopic))
	if err != nil {
		if gErr, ok := err.(*googleapi.Error); !ok || gErr.Code != 409 {
			log.Fatal(err)
		}
	}

	http.HandleFunc("/_ah/health", healthCheckHandler)
	http.HandleFunc("/track/event", trackEventHandler(topic))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func trackEventHandler(topic *pubsub.Topic) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = context.Background()
			signal = &pb.Signal{}
		)

		raw, err := proto.Marshal(signal)
		if err != nil {
			respondError(w, 500, err)
			return
		}

		if _, err := topic.Publish(ctx, &pubsub.Message{Data: raw}); err != nil {
			respondError(w, 500, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

func mustGetenv(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatalf("%s env variable not set", key)
	}

	return v
}

type apiError struct {
	Message string `json:"message"`
}

func respondError(w http.ResponseWriter, code int, err error) {
	respondJSON(w, code, apiError{
		Message: err.Error(),
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
