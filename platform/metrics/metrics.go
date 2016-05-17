package metrics

import "time"

// BucketFormat is the date format used to bucket values in timeseries.
const BucketFormat = "2006-01-02"

// Field names for metric labels.
const (
	FieldComponent = "component"
	FieldMethod    = "method"
	FieldNamespace = "namespace"
	FieldRoute     = "route"
	FieldService   = "service"
	FieldSource    = "source"
	FieldStatus    = "status"
	FieldStore     = "store"
	FieldVersion   = "version"
)

// BucketsQueue are used for Histograms observing queue latencies.
var BucketsQueue = []float64{
	.0005,
	.001,
	.0025,
	.005,
	.01,
	.025,
	.05,
	.1,
	.25,
	.5,
	1,
}

// BucketByDay is expected to be satisfied by various entity implementations.
type BucketByDay interface {
	CreatedByDay(ns string, start, end time.Time) (Timeseries, error)
}

// Datapoint describes a point in a Timeseries carrying a bucket value.
type Datapoint struct {
	Bucket string `json:"bucket"`
	Value  int    `json:"value"`
}

// Timeseries is a collection of Datapoints.
type Timeseries []Datapoint
