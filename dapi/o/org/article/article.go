package article

import (
	"myproject/dapi/o/model"
)

type Article struct {
	model.WithCode  `bson:",inline"`
	Name            *string   `bson:"name,omitempty" json:"name"`
	CategoryID      *string   `bson:"category_id,omitempty" json:"category_id"`
	Parent          *string   `bson:"parent,omitempty" json:"parent"`
	Content         *string   `bson:"content,omitempty" json:"content"`
	Title           *string   `bson:"title,omitempty" json:"title"`
	Summary         *string   `bson:"summary,omitempty" json:"summary"`
	Author          *string   `bson:"author,omitempty" json:"author"`
	Signature       *string   `bson:"signature,omitempty" json:"signature"`
	Link            *string   `bson:"link,omitempty" json:"link"`
	Source          *string   `bson:"source,omitempty" json:"source"`                   // Nguồn
	Description     *string   `bson:"description,omitempty" json:"description"`         // Description
	Alias           *string   `bson:"alias,omitempty" json:"alias"`                     // Alias
	MetaKeyword     *string   `bson:"metaKeyword,omitempty" json:"metaKeyword"`         // MetaKeyword
	MetaDescription *string   `bson:"metaDescription,omitempty" json:"metaDescription"` // MetaDescription
	Active          *bool     `bson:"active,omitempty" json:"active"`
	Viewer          *int64    `bson:"viewer,omitempty" json:"viewer"`           // Viewer
	Liker           *int64    `bson:"liker,omitempty" json:"liker"`             // Liker
	ShareFace       *int64    `bson:"shareFace,omitempty" json:"shareFace"`     // ShareFace
	ShareGoogle     *int64    `bson:"shareGoogle,omitempty" json:"shareGoogle"` // ShareGoogle
	Tags            *[]string `bson:"tags,omitempty" json:"tags"`
	CTime           int64     `bson:"ctime,omitempty" json:"ctime"` // Create Time
}

func (a *Article) EnsureUniqueArticleName(newValue *string) error {
	if err := TableArticle.NotExist(map[string]interface{}{
		"name": newValue,
	}); err != nil {
		return err
	}

	return nil
}
