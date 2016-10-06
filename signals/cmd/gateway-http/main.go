package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"

	"github.com/tapglue/multiverse/platform/generate"
	"github.com/tapglue/multiverse/signals/pb"
)

const (
	envDataset = "BIGQUERY_DATASET"
	envProject = "GCLOUD_PROJECT"
	envTopic   = "PUBSUB_TOPIC"
)

type bqSignalRaw struct {
	ID  uint64
	Raw []byte
}

func main() {
	var (
		ctx       = context.Background()
		dataset   = mustGetenv(envDataset)
		projectID = mustGetenv(envProject)
		topic     = mustGetenv(envTopic)
	)

	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	ds := bq.Dataset(dataset)

	ps, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	t, err := ps.CreateTopic(ctx, topic)
	if err != nil {
		if gErr, ok := err.(*googleapi.Error); !ok || gErr.Code != 409 {
			log.Fatal(err)
		}
	}

	schema, err := bigquery.InferSchema(bqSignalRaw{})
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range schema {
		log.Println(f.Name, f.Type)
	}

	ts := map[string]*bigquery.Table{}

	http.HandleFunc("/_ah/health", healthCheckHandler)
	http.HandleFunc(
		"/pubsub-handlers/signals-persist-bigquery",
		persistBigQueryHandler(
			ds,
			tableForNamespace(ts),
		),
	)
	http.HandleFunc("/track/signal", trackSignalHandler(t))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func persistBigQueryHandler(ds *bigquery.Dataset, tableForNamespace tableForNamespaceFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = context.Background()
			msg = struct {
				Message struct {
					Attributes map[string]string
					Data       []byte
					ID         string `json:"message_id"`
				}
				Subscription string
			}{}
			signal = &pb.Signal{}
		)

		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Errorf("message decode failed: %s", err))
			return
		}

		if err := proto.Unmarshal(msg.Message.Data, signal); err != nil {
			respondError(w, http.StatusBadRequest, fmt.Errorf("signal decode failed: %s", err))
			return
		}

		_, err := tableForNamespace(ctx, ds, signal.Namespace)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Errorf("table construction failed: %s", err))
			return
		}

		// TODO: store in BigQuery
		// TODO: cache id to avoid multiple appends of the same Signal (Memcache?)

		respondJSON(w, http.StatusNoContent, nil)
	}
}

func trackSignalHandler(topic *pubsub.Topic) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = context.Background()
			signal = &pb.Signal{
				Arrvied:   time.Now().Format(time.RFC3339Nano),
				Id:        generate.RandomString(32),
				Org:       "1",
				App:       "1",
				Namespace: "1_1",
			}
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

type apiError struct {
	Message string `json:"message"`
}

func mustGetenv(key string) string {
	v := os.Getenv(key)

	if v == "" {
		log.Fatalf("%s env variable not set", key)
	}

	return v
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

type tableForNamespaceFunc func(context.Context, *bigquery.Dataset, string) (*bigquery.Table, error)

func tableForNamespace(ts map[string]*bigquery.Table) tableForNamespaceFunc {
	return func(
		ctx context.Context,
		ds *bigquery.Dataset,
		ns string,
	) (*bigquery.Table, error) {
		t, ok := ts[ns]
		if ok {
			return t, nil
		}

		// If the table is not already been stored we should set it up locally
		// and if necessary on BigQuery itself.
		t = ds.Table(ns)

		if _, err := t.Metadata(ctx); err != nil {
			if gErr, ok := err.(*googleapi.Error); !ok || gErr.Code != 404 {
				return nil, fmt.Errorf("metadata fetch failed: %s", err)
			}

			err := t.Create(ctx, bigquery.UseStandardSQL())
			if err != nil {
				return nil, fmt.Errorf("create failed: %s", err)
			}
		}

		ts[ns] = t

		return t, nil
	}
}
