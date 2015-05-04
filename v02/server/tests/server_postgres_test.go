// @author Florin Patan <florinpatan@gmail.com>

// +build postgres

package server_test

import (
	"flag"
	"runtime"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/logger"

	. "gopkg.in/check.v1"
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	flag.Parse()

	if *doCurlLogs {
		*doLogTest = true
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")

	if *doLogResponseTimes {
		go logger.TGLogResponseTimes(mainLogChan)
		go logger.TGLogResponseTimes(errorLogChan)
	} else if *doLogTest {
		if *doCurlLogs {
			go logger.TGCurlLog(mainLogChan)
		} else {
			go logger.TGLog(mainLogChan)
		}
		go logger.TGLog(errorLogChan)
	} else {
		go logger.TGSilentLog(mainLogChan)
		go logger.TGSilentLog(errorLogChan)
	}
}
