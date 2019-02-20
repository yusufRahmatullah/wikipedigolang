package service

import (
	"log"

	"git.heroku.com/pg1-go-work/cmd/igoweb/config"
	"git.heroku.com/pg1-go-work/cmd/igoweb/model"
	"git.heroku.com/pg1-go-work/cmd/igoweb/repository"
)

// LoggerService has objective to write Log to command line and database
type LoggerService struct {
	Name       string
	Repository repository.LogRepository
}

// Debug force print Log to STDOUT if config.GinMode is debug
func (logger *LoggerService) Debug(message string) {
	if config.GetInstance().GinMode == "debug" {
		l := logger.instantiateLog(model.LevelDebug, message)
		l.InitTimeStamp()
		logger.logToStdOut(l)
	}
}

// Info always print log to STDOUT
func (logger *LoggerService) Info(message string) {
	l := logger.instantiateLog(model.LevelInfo, message)
	l.InitTimeStamp()
	logger.logToStdOut(l)
}

// Warning print log to STDOUT and save to DB
func (logger *LoggerService) Warning(message, cause string) {
	l := logger.instantiateLog(model.LevelWarning, message)
	l.InitTimeStamp()
	l.Cause = cause
	logger.logToStdOut(l)
	logger.Repository.Create(l)
}

// Fatal print log to STDOUT and save to DB
func (logger *LoggerService) Fatal(message string, err error) {
	l := logger.instantiateLog(model.LevelFatal, message)
	l.InitTimeStamp()
	l.Cause = err.Error()
	logger.logToStdOut(l)
	logger.Repository.Create(l)
}

// Error print log to STDOUT, save to DB, and force close
func (logger *LoggerService) Error(message string, err error) {
	l := logger.instantiateLog(model.LevelError, message)
	l.InitTimeStamp()
	l.Cause = err.Error()
	logger.Repository.Create(l)
	log.Panicf("%v [%v]-%v: %v cause: %v\n", l.CreatedAt, l.Name, l.Level, l.Message, l.Cause)
}

func (logger *LoggerService) instantiateLog(level model.LogLevel, message string) *model.Log {
	return &model.Log{
		Cause:   "-",
		Level:   level,
		Message: message,
		Name:    logger.Name,
	}
}

func (logger *LoggerService) logToStdOut(l *model.Log) {
	log.Printf("%v [%v]-%v: %v cause: %v\n", l.CreatedAt, l.Name, l.Level, l.Message, l.Cause)
}
