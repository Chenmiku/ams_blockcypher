package menu

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Menu, error) {
	var b Menu
	return &b, TableMenu.ReadByID(id, &b)
}

func GetParent(id string) (*Menu, error) {
	var b Menu
	return &b, TableMenu.C().Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"parent": id,
			},
		},
	}).One(&b)
}

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, menu *[]Menu) (int, error) {
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	where := bson.M{
		"dtime":  0,
		"active": true,
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return TableMenu.ReadPagingByQuery(query, menu)
}

func Query(q bson.M) ([]*Menu, error) {
	var brs []*Menu
	err := TableMenu.UnsafeReadMany(q, &brs)
	return brs, err
}
