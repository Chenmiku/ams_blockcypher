package menu

import (
	"errors"
	"myproject/dapi/o/model"
	"time"
)

//TableMenu : Table in DB
var TableMenu = model.NewTable("menus")

//Create :
func (c *Menu) Create() error {
	if err := c.EnsureUniqueMenuName(c.Name); err != nil {
		return errors.New("menu_already_exists")
	}

	c.CTime = time.Now().Unix()
	return TableMenu.Create(c)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableMenu.MarkDelete(id)
}

//Update :
func (c *Menu) Update(newValue *Menu) error {
	if newValue.Name != c.Name {
		if err := c.EnsureUniqueMenuName(newValue.Name); err != nil {
			return errors.New("menu_already_exists")
		}
	}
	return TableMenu.UpdateByID(c.ID, newValue)
}
