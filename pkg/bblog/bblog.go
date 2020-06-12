package bblog

import (
	"io/ioutil"
	"log"
	"os"
)

/*

To use this, include the following lines in your .go file.

var logErr *log.Logger
var logDebug *log.Logger
func init() {
  logErr = bblog.GetLogErr()
  logDebug = bblog.GetLogDebug(verbose)
}

Or in a function:

  logErr := bblog.GetLogErr()
  logDebug := bblog.GetLogDebug(verbose)
	logDebug.Printf("whatever")

*/

var logErr *log.Logger
var logDebug *log.Logger

func GetErr() *log.Logger {
	if logErr == nil {
		logErr = log.New(os.Stderr, "", 0)
	}
	return logErr
}

func GetDebug(visible bool) *log.Logger {
	if visible {
		logDebug = log.New(os.Stderr, "", 0)
	} else {
		// Invisible mode (i.e. display nothing)
		logDebug = log.New(ioutil.Discard, "", 0)
	}
	return logDebug
}
