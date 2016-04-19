package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/platform/metrics"
)

// AnalyticsDeprecated is the endpoint for devices to submit usage data.
func AnalyticsDeprecated() Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// Analytics is the endpoint for devices to submit usage data.
func Analytics() Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusNoContent, nil)
	}
}

// AnalyticsApp returns the analytics data for an app controlled by the
// passed where clause.
func AnalyticsApp(c *controller.AnalyticsController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		wr := whereWrapper{}

		if r.URL.Query().Get("where") != "" {
			err := json.Unmarshal([]byte(r.URL.Query().Get("where")), &wr)
			if err != nil {
				respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
				return
			}
		}

		rs, err := c.App(mux.Vars(r)["appID"], wr.where)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadAnalyticsApp{result: rs})
	}
}

type payloadAnalyticsApp struct {
	result controller.AppResult
}

func (p *payloadAnalyticsApp) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.result)
}

type whereWrapper struct {
	where *controller.AnalyticsWhere
	Start string `json:"start"`
	End   string `json:"end"`
}

func (w *whereWrapper) UnmarshalJSON(raw []byte) error {
	f := struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	start, err := time.Parse(metrics.BucketFormat, f.Start)
	if err != nil {
		return err
	}

	end, err := time.Parse(metrics.BucketFormat, f.End)
	if err != nil {
		return err
	}

	w.where = &controller.AnalyticsWhere{
		Start: start,
		End:   end,
	}

	return nil
}
