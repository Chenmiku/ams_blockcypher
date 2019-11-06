package activity

type Activity struct {
	AllDevice                 int64   `json:"all_device" bson:"all_device"`
	OnTimeAverage             int64   `json:"on_time_average" bson:"on_time_average"`
	AllTimeOn                 int64   `json:"all_on_time" bson:"all_on_time"`
	BestPerfomanceDevice      *Device `json:"best_perfomance_device" bson:"best_perfomance_device"`
	WorstPerformanceEquipment *Device `json:"worst_perfomance_device" bson:"worst_perfomance_device"`
}

type Device struct {
	ID            string `json:"device_id" bson:"_id"`
	Name          string `json:"device_name" bson:"device_name"`
	BranchName    string `json:"branch_name" bson:"branch_name"`
	BranchID      string `json:"branch_id" bson:"branch_id"`
	AllOnTime     int64  `json:"all_on_time" bson:"all_on_time"`
	OnTimeAverage int64  `json:"on_time_average" bson:"on_time_average"`
}

type Chart struct {
	ID     string `json:"_id" bson:"_id"`
	OnTime int64  `json:"on_time" bson:"on_time"`
}

type Result struct {
	Activity Aggregate `json:"activity" bson:"activity"`
	Devices  []*Chart  `json:"devices" bson:"devices"`
}

type Aggregate struct {
	Activity Activity  `json:"aggregate" bson:"aggregate"`
	Days     []*Device `json:"days" bson:"days"`
}
