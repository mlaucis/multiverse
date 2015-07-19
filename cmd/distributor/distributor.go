/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Command distributor will launch a specified consumer for Kinesis and write the received information to its target
//
// Currently it supports:
// - postgres
package main

import (
	"flag"
	"log"
	"log/syslog"
	mr "math/rand"
	"os"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"
	"github.com/tapglue/backend/v02/writer"
	"github.com/tapglue/backend/v02/writer/postgres"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_DISTRIBUTOR_CONFIG_PATH"
)

var (
	startTime       time.Time
	conf            *config.Config
	currentRevision string

	consumerTarget = flag.String("target", "", "Select the target of the consumer to be launched. Currently supported: postgres")
	myConsumer     writer.Writer
	pgConsumer     writer.Writer

	mainLogChan  = make(chan *logger.LogMsg, 100000)
	errorLogChan = make(chan *logger.LogMsg, 100000)
)

func init() {
	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

	log.SetFlags(0)

	if conf.UseSysLog {
		syslogWriter, err := syslog.New(syslog.LOG_INFO, "intaker")
		if err == nil {
			log.Printf("logging to syslog is enabled. Please tail your syslog for intaker app for further logs\n")
			log.SetOutput(syslogWriter)
		} else {
			log.Printf("%v\n", err)
			log.Printf("logging to syslog failed reverting to stdout logging\n")
		}
		conf.UseArtwork = false
	}

	errors.Init(conf.Environment != "prod")

	var v02KinesisClient v02_kinesis.Client
	if conf.Environment == "prod" {
		v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment)
	} else {
		if conf.Kinesis.Endpoint != "" {
			v02KinesisClient = v02_kinesis.NewTest(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment)
		} else {
			v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment)
		}
	}

	v02PgClient := v02_postgres.New(conf.Postgres)

	pgConsumer = postgres.New(v02KinesisClient, v02PgClient)
}

func main() {
	flag.Parse()

	if *consumerTarget == "" {
		flag.PrintDefaults()
		os.Exit(64)
	}

	switch *consumerTarget {
	case "postgres":
		myConsumer = pgConsumer
	default:
		flag.PrintDefaults()
		os.Exit(64)
	}

	if conf.JSONLogs {
		go logger.JSONLog(mainLogChan)
		go logger.JSONLog(errorLogChan)
	} else {
		go logger.TGLog(mainLogChan)
		go logger.TGLog(errorLogChan)
	}

	if conf.UseArtwork {
		log.Printf(`

88888888888                         888                        8888888b.  d8b          888            d8b 888               888
    888                             888                        888  "Y88b Y8P          888            Y8P 888               888
    888                             888                        888    888              888                888               888
    888   8888b.  88888b.   .d88b.  888 888  888  .d88b.       888    888 888 .d8888b  888888 888d888 888 88888b.  888  888 888888 .d88b.  888d888
    888      "88b 888 "88b d88P"88b 888 888  888 d8P  Y8b      888    888 888 88K      888    888P"   888 888 "88b 888  888 888   d88""88b 888P"
    888  .d888888 888  888 888  888 888 888  888 88888888      888    888 888 "Y8888b. 888    888     888 888  888 888  888 888   888  888 888
    888  888  888 888 d88P Y88b 888 888 Y88b 888 Y8b.          888  .d88P 888      X88 Y88b.  888     888 888 d88P Y88b 888 Y88b. Y88..88P 888
    888  "Y888888 88888P"   "Y88888 888  "Y88888  "Y8888       8888888P"  888  88888P'  "Y888 888     888 88888P"   "Y88888  "Y888 "Y88P"  888
                  888           888
                  888      Y8b d88P
                  888       "Y88P"

`)
	}

	log.Printf("Starting server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))

	myConsumer.Execute(conf.Environment, mainLogChan, errorLogChan)
}
