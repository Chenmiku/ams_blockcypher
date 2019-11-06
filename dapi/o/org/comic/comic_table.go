package comic

import (
	"errors"
	"myproject/dapi/o/model"
	"time"
)

//Tablecomic : Table in DB
var TableComic = model.NewTable("comics")

//Create :
func (c *Comic) Create() error {
	if err := c.EnsureUniqueComicName(c.Name); err != nil {
		return errors.New("comic_already_exists")
	}

	c.CTime = time.Now().Unix()
	return TableComic.Create(c)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableComic.MarkDelete(id)
}

//Update :
func (c *Comic) Update(newValue *Comic) error {
	if newValue.Name != c.Name {
		if err := c.EnsureUniqueComicName(newValue.Name); err != nil {
			return errors.New("comic_already_exists")
		}
	}
	return TableComic.UpdateByID(c.ID, newValue)
}
