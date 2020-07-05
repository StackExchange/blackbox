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
  logErr = bblog.GetErr()
  logDebug = bblog.GetDebug(debug)
}

Or in a function:

  logErr := bblog.GetErr()
  logDebug := bblog.GetDebug(debug)
	logDebug.Printf("whatever: %v", err)

*/

var logErr *log.Logger
var logDebug *log.Logger

// GetErr returns a logger handle used for errors
func GetErr() *log.Logger {
	if logErr == nil {
		logErr = log.New(os.Stderr, "", 0)
	}
	return logErr
}

// GetDebug returns a Logger handle used for debug info (output is discarded if viable=false)
func GetDebug(visible bool) *log.Logger {
	if visible {
		logDebug = log.New(os.Stderr, "", 0)
	} else {
		// Invisible mode (i.e. display nothing)
		logDebug = log.New(ioutil.Discard, "", 0)
	}
	return logDebug
}
