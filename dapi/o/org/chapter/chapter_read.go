package chapter

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Chapter, error) {
	var b Chapter
	return &b, TableChapter.ReadByID(id, &b)
}

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, chapter *[]Chapter) (int, error) {
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	where := bson.M{
		"dtime": 0,
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return TableChapter.ReadPagingByQuery(query, chapter)
}

func GetByComicID(pageSize int, pageNumber int, sortBy string, sortOrder string, chapter *[]Chapter, id string) (int, error) {
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	where := bson.M{
		"dtime":    0,
		"comic_id": id,
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return TableChapter.ReadPagingByQuery(query, chapter)
}

func Query(q bson.M) ([]*Chapter, error) {
	var brs []*Chapter
	err := TableChapter.UnsafeReadMany(q, &brs)
	return brs, err
}

func GetLatestChapterByComicID(id string, chapter *Chapter) error {
	var query []bson.M
	var filterChapter bson.M

	filterChapter = bson.M{
		"dtime":    0,
		"comic_id": id,
	}

	query = append(query, bson.M{
		"$match": filterChapter,
	}, bson.M{
		"$sort": bson.M{"display_order": -1},
	}, bson.M{
		"$limit": 1,
	})

	return TableChapter.ReadOneByQuery(query, chapter)
}
