package upgrade

import (
	"ams_system/dapi/o/model"
)

type Upgrade struct {
	model.WithBranchID `bson:",inline"`
	Name               string `bson:"name" json:"name"`
	Version            string `bson:"version" json:"version"`
	AndroidAPK         string `bson:"android_apk" json:"android_apk"`
	TimeUpgrade        int64  `bson:"time_upgrade" json:"time_upgrade"`
	Details            string `bson:"details" json:"details"`
}
