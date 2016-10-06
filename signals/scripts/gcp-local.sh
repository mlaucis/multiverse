#!/bin/sh

export BIGQUERY_DATASET=signals_persist
export GCLOUD_PROJECT=tapglue-signals
export PUBSUB_TOPIC=signals-raw
$(gcloud beta emulators datastore env-init)
$(gcloud beta emulators pubsub env-init)

$@