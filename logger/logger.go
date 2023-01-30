package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type LogLevel string

const (
	Info  LogLevel = "info"
	Error LogLevel = "error"
	Panic LogLevel = "panic"
)

type LogString struct {
	strings.Builder
}

func (l *LogString) LevelBegin(level LogLevel, function string) *LogString {
	l.WriteString(fmt.Sprintf("%s: function: %s", level, function))
	return l
}

func (l *LogString) ErrorCheck(err error) *LogString {
	if err != nil {
	  l.WriteString(fmt.Sprintf(", error: %s", err))
	}
	return l
}

func (l *LogString) PanicCheck(panic interface{}) *LogString {
	if panic != nil {
		l.WriteString(fmt.Sprintf(", panic: %s", panic))
	}
	return l
}
func (l *LogString) HttpStatusWrite(status int) *LogString {
	if status != 0 {
		l.WriteString(fmt.Sprintf(", status: %d", status))
	}
	return l
}

func (l *LogString) RequestWrite(request interface{}) *LogString {
	if request != nil {
		breq, _ := json.Marshal(request)
		l.WriteString(fmt.Sprintf(", request: %s", breq))
	}
	return l
}

func (l *LogString) ResponseWrite(response interface{}) *LogString {
	if response != nil {
		bresp, _ := json.Marshal(response)
		l.WriteString(fmt.Sprintf(", response: %s", bresp))
	}
	return l
}

func (l *LogString) String() string {
	return l.Builder.String()
}

func Log(level LogLevel, function string, httpStatus int, err error, panic interface{}, req interface{}, resp interface{}) {
	logString := new(LogString)
	s := logString.LevelBegin(level, function).PanicCheck(panic).ErrorCheck(err).
		HttpStatusWrite(httpStatus).RequestWrite(req).ResponseWrite(resp).String()
	log.Println(s)
}
