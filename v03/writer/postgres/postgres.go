// Package postgres implements the writer logic for writing data to postgres
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/v03/core"
	postgresCore "github.com/tapglue/backend/v03/core/postgres"
	ksis "github.com/tapglue/backend/v03/storage/kinesis"
	"github.com/tapglue/backend/v03/storage/postgres"
	"github.com/tapglue/backend/v03/writer"
)

// This allows us to process at most x number of entries per stream
// TODO Does it make sense to have something different? Investigate how this will behave in production
const maxEntriesPerStream = 50

type pg struct {
	ksis            ksis.Client
	pg              postgres.Client
	organization    core.Organization
	member          core.Member
	application     core.Application
	applicationUser core.ApplicationUser
	connection      core.Connection
	event           core.Event
}

const (
	getConsumerPositionQuery    = `SELECT consumer_position FROM tg.consumers WHERE consumer_name='distributor'`
	updateConsumerPositionQuery = `UPDATE tg.consumers SET consumer_position=$1, updated_at=now() WHERE consumer_name='distributor'`
	insertConsumerInitialQuery  = `INSERT INTO tg.consumers(consumer_name, consumer_position, updated_at) VALUES ('distributor', '', now())`
)

var (
	errUnknownMessage     = errors.New(http.StatusInternalServerError, 1, "unknown message retrieved", "", false)
	hostname, hostnameErr = os.Hostname()
)

func init() {
	if hostnameErr != nil {
		fmt.Println("failed to fetcht the hostname")
		panic(hostnameErr)
	}
}

func (p *pg) consumeStream(streamName, position string) {
	output, sequenceNumber, errs, done := p.ksis.StreamRecords(streamName, fmt.Sprintf("%s-%s", hostname, streamName), position, maxEntriesPerStream)
	internalDone := make(chan struct{}, 2)

	// We want to process errors in background
	go p.processErrors(streamName, errs, internalDone)

	go p.progressSaver(sequenceNumber, internalDone)

	go p.processMessages(output, errs, internalDone)

	<-done
	internalDone <- struct{}{}
	internalDone <- struct{}{}
	time.After(1 * time.Second)
}

func (p *pg) Execute(env string, mainLogChan, errorLogChan chan *logger.LogMsg) {
	/*for idx := range kinesis.Streams {
		go p.consumeStream(kinesis.Streams[idx])
	}*/

	streamName := ""
	switch env {
	case "dev":
		streamName = ksis.PackedStreamNameDev
	case "test":
		streamName = ksis.PackedStreamNameTest
	case "prod":
		streamName = ksis.PackedStreamNameProduction
	}

	consumerPosition := ""
	if err := p.pg.SlaveDatastore(-1).QueryRow(getConsumerPositionQuery).Scan(&consumerPosition); err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		} else {
			_, err = p.pg.MainDatastore().Exec(insertConsumerInitialQuery)
			if err != nil {
				panic(err)
			}
		}
	}

	p.consumeStream(streamName, consumerPosition)
}

func (p *pg) processErrors(streamName string, errs chan errors.Error, internalDone chan struct{}) {
	func() {
		for {
			select {
			case err, ok := <-errs:
				if !ok {
					return
				}
				log.Printf("ERROR\t%s\t%s", streamName, err.InternalErrorWithLocation())
			case <-internalDone:
				return
			}
		}
	}()
}

func (p *pg) progressSaver(sequenceNumber <-chan string, internalDone chan struct{}) {
	var (
		sequenceNo, prevPosition = "", ""
		ok                       bool
	)

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sequenceNo, ok = <-sequenceNumber:
			if !ok {
				break
			}
			log.Printf("GOT SEQUENCE:\t%s", sequenceNo)

		case <-ticker.C:
			sequence := sequenceNo
			if prevPosition == sequence {
				break
			}
			log.Printf("SAVING SEQUENCE:\t%s", sequence)
			_, err := p.pg.MainDatastore().Exec(updateConsumerPositionQuery, sequence)
			if err != nil {
				log.Printf("Error %s while saving sequence: %s", err, sequence)
			}
		case <-internalDone:
			return
		}
	}
}

func (p *pg) processMessages(output <-chan string, errs chan errors.Error, internalDone chan struct{}) {
	for {
		select {
		case msg, ok := <-output:
			{
				if !ok {
					return
				}

				channelName, msg, err := p.ksis.UnpackRecord(msg)
				if err != nil {
					errs <- err
					break
				}

				var ers []errors.Error
				switch channelName {
				case ksis.StreamAccountUpdate:
					ers = p.organizationUpdate(msg)
				case ksis.StreamAccountDelete:
					ers = p.organizationDelete(msg)
				case ksis.StreamAccountUserCreate:
					ers = p.memberCreate(msg)
				case ksis.StreamAccountUserUpdate:
					ers = p.memberUpdate(msg)
				case ksis.StreamAccountUserDelete:
					ers = p.memberDelete(msg)
				case ksis.StreamApplicationCreate:
					ers = p.applicationCreate(msg)
				case ksis.StreamApplicationUpdate:
					ers = p.applicationUpdate(msg)
				case ksis.StreamApplicationDelete:
					ers = p.applicationDelete(msg)
				case ksis.StreamApplicationUserUpdate:
					ers = p.applicationUserUpdate(msg)
				case ksis.StreamApplicationUserDelete:
					ers = p.applicationUserDelete(msg)
				case ksis.StreamConnectionCreate:
					ers = p.connectionCreate(msg)
				case ksis.StreamConnectionConfirm:
					ers = p.connectionConfirm(msg)
				case ksis.StreamConnectionUpdate:
					ers = p.connectionUpdate(msg)
				case ksis.StreamConnectionAutoConnect:
					ers = p.connectionAutoConnect(msg)
				case ksis.StreamConnectionSocialConnect:
					ers = p.connectionSocialConnect(msg)
				case ksis.StreamConnectionDelete:
					ers = p.connectionDelete(msg)
				case ksis.StreamEventCreate:
					ers = p.eventCreate(msg)
				case ksis.StreamEventUpdate:
					ers = p.eventUpdate(msg)
				case ksis.StreamEventDelete:
					ers = p.eventDelete(msg)
				default:
					{
						errs <- errUnknownMessage.UpdateInternalMessage(msg)
					}
				}

				// TODO we should really do something with the error here, not just ignore it like this, maybe?
				if ers != nil {
					for idx := range ers {
						errs <- ers[idx].UpdateInternalMessage(ers[idx].InternalErrorWithLocation() + "\t" + msg)
					}
				} else {
					log.Printf("processed message\t%s\t%s", channelName, msg)
				}
			}
		case <-internalDone:
			return
		}
	}
}

// New will return a new PosgreSQL writer
func New(kinesis ksis.Client, pgsql postgres.Client) writer.Writer {
	return &pg{
		ksis:            kinesis,
		pg:              pgsql,
		organization:    postgresCore.NewOrganization(pgsql),
		member:          postgresCore.NewMember(pgsql),
		application:     postgresCore.NewApplication(pgsql),
		applicationUser: postgresCore.NewApplicationUser(pgsql),
		connection:      postgresCore.NewConnection(pgsql),
		event:           postgresCore.NewEvent(pgsql),
	}
}
