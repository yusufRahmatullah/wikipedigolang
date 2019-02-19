package database

// Database is interface for connect to database
type Database interface {
	Prepare() error
	Close()
}
