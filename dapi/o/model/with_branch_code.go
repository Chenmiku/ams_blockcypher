package model

type WithBranchCode struct {
	WithCode `bson:",inline"`
	BranchID string `bson:"branch_id" json:"branch_id"`
}

type IWithBranchCode interface {
	GetBranchID() string
	IWithCode
}

func (v *WithBranchCode) GetBranchID() string {
	return v.BranchID
}

func (t *TableWithBranchCode) Create(v IWithBranchCode) error {
	// check code
	err := t.NotExist(map[string]interface{}{"branch_id": v.GetBranchID(), "code": v.GetCode(), "dtime": 0})
	// var tableName = t.Table.Name
	// if tableName == "screens" || tableName == "kiosks" || tableName == "counters" {
	// 	var count, _ = t.UnsafeCount(map[string]interface{}{})
	// 	if limit := license.GetLimitByName(tableName); count > limit {
	// 		glog.Info("Limit license ", limit, " count ", count)
	// 		return web.BadRequest("Limited license")
	// 	}
	// }
	if err != nil {
		return err
	}
	return t.Table.Create(v)
}

func (t *TableWithBranchCode) UpdateByID(oldID string, oldBranchID string, oldCode string, v IWithBranchCode) error {
	if v.GetBranchID() != oldBranchID || v.GetCode() != oldCode {
		// check code
		err := t.NotExist(map[string]interface{}{"branch_id": v.GetBranchID(), "code": v.GetCode(), "dtime": 0})
		if err != nil {
			return err
		}
	}
	return t.Table.UpdateByID(oldID, v)
}
