package upgrade

import (
	"db/mgo"
	"myproject/dapi/o/model"
)

var TableUpgrade = model.NewTable("upgrade")

func (b *Upgrade) Create() error {
	return TableUpgrade.Create(b)
}

func MarkDelete(id string) error {
	return TableUpgrade.MarkDelete(id)
}

func (v *Upgrade) Update(newValue *Upgrade) error {
	return TableUpgrade.UpdateByID(v.ID, newValue)
}

func CheckNameDuplicate(name string, branchID string) bool {
	var d []*Upgrade
	err := TableUpgrade.ReadMany(mgo.M{
		"name":      name,
		"branch_id": branchID,
	}, &d)

	if err == nil && len(d) > 0 {
		return true
	}

	return false
}
