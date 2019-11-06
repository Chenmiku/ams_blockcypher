package contact

import (
	"myproject/dapi/o/model"
	"time"
)

//TableContact : Table in DB
var TableContact = model.NewTable("contacts")

//Create :
func (c *Contact) Create() error {
	c.CTime = time.Now().Unix()
	return TableContact.Create(c)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableContact.MarkDelete(id)
}

//Update :
func (c *Contact) Update(newValue *Contact) error {
	return TableContact.UpdateByID(c.ID, newValue)
}
