package logging

import (
	"fmt"
	"log"
	"runtime"
)

const (
	DEBUG = "[DEBUG] "
	INFO  = "[INFO] "
	WARN  = "[WARN] "
	ERR   = "[ERROR] "
)

const (
	COLOR_RED = uint8(iota + 91)
	COLOR_GREEN
	COLOR_YELLOW
	COLOR_BLUE
)

func Debug(v string) {
	log.Printf(blue(DEBUG)+"%s", v)
}

func Debugf(format string, v ...interface{}) {
	log.Printf(blue(DEBUG)+format, v...)
}

func Info(v string) {
	log.Printf(green(INFO)+"%s", v)
}

func Infof(format string, v ...interface{}) {
	log.Printf(green(INFO)+format, v...)
}

func Warn(v string) {
	log.Printf(yellow(WARN)+"%s", v)
}

func Warnf(format string, v ...interface{}) {
	log.Printf(yellow(WARN)+format, v...)
}

func Error(v string) {
	log.Printf(red(ERR)+"%s \n %v", v, getStacksMsg(stacks(4)))
}

func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(red(ERR)+format, v...)
	log.Printf(msg+"\n %v", getStacksMsg(stacks(4)))
}

func Fatal(v string) {
	log.Fatalf(red(ERR)+"%s", v)
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf(red(ERR)+format, v...)
}

// -----------------------------Stack Info-----------------------------------

func getStacksMsg(stacks []*stack) string {
	var stacksMsg string
	for _, stack := range stacks {
		stacksMsg += fmt.Sprintf("\tat %v \n", stack.String())
	}
	return stacksMsg
}

func stacks(skip int) []*stack {
	var stacks []*stack
	pc := make([]uintptr, 32)
	n := runtime.Callers(skip, pc)
	for i := 0; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		stacks = append(stacks, &stack{Source: file, Line: line, Func: f.Name()})
	}
	return stacks
}

type stack struct {
	Func   string `json:"func"`
	Line   int    `json:"line"`
	Source string `json:"source"`
}

func (s *stack) String() string {
	return fmt.Sprintf("%s:%d (Method: %s)", s.Source, s.Line, s.Func)
}

// -------------------------------COLOR---------------------------------

func red(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", COLOR_RED, s)
}

func green(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", COLOR_GREEN, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", COLOR_YELLOW, s)
}

func blue(s string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", COLOR_BLUE, s)
}
