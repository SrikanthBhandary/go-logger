package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int8

var GlobalLevel Level

type Hook func() string

func SetGlobalLevel(level Level) {
	GlobalLevel = level
}

const (
	DebugLevel Level = iota
	WarningLevel
	InfoLevel
	ErrorLevel
	Disabled = -1
)

type Logger struct {
	w    io.Writer
	hook []Hook
}

type Event struct {
	level Level
	w     io.Writer
	hook  []Hook
}

func NewEvent(level Level, w io.Writer, hook []Hook) *Event {
	e := &Event{level: level, w: w}
	e.hook = append(hook, e.levelHook)
	return e
}

func (e *Event) levelHook() string {
	return fmt.Sprintf("level=%s", GetPrefix(e.level))

}

func (e *Event) Msg(msg []byte) {
	if (GlobalLevel != DebugLevel) && (e.level != GlobalLevel) {
		return
	}
	for _, fn := range e.hook {
		e.w.Write([]byte(fn()))
		e.w.Write([]byte(","))
	}
	e.w.Write([]byte("msg: "))
	e.w.Write(msg)
	e.w.Write([]byte("\n"))
}

func GetPrefix(level Level) string {
	switch level {
	case DebugLevel:
		return "Debug"
	case WarningLevel:
		return "Warn"
	case InfoLevel:
		return "Info"
	case ErrorLevel:
		return "Error"
	default:
		return ""
	}
}

func timeHook() string {
	return fmt.Sprintf("time=%s", time.Now().Format("2006-01-02 15:04:05"))
}

func New(w io.Writer) *Logger {
	return &Logger{w: w, hook: []Hook{timeHook}}
}

func (l *Logger) Info() *Event {
	return NewEvent(InfoLevel, l.w, l.hook)
}

func (l *Logger) Warn() *Event {
	return NewEvent(WarningLevel, l.w, l.hook)
}

func (l *Logger) Error() *Event {
	return NewEvent(ErrorLevel, l.w, l.hook)
}

func (l *Logger) Write(msg []byte) {
	for _, fn := range l.hook {
		l.w.Write([]byte(fn()))
		l.w.Write([]byte(","))
	}
	l.w.Write([]byte("msg: "))
	l.w.Write(msg)
	l.w.Write([]byte("\n"))
}

func main() {
	logger := New(os.Stdout)
	SetGlobalLevel(InfoLevel)
	logger.Info().Msg([]byte("Testing"))
	logger.Info().Msg([]byte("Testing1"))
	logger.Warn().Msg([]byte("Testing"))
	logger.Warn().Msg([]byte("Testing1"))
	logger.Error().Msg([]byte("Testing"))
	logger.Error().Msg([]byte("Testing1"))
}
