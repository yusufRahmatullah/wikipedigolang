package model

// FindRequest holds query params for handler with find capability
type FindRequest struct {
	Offset int
	Limit  int
	Query  string
	Sort   string
}
