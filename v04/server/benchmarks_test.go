// +build postgres bench

package server_test

import (
	"flag"
	"log"
	"os"

	. "gopkg.in/check.v1"
)

var (
	organizations    = flag.Int("org", 50, "number of organizations")
	organizationID   = flag.Int64("orgid", 1, "organization id")
	members          = flag.Int("orgmem", 5, "number of members per organization")
	applications     = flag.Int("app", 1, "number of applications per organization")
	applicationID    = flag.Int64("appid", 1, "application id")
	applicationUsers = flag.Int("appusr", 100000, "number of application users")
	connections      = flag.Int("conn", 10, "number of connections per application user (max 50 or appusr)")
	events           = flag.Int("evt", 500, "number of events per application user")
)

func (s *BenchSuite) TestCreateUsers(c *C) {
	t.SkipNow()
	if !flag.Parsed() {
		flag.Parse()
	}

	if os.Getenv("CI") == "true" {
		c.Skip("not to be run inside the CI suite")
	}

	if os.Getenv("CORRECT_DEPLOY") != "true" {
		c.Skip("this can run only under CORRECT_DEPLOY=true")
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

func (s *BenchSuite) TestCreateUserEvents(c *C) {
	t.SkipNow()
	if !flag.Parsed() {
		flag.Parse()
	}

	if os.Getenv("CI") == "true" {
		c.Skip("not to be run inside the CI suite")
	}

	if os.Getenv("CORRECT_DEPLOY") != "true" {
		c.Skip("this can run only under CORRECT_DEPLOY=true")
	}

	log.Printf("\n\nCreating events for organization = %d, application = %d\n\n",
		*organizationID, *applicationID)

	populateEventsForUsers(*organizationID, *applicationID, *events)
}
