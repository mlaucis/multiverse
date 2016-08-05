package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/bluele/slack"
	klog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
)

const (
	component = "reporter"
	sub       = "customers-weekly"
)

const (
	pgClauseByWeek            = `(current_date - interval '7 days')`
	pgClauseByMonth           = `(current_date - interval '30 days')`
	pgClauseMTD               = `cast(date_trunc('month', current_date) as date)`
	pgActiveUserCountByPeriod = `
    SELECT
    	count(DISTINCT userid)
    FROM (

    	SELECT DISTINCT
    		(json_data->>'user_from_id')::BIGINT AS userid
      FROM
    		{{.Namespace}}.connections
      WHERE
      	(json_data->>'updated_at')::DATE >= {{.Period}}
      UNION ALL

    	SELECT DISTINCT
    		(json_data->>'user_id')::BIGINT AS userid
      FROM
    		{{.Namespace}}.events
      WHERE
      	(json_data->>'updated_at')::DATE >= {{.Period}}
    	UNION ALL

    	SELECT DISTINCT
    		(json_data->>'owner_id')::BIGINT AS userid
    	FROM
    		{{.Namespace}}.objects
    	WHERE
    		(json_data->>'updated_at')::DATE >= {{.Period}}
    	UNION ALL

    	SELECT DISTINCT
    		user_id AS userid
      FROM
    		{{.Namespace}}.sessions
      WHERE
      	created_at >= {{.Period}}
      UNION ALL

    	SELECT DISTINCT
    		(json_data->>'id')::BIGINT AS userid
      FROM
    		{{.Namespace}}.users
      WHERE
      	(json_data->>'updated_at')::DATE >= {{.Period}} OR
          last_read >= {{.Period}}

    ) AS "count"`
)

var (
	currentRevision = "0000000-dev"
	defaultTrue     = true
)

var appNames = map[string]string{
	"app_26_187":  "DailyMe",
	"app_309_428": "Gambify Travis Target",
	"app_309_425": "Gambify Testing",
	"app_309_443": "Gambify Production",
	"app_374_501": "Stepz",
	"app_409_652": "DaWanda iOS",
}

type Report struct {
	AppName              string
	AverageConnections   float64
	AverageEvents        float64
	AverageObjects       float64
	ConnectionsTotal     uint64
	EventsTotal          uint64
	ObjectsTotal         uint64
	UsersActiveLastWeek  uint64
	UsersActiveLastMonth uint64
	UsersActiveMTD       uint64
	UsersTotal           uint64
}

func main() {
	var (
		hostname, _ = os.Hostname()
		startTime   = time.Now()

		pgURL        = flag.String("pg.url", "postgres://127.0.0.1:5432/test?sslmode=disable&connect_timeout=5", "Postgres URL")
		slackChannel = flag.String("slack.channel", "", "Slack channel to post reports to.")
		slackToken   = flag.String("slack.token", "", "Token used for authentication against the Slack API.")
	)
	flag.Parse()

	logger := klog.NewContext(
		klog.NewJSONLogger(os.Stdout),
	).With(
		"caller", klog.Caller(3),
		"component", component,
		"host", hostname,
		"revision", currentRevision,
		"sub", sub,
	)

	log.SetFlags(log.Ldate | log.Lshortfile)

	slackAPI := slack.New(*slackToken)
	channel, err := slackAPI.FindChannelByName(*slackChannel)
	if err != nil {
		fatalf(logger, "Slack channel lookup failed: %s", err)
	}

	db, err := sqlx.Connect("postgres", *pgURL)
	if err != nil {
		fatalf(logger, "Postgres connect failed: %s", err)
	}

	var apps app.Service
	apps = app.NewPostgresService(db)

	var connections connection.Service
	connections = connection.NewPostgresService(db)

	var events event.Service
	events = event.NewPostgresService(db)

	var objects object.Service
	objects = object.NewPostgresService(db)

	var users user.Service
	users = user.NewPostgresService(db)

	logger.Log(
		"duration", time.Now().Sub(startTime),
		"lifecycle", "start",
	)

	as, err := apps.Query(app.NamespaceDefault, app.QueryOptions{
		Enabled:      &defaultTrue,
		InProduction: &defaultTrue,
	})
	if err != nil {
		fatalf(logger, "App query failed: %s", err)
	}

	rs := []Report{}

	for _, a := range as {
		r := Report{
			AppName: a.Name,
		}

		if name, ok := appNames[a.Namespace()]; ok {
			r.AppName = name
		}

		allUsers, err := users.Count(a.Namespace(), user.QueryOptions{
			Enabled: &defaultTrue,
		})
		if err != nil {
			fatalf(logger, "User counting for %s failed: %s", r.AppName, err)
		}

		r.UsersTotal = uint64(allUsers)

		allConnections, err := connections.Count(a.Namespace(), connection.QueryOptions{
			Enabled: &defaultTrue,
		})
		if err != nil {
			fatalf(logger, "Connection counting for %s failed: %s", r.AppName, err)
		}

		r.AverageConnections = float64(allConnections) / float64(allUsers)
		r.ConnectionsTotal = uint64(allConnections)

		allEvents, err := events.Count(a.Namespace(), event.QueryOptions{
			Enabled: &defaultTrue,
		})
		if err != nil {
			fatalf(logger, "Event counting for %s failed: %s", r.AppName, err)
		}

		r.AverageEvents = float64(allEvents) / float64(allUsers)
		r.EventsTotal = uint64(allEvents)

		allObjects, err := objects.Count(a.Namespace(), object.QueryOptions{
			Deleted: false,
		})
		if err != nil {
			fatalf(logger, "Object counting for %s failed: %s", r.AppName, err)
		}

		r.AverageObjects = float64(allObjects) / float64(allUsers)
		r.ObjectsTotal = uint64(allObjects)

		r.UsersActiveLastWeek, err = countActive(db, a.Namespace(), pgClauseByWeek)
		if err != nil {
			fatalf(logger, "Active user counting for %s failed: %s", r.AppName, err)
		}

		r.UsersActiveLastMonth, err = countActive(db, a.Namespace(), pgClauseByMonth)
		if err != nil {
			fatalf(logger, "Active user counting for %s failed: %s", r.AppName, err)
		}

		r.UsersActiveMTD, err = countActive(db, a.Namespace(), pgClauseMTD)
		if err != nil {
			fatalf(logger, "Active user counting for %s failed: %s", r.AppName, err)
		}

		rs = append(rs, r)
	}

	err = reportToSlack(slackAPI, channel, rs)
	if err != nil {
		fatalf(logger, "Slack report failed: %s", err)
	}

	// for _, r := range rs {
	// send report to email
	// send report to s3
	// send report to spreadsheet
	// }

	logger.Log(
		"duration", time.Now().Sub(startTime),
		"lifecycle", "stop",
	)
}

func countActive(db *sqlx.DB, ns string, clause string) (uint64, error) {
	t, err := template.New("activeCountQuery").Parse(pgActiveUserCountByPeriod)
	if err != nil {
		return 0, err
	}

	var (
		buf = []byte{}
		w   = bytes.NewBuffer(buf)
	)

	err = t.Execute(w, struct {
		Namespace string
		Period    string
	}{
		Namespace: ns,
		Period:    clause,
	})
	if err != nil {
		return 0, err
	}

	var (
		active uint64
	)

	err = db.Get(&active, w.String())

	return active, err
}

func fatalf(logger *klog.Context, format string, vs ...interface{}) {
	logger.Log("err", fmt.Sprintf(format, vs...))
	os.Exit(1)
}

func reportToSlack(api *slack.Slack, channel *slack.Channel, rs []Report) error {
	var (
		year, week = time.Now().ISOWeek()
		msg        = fmt.Sprintf(
			"Hello <!channel>, reporting in with our weekly customer statistics for week %d in %d.",
			week,
			year,
		)
		opts = &slack.ChatPostMessageOpt{
			Attachments: []*slack.Attachment{},
			IconUrl:     "http://downloadpack.net/skin_pack/ironman/ironman.png",
			Username:    "Jarvis",
		}
	)

	for _, r := range rs {
		a := &slack.Attachment{
			Color: "#439FE0",
			Fields: []*slack.AttachmentField{
				{Title: "Last7", Value: fmt.Sprintf("%d", r.UsersActiveLastWeek), Short: true},
				{Title: "Last30", Value: fmt.Sprintf("%d", r.UsersActiveLastMonth), Short: true},
				{Title: "MTD", Value: fmt.Sprintf("%d", r.UsersActiveMTD), Short: true},
				{Title: "Users", Value: fmt.Sprintf("%d", r.UsersTotal), Short: true},
				{Title: "Connections", Value: fmt.Sprintf("%d", r.ConnectionsTotal), Short: true},
				{Title: "Events", Value: fmt.Sprintf("%d", r.EventsTotal), Short: true},
				{Title: "Objects", Value: fmt.Sprintf("%d", r.ObjectsTotal), Short: true},
				{Title: "Conn/User", Value: fmt.Sprintf("%f", r.AverageConnections), Short: true},
				{Title: "Event/User", Value: fmt.Sprintf("%f", r.AverageEvents), Short: true},
				{Title: "Object/User", Value: fmt.Sprintf("%f", r.AverageObjects), Short: true},
			},
			MarkdownIn: []string{
				"fields",
				"pretext",
				"text",
			},
			Pretext: fmt.Sprintf("*%s*", r.AppName),
		}

		opts.Attachments = append(opts.Attachments, a)
	}

	return api.ChatPostMessage(channel.Id, msg, opts)
}
