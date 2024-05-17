package loger

import "log"

var LogPrint LoggerAll

type Logger interface {
	Log(msg string)
	Package(name string) Logger
	Module(name string)
}

type LoggerAll struct {
	NAME_MODULE  string
	NAME_PACKAGE string
}

func (l *LoggerAll) Log(msg string) {
	log.Printf("[%s] [%s] << %s >>\n", l.NAME_MODULE, l.NAME_PACKAGE, msg)
}

func (l *LoggerAll) Package(name string) Logger {
	l.NAME_PACKAGE = name
	return l
}

func (l *LoggerAll) Module(name string) {
	l.NAME_MODULE = name
}
