package activity

import (
	"db/mgo"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func GetAllFromDB(stime string, etime string, branchID, deviceID, deviceType string) *[]Activity {
	var res = make([]Activity, 0)

	var query = mgo.M{}
	query["bid"] = branchID
	query["eid"] = deviceID
	query["cat"] = deviceType

	query["date"] = mgo.M{"$gte": stime, "$lte": etime}
	err := ActivityTable.C().Find(query).Sort("s_at").All(&res)
	if err != nil {
		fmt.Println("ERROR: Activity_table.go	-> GetAllFromDB: ", err.Error())
	}
	return &res
}

func ReadActivitiesByFilter(filter bson.M, limit int, skip int) []*Activity {
	var res = make([]*Activity, 0)
	ActivityTable.C().Find(filter).Sort("s_at").Skip(skip).Limit(limit).All(&res)
	return res
}

func CountActivitiesByFilter(filter bson.M) int {
	var count, _ = ActivityTable.C().Find(filter).Count()
	return count
}

func ReadActivitiesByPipeline(pipeline []bson.M) *[]Activity {
	var res = make([]Activity, 0)
	ActivityTable.C().Pipe(pipeline).All(&res)
	return &res
}

func KioskTime(pipeline []bson.M) *[]Activity {
	var res = make([]Activity, 0)
	err := ActivityTable.C().Pipe(pipeline).All(&res)
	if err != nil {
		fmt.Println(err)
	}
	return &res
}
