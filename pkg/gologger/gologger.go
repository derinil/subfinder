package gologger

import "fmt"

type Event struct {
	level Level
}

type Level int

const (
	LevelFatal Level = iota
	LevelSilent
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelVerbose
)

var MaxLevel Level = LevelInfo

func (l Level) String() string {
	return [...]string{"fatal", "silent", "error", "info", "warning", "debug", "verbose"}[l]
}

func (e *Event) Label(label string) *Event { return e }
func (e *Event) Msg(message string) {
	if e.level > MaxLevel {
		return
	}
	fmt.Println(e.level.String() + " - " + message)
}
func (e *Event) Msgf(format string, args ...any) {
	if e.level > MaxLevel {
		return
	}
	fmt.Println(e.level.String() + " - " + fmt.Sprintf(format, args...))
}

func SetMaxLevel(l Level) { MaxLevel = l }

func Warning() *Event { return &Event{level: LevelWarning} }
func Info() *Event    { return &Event{level: LevelInfo} }
func Error() *Event   { return &Event{level: LevelError} }
func Debug() *Event   { return &Event{level: LevelDebug} }
func Fatal() *Event   { return &Event{level: LevelFatal} }
func Silent() *Event  { return &Event{level: LevelSilent} }
func Print() *Event   { return &Event{level: LevelInfo} }
func Verbose() *Event { return &Event{level: LevelVerbose} }
