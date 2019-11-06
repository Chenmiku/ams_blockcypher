package menu

import (
	"db/mgo"
)

type Menu struct {
	mgo.BaseModel `bson:",inline"`
	Name          *string `bson:"name,omitempty" json:"name"`                   //
	Parent        *string `bson:"parent,omitempty" json:"parent"`               //
	DisplayOrder  *int64  `bson:"display_order,omitempty" json:"display_order"` //
	Link          *string `bson:"link,omitempty" json:"link"`
	Active        *bool   `bson:"active,omitempty" json:"active"`
	CTime         int64   `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func (a *Menu) EnsureUniqueMenuName(newValue *string) error {
	if err := TableMenu.NotExist(map[string]interface{}{
		"name": newValue,
	}); err != nil {
		return err
	}

	return nil
}
