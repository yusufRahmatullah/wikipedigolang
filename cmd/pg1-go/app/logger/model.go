package logger

import (
	"log"
	"os"
	"time"

	"git.heroku.com/pg1-go-work/cmd/pg1-go/app/base"

	"github.com/globalsign/mgo"

	"github.com/globalsign/mgo/bson"
)

// dbLogs holds information of logs that written to database
type dbLogs struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"created_at,omitempty" bson:"created_at"`
	Name      string        `json:"name" bson:"name"`
	Level     string        `json:"level" bson:"level"`
	Message   string        `json:"message" bson:"message"`
	Cause     string        `json:"cause" bson:"cause"`
}

// Logger has objective to write logs to command line and database
type Logger struct {
	Name       string
	IsToDB     bool
	IsToStdOut bool
}

const logsCol = "logs"

var ginMode string

func init() {
	ginMode = os.Getenv("GIN_MODE")
	if ginMode != "release" {
		ginMode = "debug"
		log.Println("[Logger.init] Using debug mode, all logs will be shown")
	}
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(logsCol)
	index := mgo.Index{
		Key:        []string{"name", "level"},
		Background: true,
	}
	err := col.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// NewLogger instantiate new Logger Logger
func NewLogger(name string, isToDB, isToStdout bool) *Logger {
	return &Logger{
		Name:       name,
		IsToDB:     isToDB,
		IsToStdOut: isToStdout,
	}
}

func logToStdOut(level, name, message string) {
	log.Printf("[%v] %v: %v\n", level, name, message)
}

func errToStdOut(level, name, message string, err error) {
	log.Printf("[%v] %v: %v cause: %v\n", level, name, message, err)
}

func logToDB(level, name, message string) {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(logsCol)
	err := col.Insert(&dbLogs{
		CreatedAt: time.Now(),
		Name:      name,
		Level:     level,
		Message:   message,
	})
	if err != nil {
		logToStdOut("WARN", name, "Cannot write to database")
	}
}

func errToDB(level, name, message string, err error) {
	dataAccess := base.NewDataAccess()
	defer dataAccess.Close()
	col := dataAccess.GetCollection(logsCol)
	cause := err.Error()
	if err == nil {
		cause = ""
	}
	ierr := col.Insert(&dbLogs{
		CreatedAt: time.Now(),
		Name:      name,
		Level:     level,
		Message:   message,
		Cause:     cause,
	})
	if ierr != nil {
		logToStdOut("WARN", name, "Cannot write to database")
	}
}

func isDebug() bool {
	return ginMode == "debug"
}

// Debug force print logs to Stdout if GIN_MODE is not release
func (lg *Logger) Debug(msg string) {
	if isDebug() {
		logToStdOut("DEBUG", lg.Name, msg)
	}
}

// Info print logs to StdOut if IsToStdOur true or
// write to database if IsToDB is true
func (lg *Logger) Info(msg string) {
	if lg.IsToStdOut {
		logToStdOut("INFO", lg.Name, msg)
	}
	// info level will not be used on release version
	if lg.IsToDB && !isDebug() {
		logToDB("INFO", lg.Name, msg)
	}
}

// Warning force print logs to StdOut and
// write to database if IsToDB is true
func (lg *Logger) Warning(msg string) {
	logToStdOut("WARN", lg.Name, msg)
	if lg.IsToDB {
		logToDB("WARN", lg.Name, msg)
	}
}

// Fatal force print logs to StdOut and database
func (lg *Logger) Fatal(msg string, err error) {
	errToStdOut("FATAL", lg.Name, msg, err)
	errToDB("FATAL", lg.Name, msg, err)
}

// Error force print logs to StdOut and database
// and stop the programs
func (lg *Logger) Error(msg string, err error) {
	errToStdOut("ERROR", lg.Name, msg, err)
	errToDB("ERROR", lg.Name, msg, err)
	panic(msg)
}
