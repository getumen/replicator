package log

import "log"

type stderrLogger struct {
	debugLogger  *log.Logger
	infoLogger   *log.Logger
	waningLogger *log.Logger
	fatalLogger  *log.Logger
	panicLogger  *log.Logger
}
