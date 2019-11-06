package category

import (
	"db/mgo"
)

type Category struct {
	mgo.BaseModel `bson:",inline"`
	Name          *string `bson:"name,omitempty" json:"name"`     //
	Parent        *string `bson:"parent,omitempty" json:"parent"` //
	Level         *int64  `bson:"level,omitempty" json:"level"`   //
	Active        *bool   `bson:"active,omitempty" json:"active"`
	CTime         int64   `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func (a *Category) EnsureUniqueCategoryName(newValue *string) error {
	if err := TableCategory.NotExist(map[string]interface{}{
		"name": newValue,
	}); err != nil {
		return err
	}

	return nil
}
