package shared

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

type Query struct {
	BranchID    []string
	DeviceID    []string
	Map         []string
	PageSize    int
	PageNumber  int
	SortBy      string
	SortOrder   string
	SearchBy    string
	SearchValue string
}

func NewQuery(r *http.Request) *Query {

	type bodyData struct {
		BranchID []string `json:"branch_id"`
		DeviceID []string `json:"device_id"`
		Map      []string `json:"map"`
	}

	var db = &bodyData{}

	json.NewDecoder(r.Body).Decode(db)

	values := r.URL.Query()

	var pageSize = values.Get("page_size")
	var pageNumber = values.Get("page_number")
	var sortBy = values.Get("sort_by")
	var sortOrder = values.Get("sort_order")

	var searchBy = values.Get("search_by")
	var searchValue = values.Get("search_value")

	return &Query{
		BranchID:    db.BranchID,
		DeviceID:    db.DeviceID,
		PageSize:    StrToInt(pageSize),
		PageNumber:  StrToInt(pageNumber),
		SortBy:      sortBy,
		SortOrder:   sortOrder,
		Map:         db.Map,
		SearchBy:    searchBy,
		SearchValue: searchValue,
	}
}

func (q *Query) CheckMapTable(table string) bool {
	for _, m := range q.Map {
		if m == table {
			return true
		}
	}

	return false
}

func (q *Query) Sum(value interface{}) bson.M {
	return bson.M{"$sum": value}
}

func (q *Query) First(value interface{}) bson.M {
	return bson.M{"$first": value}
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
