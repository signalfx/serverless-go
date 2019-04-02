package sfxgcf

import (
    "os"
    "log"
)

type Logger struct {
	log *log.Logger
	error *log.Logger
}

// When log : logger.log("Logging.")
// When errored : logger.error("Error!")
var logger = &Logger{
	log: log.New(os.Stdout, "sfx_wrapper", 0),
	error: log.New(os.Stderr, "sfx_wrapper", 0),
}