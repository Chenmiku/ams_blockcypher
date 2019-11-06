package monitor

import (
	"time"
)

type MonitorData struct {
	DeviceID      string    `json:"device_id" bson:"device_id"`
	CampainID     string    `json:"campain_id" bson:"campain_id"`
	NetworkSpeed  string    `json:"network_speed" bson:"network_speed"`
	ContentLoaded bool      `json:"content_loaded" bson:"content_loaded"`
	OnlineTime    time.Time `json:"online_time" bson:"online_time"`
}
