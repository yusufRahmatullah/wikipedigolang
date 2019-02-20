package model

// LogLevel is the level of Log
type LogLevel string

const (
	// LevelDebug means the log level is for debugging
	LevelDebug LogLevel = "DEBUG"
	// LevelInfo means the log level is for show info
	LevelInfo LogLevel = "INFO"
	// LevelWarning means the log level is for show warning
	LevelWarning LogLevel = "WARNING"
	// LevelFatal means the log level for fatal condition
	LevelFatal LogLevel = "FATAL"
	// LevelError means the log level for error condition and will stop the app
	LevelError LogLevel = "ERROR"
)

// Log holds information about logging data
type Log struct {
	MongoID
	TimeStamp
	Name    string   `json:"name" bson:"name"`
	Level   LogLevel `json:"level" bson:"level"`
	Message string   `json:"message" bson:"message"`
	Cause   string   `json:"cause" bson:"cause"`
}
