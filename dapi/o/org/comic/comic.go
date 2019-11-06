package comic

import (
	"db/mgo"
)

// Table Comic
type Comic struct {
	mgo.BaseModel   `bson:",inline"`
	Name            *string   `bson:"name,omitempty" json:"name"`
	OrtherName      *string   `bson:"orther_name,omitempty" json:"orther_name"`
	Thumbnail       *string   `bson:"thumbnail,omitempty" json:"thumbnail"`
	Summary         *string   `bson:"summary,omitempty" json:"summary"`
	Categories      *[]string `bson:"categories,omitempty" json:"categories"`
	Author          *string   `bson:"author,omitempty" json:"author"`
	Source          *string   `bson:"source,omitempty" json:"source"`
	Alias           *string   `bson:"alias,omitempty" json:"alias"`
	MetaKeyword     *string   `bson:"metakeyword,omitempty" json:"metakeyword"`
	MetaDescription *string   `bson:"metadescription,omitempty" json:"metadescription"`
	Link            *string   `bson:"link,omitempty" json:"link"`
	State           *bool     `bson:"state,omitempty" json:"state"`
	Vote            *int64    `bson:"vote,omitempty" json:"vote"`
	Viewer          *int64    `bson:"viewer,omitempty" json:"viewer"`
	Sharer          *int64    `bson:"sharer,omitempty" json:"sharer"`
	Liker           *int64    `bson:"liker,omitempty" json:"liker"`
	Hot             *bool     `bson:"hot,omitempty" json:"hot"` // trend, most viewer
	New             *bool     `bson:"new,omitempty" json:"new"` // new comic
	Top             *bool     `bson:"top,omitempty" json:"top"` // must view
	Approve         *bool     `bson:"approve,omitempty" json:"approve"`
	Active          *bool     `bson:"active,omitempty" json:"active"`
	CTime           int64     `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func (a *Comic) EnsureUniqueComicName(newValue *string) error {
	if err := TableComic.NotExist(map[string]interface{}{
		"name": newValue,
	}); err != nil {
		return err
	}

	return nil
}
