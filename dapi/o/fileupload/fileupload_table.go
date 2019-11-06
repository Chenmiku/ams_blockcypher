package fileupload

import (
	"db/mgo"
	"myproject/dapi/o/model"
)

var TableFileUpload = model.NewTable("fileupload")

func (f *FileUpload) Create() error {
	return TableFileUpload.Create(f)
}

func (f *FileUpload) Update(newValue *FileUpload) error {
	return TableFileUpload.UpdateByID(f.ID, newValue)
}

func MarkDelete(id string) error {
	return TableFileUpload.MarkDelete(id)
}

func CheckNameDuplicate(name string, chapterID string) bool {
	var f []*FileUpload
	err := TableFileUpload.ReadMany(mgo.M{
		"name":       name,
		"chapter_id": chapterID,
	}, &f)

	if err == nil && len(f) > 0 {
		return true
	}

	return false
}
