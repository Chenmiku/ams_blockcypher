package chapter

import (
	"myproject/dapi/o/model"
	"time"
)

//TableChapter: Table in DB
var TableChapter = model.NewTable("chapters")

//Create :
func (c *Chapter) Create() error {
	c.CTime = time.Now().Unix()
	return TableChapter.Create(c)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableChapter.MarkDelete(id)
}

//Update :
func (c *Chapter) Update(newValue *Chapter) error {
	return TableChapter.UpdateByID(c.ID, newValue)
}
