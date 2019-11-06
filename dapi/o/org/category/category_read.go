package category

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Category, error) {
	var b Category
	return &b, TableCategory.ReadByID(id, &b)
}

func GetParent(id string) (*Category, error) {
	var b Category
	return &b, TableCategory.C().Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"parent": id,
			},
		},
	}).One(&b)
}

// func Search(branchIDs []string) ([]*Branch, error) {
// 	var branch = []*Branch{}
// 	var err = TableBranch.ReadManyIn("branch_id", branchIDs, &branch)
// 	return branch, err
// }

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, cate *[]Category) (int, error) {
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

	return TableCategory.ReadPagingByQuery(query, cate)
}

func Query(q bson.M) ([]*Category, error) {
	var brs []*Category
	err := TableCategory.UnsafeReadMany(q, &brs)
	return brs, err
}

func GetByParent(parent []string) ([]*Category, error) {
	var category []*Category
	var err = TableCategory.ReadManyIn("parent", parent, &category)
	return category, err
}
