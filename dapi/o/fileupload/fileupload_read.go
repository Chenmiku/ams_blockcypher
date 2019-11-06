package fileupload

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*FileUpload, error) {
	var f FileUpload
	return &f, TableFileUpload.ReadByID(id, &f)
}

func GetByName(name string) (*FileUpload, error) {
	var f FileUpload
	var where = mgo.M{
		"name": name,
	}
	return &f, TableFileUpload.ReadOne(where, &f)
}

func GetFileUploadByChapterID(chapterID string, pageSize int, pageNumber int, sortBy string, sortOrder string, fileupload *[]FileUpload) (int, error) {
	var where bson.M
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	if chapterID != "" {
		where = bson.M{
			"dtime":      0,
			"chapter_id": chapterID,
		}
	} else {
		where = bson.M{
			"dtime": 0,
		}
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return TableFileUpload.ReadPagingByQuery(query, fileupload)
}
