package block

import (
	"ams_system/dapi/o/model"
	"time"
)

var TableBlock = model.NewTable("blocks")

func (w *Block) Create() error {
	w.CTime = time.Now().Unix()

	return TableBlock.Create(w)
}

func MarkDelete(id string) error {
	return TableBlock.MarkDelete(id)
}

func (w *Block) UpdateById(newvalue *Block) error {
	return TableBlock.UnsafeUpdateByID(w.ID, newvalue)
}
