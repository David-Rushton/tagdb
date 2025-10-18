package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	panicLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO  ", log.LstdFlags)
	warnLogger = log.New(os.Stdout, "WARN  ", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "ERR   ", log.LstdFlags)
	fatalLogger = log.New(os.Stderr, "FATAL ", log.LstdFlags)
	panicLogger = log.New(os.Stderr, "PANIC ", log.LstdFlags)
}

func Info(v ...any) {
	infoLogger.Print(v...)
}

func Warn(v ...any) {
	warnLogger.Print(v...)
}

func Error(v ...any) error {
	errorLogger.Print(v...)
	return errors.New(fmt.Sprint(v...))
}

func Fatal(v ...any) {
	fatalLogger.Fatal(v...)
}

func Panic(v ...any) {
	panicLogger.Panic(v...)
}

func Infof(msg string, a ...any) {
	infoLogger.Printf(msg, a...)
}

func Warnf(msg string, a ...any) {
	warnLogger.Printf(msg, a...)
}

func Errorf(msg string, a ...any) error {
	err := fmt.Errorf(msg, a...)
	errorLogger.Print(err)
	return err
}

func Fatalf(msg string, a ...any) {
	fatalLogger.Fatalf(msg, a...)
}

func Panicf(msg string, a ...any) {
	panicLogger.Panicf(msg, a...)
}
