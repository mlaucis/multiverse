package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/device"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/session"
	"github.com/tapglue/multiverse/service/user"
)

var enabled = true

func main() {
	var (
		namespace   = flag.String("namespace", "", "Namespace to dump data from")
		postgresURL = flag.String("postgres.url", "", "Postgres URL to connect to")
	)
	flag.Parse()

	db, err := sqlx.Connect("postgres", *postgresURL)
	if err != nil {
		log.Fatal(err)
	}

	var (
		apps        = app.NewPostgresService(db)
		connections = connection.NewPostgresService(db)
		devices     = device.PostgresService(db)
		events      = event.NewPostgresService(db)
		objects     = object.NewPostgresService(db)
		sessions    = session.NewPostgresService(db)
		users       = user.NewPostgresService(db)
		ns          = *namespace
	)

	err = os.MkdirAll(ns, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpApp(apps, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpConnections(connections, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpDevices(devices, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpEvents(events, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpObjects(objects, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpSessions(sessions, ns)
	if err != nil {
		log.Fatal(err)
	}

	err = dumpUsers(users, ns)
	if err != nil {
		log.Fatal(err)
	}
}

func dumpApp(apps app.Service, ns string) error {
	ps := strings.SplitN(ns, "_", 3)

	if len(ps) != 3 {
		return fmt.Errorf("Invalid namespace: %s", ns)
	}

	id, err := strconv.ParseUint(ps[2], 10, 64)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/apps.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	as, err := apps.Query(app.NamespaceDefault, app.QueryOptions{
		Enabled: &enabled,
		IDs: []uint64{
			id,
		},
	})
	if err != nil {
		return err
	}

	for _, a := range as {
		err := out.Encode(a)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpConnections(connections connection.Service, ns string) error {
	cs, err := connections.Query(ns, connection.QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/connections.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, con := range cs {
		err := out.Encode(con)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpDevices(devices device.Service, ns string) error {
	ds, err := devices.Query(ns, device.QueryOptions{})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/devices.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, d := range ds {
		err := out.Encode(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpEvents(events event.Service, ns string) error {
	es, err := events.Query(ns, event.QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/events.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, ev := range es {
		err := out.Encode(ev)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpObjects(objects object.Service, ns string) error {
	ls, err := objects.Query(ns, object.QueryOptions{})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/objects.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, o := range ls {
		err := out.Encode(o)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpSessions(sessions session.Service, ns string) error {
	ss, err := sessions.Query(ns, session.QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/sessions.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, s := range ss {
		err := out.Encode(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func dumpUsers(users user.Service, ns string) error {
	us, err := users.Query(ns, user.QueryOptions{
		Enabled: &enabled,
	})
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("./%s/users.json", ns))
	if err != nil {
		return err
	}

	out := json.NewEncoder(f)

	for _, u := range us {
		err := out.Encode(u)
		if err != nil {
			return err
		}
	}

	return nil
}
