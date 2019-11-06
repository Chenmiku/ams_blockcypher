package chapter

import (
	"db/mgo"
)

// Table Chapter
type Chapter struct {
	mgo.BaseModel   `bson:",inline"`
	Name            *string   `bson:"name,omitempty" json:"name"`
	ComicID         *string   `bson:"comic_id,omitempty" json:"comic_id"`
	Images          *[]string `bson:"images,omitempty" json:"images"`
	Videos          *[]string `bson:"videos,omitempty" json:"videos"`
	Content         *string   `bson:"content,omitempty" json:"content"`
	DisplayOrder    *int64    `bson:"display_order,omitempty" json:"display_order"`
	Sharer          *int64    `bson:"sharer,omitempty" json:"sharer"`
	Liker           *int64    `bson:"liker,omitempty" json:"liker"`
	Viewer          *int64    `bson:"viewer,omitempty" json:"viewer"`
	Alias           *string   `bson:"alias,omitempty" json:"alias"`
	MetaKeyword     *string   `bson:"metakeyword,omitempty" json:"metakeyword"`
	MetaDescription *string   `bson:"metadescription,omitempty" json:"metadescription"`
	Link            *string   `bson:"link,omitempty" json:"link"`
	Active          *bool     `bson:"active,omitempty" json:"active"`
	CTime           int64     `bson:"ctime,omitempty" json:"ctime"` // Create Time
}
