package gcfwrapper

import (
	"log"
	"os"
)

type gcfLogger struct {
	log *log.Logger
	err *log.Logger
}

// When log : logger.Log.("Logging.")
// When errored : logger.Error("Error!")
var logger = &gcfLogger{
	log: log.New(os.Stdout, "sfx_wrapper", 0),
	err: log.New(os.Stderr, "sfx_wrapper", 0),
}

func (l *gcfLogger) Log(str string) {
	l.log.Println(str)
}

func (l *gcfLogger) Error(str string) {
	l.err.Println(str)
}
