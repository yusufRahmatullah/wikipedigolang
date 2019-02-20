package model

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultOffset = 0
	defaultLimit  = 24
)

// FindRequest holds query params for handler with find capability
type FindRequest struct {
	Offset int
	Limit  int
	Query  string
	Sort   string
}

// GetFindRequest extract FindRequest from given context
func GetFindRequest(c *gin.Context) *FindRequest {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")
	query := c.Query("query")
	offset := convertIntOrDefault(offsetStr, defaultOffset)
	limit := convertIntOrDefault(limitStr, defaultLimit)

	sort := generateSortOrder(c)
	return &FindRequest{
		Offset: offset,
		Limit:  limit,
		Sort:   sort,
		Query:  query,
	}
}

func convertIntOrDefault(text string, def int) int {
	if text == "" {
		return def
	}
	num, err := strconv.Atoi(text)
	if err != nil {
		return def
	}
	return num
}

func generateSortOrder(c *gin.Context) string {
	sort := c.Query("sort")
	orderStr := c.Query("order")
	order := convertIntOrDefault(orderStr, -1)
	if sort == "" {
		sort = "_id"
	}
	if order == -1 {
		sort = "-" + sort
	}
	return sort
}
