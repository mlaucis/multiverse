package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/api/googleapi"

	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/signals/pb"
)

const (
	envDataset      = "BIGQUERY_DATASET"
	envMemcacheHost = "MEMCACHE_PORT_11211_TCP_ADDR"
	envMemcachePort = "MEMCACHE_PORT_11211_TCP_PORT"
	envProject      = "GCLOUD_PROJECT"
	envTopic        = "PUBSUB_TOPIC"
)

func main() {
	var (
		ctx          = context.Background()
		dataset      = mustGetenv(envDataset, "signals_persist")
		memcacheHost = mustGetenv(envMemcacheHost, "localhost")
		memcachePort = mustGetenv(envMemcachePort, "11211")
		projectID    = mustGetenv(envProject, "tapglue-signals")
		topic        = mustGetenv(envTopic, "signals-raw")
	)

	bq, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	ds := bq.Dataset(dataset)

	mc := memcache.New(fmt.Sprintf("%s:%s", memcacheHost, memcachePort))

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

	http.HandleFunc("/_ah/health", healthCheckHandler)
	http.HandleFunc(
		"/pubsub-handlers/signals-persist-bigquery-raw",
		persistBigRawQueryHandler(
			ds,
			mc,
			uploaderForNamespace(map[string]*bigquery.Uploader{}),
		),
	)
	http.HandleFunc("/track/signal", trackSignalHandler(t))

	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func persistBigQueryRawHandler(
	ds *bigquery.Dataset,
	mc *memcache.Client,
	uploader uploaderForNamespaceFunc,
) http.HandlerFunc {
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
			log.Println(fmt.Errorf("message decode failed: %s", err))

			respondError(w, http.StatusBadRequest, fmt.Errorf("message decode failed: %s", err))
			return
		}

		if _, err := mc.Get(fmt.Sprintf("signals-persist-raw|%s", msg.Message.ID)); err != memcache.ErrCacheMiss {
			if err == nil {
				respondJSON(w, http.StatusNoContent, nil)
				return
			}

			log.Println(fmt.Errorf("memcache lookup failed: %s", err))

			respondError(w, http.StatusBadRequest, fmt.Errorf("message lookup failed: %s", err))
			return
		}

		if err := proto.Unmarshal(msg.Message.Data, signal); err != nil {
			log.Println(fmt.Errorf("signal decode failed: %s", err))

			respondError(w, http.StatusBadRequest, fmt.Errorf("signal decode failed: %s", err))
			return
		}

		u, err := uploader(ctx, ds, signal.Namespace)
		if err != nil {
			log.Println(fmt.Errorf("table construction failed: %s", err))

			respondError(w, http.StatusInternalServerError, fmt.Errorf("uploader retrieval failed: %s", err))
			return
		}

		v := &bqSignalRaw{
			ID:  signal.Id,
			Raw: msg.Message.Data,
		}

		v.Arrived, err = time.Parse(time.RFC3339Nano, signal.Arrvied)
		if err != nil {
			log.Println(fmt.Errorf("arrived parse failed: %s", err))

			respondError(w, http.StatusInternalServerError, fmt.Errorf("arrived parse failed: %s", err))
			return
		}

		if err := u.Put(ctx, v); err != nil {
			log.Println(fmt.Errorf("bq append failed: %s", err))

			respondError(w, http.StatusInternalServerError, fmt.Errorf("bq append failed: %s", err))
			return
		}

		if err := mc.Add(&memcache.Item{
			Expiration: 604800, // Expire in 7 days.
			Key:        fmt.Sprintf("signals-persist-raw|%s", msg.Message.ID),
			Value:      []byte("processed"),
		}); err != nil {
			log.Println(fmt.Errorf("memcache add failed: %s", err))

			respondError(w, http.StatusInternalServerError, fmt.Errorf("memcache add failed: %s", err))
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

func trackSignalHandler(topic *pubsub.Topic) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    = context.Background()
			signal = &pb.Signal{
				Arrvied:   time.Now().Format(time.RFC3339Nano),
				Org:       "1",
				App:       "1",
				Namespace: "1_1",
			}
		)

		id, err := flake.NextID("signals")
		if err != nil {
			respondError(w, 500, fmt.Errorf("id creation failed: %s", err))
			return
		}

		signal.Id = id

		raw, err := proto.Marshal(signal)
		if err != nil {
			respondError(w, 500, fmt.Errorf("marshal failed: %s", err))
			return
		}

		if _, err := topic.Publish(ctx, &pubsub.Message{Data: raw}); err != nil {
			respondError(w, 500, fmt.Errorf("publish failed: %s", err))
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

type apiError struct {
	Message string `json:"message"`
}

type bqSignalRaw struct {
	Arrived time.Time
	ID      uint64
	Raw     []byte
}

func (v *bqSignalRaw) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"arrived": v.Arrived,
		"id":      v.ID,
		"raw":     v.Raw,
	}, strconv.FormatUint(v.ID, 10), nil
}

func mustGetenv(key, def string) string {
	v := os.Getenv(key)

	if v == "" {
		if def != "" {
			return def
		}

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

type uploaderForNamespaceFunc func(context.Context, *bigquery.Dataset, string) (*bigquery.Uploader, error)

func uploaderForNamespace(um map[string]*bigquery.Uploader) uploaderForNamespaceFunc {
	return func(
		ctx context.Context,
		ds *bigquery.Dataset,
		ns string,
	) (*bigquery.Uploader, error) {
		u, ok := um[ns]
		if ok {
			return u, nil
		}

		// If the Uploader is not present we need to set it up including
		// table creation if not present.
		t := ds.Table(ns)

		if _, err := t.Metadata(ctx); err != nil {
			if gErr, ok := err.(*googleapi.Error); !ok || gErr.Code != 404 {
				return nil, fmt.Errorf("metadata fetch failed: %s", err)
			}

			err := t.Create(ctx, bigquery.UseStandardSQL())
			if err != nil {
				return nil, fmt.Errorf("create failed: %s", err)
			}
		}

		// In order to make append operations succeed we need to set the schema
		// upfront.

		p := t.Patch()

		s := bigquery.Schema{
			&bigquery.FieldSchema{
				Name:        "arrived",
				Description: "signal arrival in system",
				Required:    true,
				Type:        bigquery.TimestampFieldType,
			},
			&bigquery.FieldSchema{
				Name:        "id",
				Description: "signal unique identifier",
				Required:    true,
				Type:        bigquery.IntegerFieldType,
			},
			&bigquery.FieldSchema{
				Name:        "raw",
				Description: "raw signal payload",
				Required:    true,
				Type:        bigquery.StringFieldType,
			},
		}

		p.Schema(s)

		if _, err := p.Apply(ctx); err != nil {
			return nil, fmt.Errorf("schema apply failed: %s", err)
		}

		u = t.NewUploader()
		um[ns] = u

		return u, nil
	}
}
