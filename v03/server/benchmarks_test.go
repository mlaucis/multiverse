// +build !kinesis
// +build postgres bench

package server_test

import (
	"flag"
	"log"

	. "gopkg.in/check.v1"
	"os"
)

var (
	organizations    = flag.Int("org", 50, "number of organizations")
	members          = flag.Int("orgmem", 5, "number of members per organization")
	applications     = flag.Int("app", 1, "number of applications per organization")
	applicationUsers = flag.Int("appusr", 100000, "number of application users")
	connections      = flag.Int("conn", 10, "number of connections per application user (max 50 or appusr)")
	events           = flag.Int("evt", 500, "number of events per application user")
)

func (s *EventSuite) TestCreateUsers(c *C) {
	flag.Parse()

	if os.Getenv("CI") == "true" {
		c.Skip("not to be run inside the CI suite")
	}

	if *connections > *applicationUsers {
		connections = applicationUsers
	}

	if *connections > 50 {
		*connections = 50
	}

	log.Printf("\n\nCreating organizations = %d, members = %d, applications = %d, applicationUsers = %d, connections = %d, events = %d\n\n",
		*organizations, *members, *applications, *applicationUsers, *connections, *events)

	CorrectDeployBench(*organizations, *members, *applications, *applicationUsers, *connections, *events, true, true)
}
