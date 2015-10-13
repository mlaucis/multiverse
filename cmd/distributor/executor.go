package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/logger"
)

// This allows us to process at most x number of entries per stream
// TODO Does it make sense to have something different? Investigate how this will behave in production
const maxEntriesPerStream = 50

const (
	getConsumerPositionQuery    = `SELECT consumer_position FROM tg.consumers WHERE consumer_name='distributor'`
	updateConsumerPositionQuery = `UPDATE tg.consumers SET consumer_position=$1, updated_at=now() WHERE consumer_name='distributor'`
	insertConsumerInitialQuery  = `INSERT INTO tg.consumers(consumer_name, consumer_position, updated_at) VALUES ('distributor', '', now())`
)

func consumeStream(streamName, position string) {
	output, sequenceNumber, errs, done := ksis.StreamRecords(streamName, fmt.Sprintf("%s-%s", hostname, streamName), position, maxEntriesPerStream)
	internalDone := make(chan struct{}, 2)

	// We want to process errors in background
	go processErrors(streamName, errs, internalDone)

	go progressSaver(sequenceNumber, internalDone)

	go processMessages(output, errs, internalDone)

	<-done
	internalDone <- struct{}{}
	internalDone <- struct{}{}
	time.After(1 * time.Second)
}

func processErrors(streamName string, errs chan errors.Error, internalDone chan struct{}) {
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

func progressSaver(sequenceNumber <-chan string, internalDone chan struct{}) {
	var (
		sequenceNo = ""
		ok         bool
	)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case sequenceNo, ok = <-sequenceNumber:
			if !ok {
				break
			}
			log.Printf("GOT SEQUENCE:\t%s", sequenceNo)

		case <-ticker.C:
			if sequenceNo == "" {
				continue
			}
			log.Printf("SAVING SEQUENCE:\t%s", sequenceNo)
			_, err := pg.MainDatastore().Exec(updateConsumerPositionQuery, sequenceNo)
			if err != nil {
				log.Printf("Error %s while saving sequence: %s", err, sequenceNo)
			}
		case <-internalDone:
			return
		}
	}
}

func processMessages(output <-chan string, errs chan errors.Error, internalDone chan struct{}) {
	for {
		select {
		case msg, ok := <-output:
			{
				if !ok {
					return
				}

				channelName, msg, err := ksis.UnpackRecord(msg)
				if err != nil {
					errs <- err
					break
				}

				var ers errors.Error
				if strings.HasPrefix(channelName, "v02_") {
					ers = v02PgConsumer.ProcessMessages(channelName, msg)
				} else if strings.HasPrefix(channelName, "v03_") {
					ers = v03PgConsumer.ProcessMessages(channelName, msg)
				} else {

				}

				// TODO we should really do something with the error here, not just ignore it like this, maybe?
				if ers != nil {
					errs <- ers
				} else {
					log.Printf("processed message\t%s\t%s", channelName, msg)
				}
			}
		case <-internalDone:
			return
		}
	}
}

func execute(streamName string, mainLogChan, errorLogChan chan *logger.LogMsg) {
	consumerPosition := ""
	if err := pg.SlaveDatastore(-1).QueryRow(getConsumerPositionQuery).Scan(&consumerPosition); err != nil {
		if err != sql.ErrNoRows {
			panic(err)
		} else {
			_, err = pg.MainDatastore().Exec(insertConsumerInitialQuery)
			if err != nil {
				panic(err)
			}
		}
	}

	consumeStream(streamName, consumerPosition)
}
