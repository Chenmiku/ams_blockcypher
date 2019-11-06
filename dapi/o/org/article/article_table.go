package article

import (
	"errors"
	"myproject/dapi/o/model"
	"time"
)

//TableArticle : Table in DB
var TableArticle = model.NewTable("articles")

//Create :
func (a *Article) Create() error {
	if err := a.EnsureUniqueArticleName(a.Name); err != nil {
		return errors.New("article_already_exists")
	}
	a.CTime = time.Now().Unix()
	return TableArticle.Create(a)
}

//MarkDelete :
func MarkDelete(id string) error {
	return TableArticle.MarkDelete(id)
}

//Update :
func (a *Article) Update(newValue *Article) error {
	if newValue.Name != a.Name {
		if err := a.EnsureUniqueArticleName(newValue.Name); err != nil {
			return errors.New("article_already_exists")
		}
	}
	return TableArticle.UpdateByID(a.ID, newValue)
}
