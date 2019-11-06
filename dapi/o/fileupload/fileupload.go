package fileupload

import (
	"db/mgo"
)

type FileUpload struct {
	mgo.BaseModel         `bson:",inline"`
	Name                  string `bson:"name" json:"name"`
	ChapterID             string `bson:"chapter_id" json:"chapter_id"`
	ComicID             string `bson:"comic_id" json:"comic_id"`
	Type                  string `bson:"type" json:"type"`
	Size                  int64  `bson:"size" json:"size"`
	Hash                  string `bson:"hash" json:"hash"`
	Details               string `bson:"details" json:"details"`
	PhysicalPath          string `bson:"physical_path" json:"physical_path"`
	RelativePath          string `bson:"relative_path" json:"relative_path"`
	PhysicalPathThumbnail string `bson:"physical_path_thumbnail" json:"physical_path_thumbnail"`
	RelativePathThumbnail string `bson:"relative_path_thumbnail" json:"relative_path_thumbnail"`
}
