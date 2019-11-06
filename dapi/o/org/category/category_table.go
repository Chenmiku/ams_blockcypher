package category

import (
	"errors"
	"myproject/dapi/o/model"
	"time"
)

//TableCategory : Table in DB
var TableCategory = model.NewTable("categories")

//Create :
func (c *Category) Create() error {
	if err := c.EnsureUniqueCategoryName(c.Name); err != nil {
		return errors.New("category_already_exists")
	}

	c.CTime = time.Now().Unix()
	return TableCategory.Create(c)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableCategory.MarkDelete(id)
}

//Update :
func (c *Category) Update(newValue *Category) error {
	if newValue.Name != c.Name {
		if err := c.EnsureUniqueCategoryName(newValue.Name); err != nil {
			return errors.New("category_already_exists")
		}
	}
	return TableCategory.UpdateByID(c.ID, newValue)
}
