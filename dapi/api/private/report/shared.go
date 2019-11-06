package report

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"util/timer"
)

type Query struct {
	values     url.Values
	Start      time.Time
	End        time.Time
	BranchID   []string
	DeviceID   []string
	PageSize   int
	PageNumber int
	SortBy     string
	SortOrder  string
}

func NewQuery(r *http.Request) *Query {

	type bodyData struct {
		BranchID []string `json:"branch_id"`
		DeviceID []string `json:"device_id"`
	}

	var db = &bodyData{}

	json.NewDecoder(r.Body).Decode(db)

	values := r.URL.Query()

	var startTime = StrToTime(values.Get("start"))
	var endTime = StrToTime(values.Get("end"))
	var pageSize = values.Get("page_size")
	var pageNumber = values.Get("page_number")
	var sortBy = values.Get("sort_by")
	var sortOrder = values.Get("sort_order")

	endTime = AddTime(endTime, 86399)

	return &Query{
		Start:      startTime,
		End:        endTime,
		BranchID:   db.BranchID,
		DeviceID:   db.DeviceID,
		PageSize:   StrToInt(pageSize),
		PageNumber: StrToInt(pageNumber),
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}
}

func (q *Query) Sum(value interface{}) bson.M {
	return bson.M{"$sum": value}
}

func (q *Query) First(value interface{}) bson.M {
	return bson.M{"$first": value}
}

func (q *Query) GetTimeFilter() int64 {
	return q.End.Unix() - q.Start.Unix()
}

func AddTime(t time.Time, sec int64) time.Time {
	return time.Unix(t.Unix()+sec, 0)
}

func StrToTime(s string) time.Time {
	var t, _ = time.Parse("2006-01-02", s)
	z := t.Unix() - timer.TimeZone()
	time := time.Unix(z, 0)
	return time
}

func StrToTimeFormat(s string, format string) time.Time {
	var t, _ = time.Parse(format, s)
	z := t.Unix() - timer.TimeZone()
	time := time.Unix(z, 0)
	return time
}

func TimeToDay(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetTimeOfDay(timestamp int64) string {
	var t = time.Unix(timestamp, 0)
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func GetTimeDuration(d int64) string {
	s := (d % 3600) % 60
	m := (d % 3600) / 60
	h := d / 3600
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func StrToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}
