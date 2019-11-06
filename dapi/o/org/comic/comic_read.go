package comic

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Comic, error) {
	var b Comic
	return &b, TableComic.ReadByID(id, &b)
}

// func Search(branchIDs []string) ([]*Branch, error) {
// 	var branch = []*Branch{}
// 	var err = TableBranch.ReadManyIn("branch_id", branchIDs, &branch)
// 	return branch, err
// }

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, cate *[]Comic) (int, error) {
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

	return TableComic.ReadPagingByQuery(query, cate)
}

func GetByCategory(pageSize int, pageNumber int, sortBy string, sortOrder string, cate []string) ([]Comic, error) {
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	var res []Comic

	where := bson.M{
		"dtime":      0,
		"categories": bson.M{"$in": cate},
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return res, TableComic.ReadPaginationByQuery(query, &res)
}

func Query(q bson.M) ([]*Comic, error) {
	var brs []*Comic
	err := TableComic.UnsafeReadMany(q, &brs)
	return brs, err
}

func GetByType(pageSize int, pageNumber int, sortBy string, sortOrder string, types string) ([]Comic, error) {
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	var where bson.M
	var res []Comic

	if types == "hot" {
		where = bson.M{
			"dtime": 0,
			"hot":   true,
		}
	} else if types == "new" {
		where = bson.M{
			"dtime": 0,
			"new":   true,
		}
	} else if types == "top" {
		where = bson.M{
			"dtime": 0,
			"top":   true,
		}
	} else {
		where = bson.M{
			"dtime": 0,
			"hot":   false,
			"new":   false,
			"top":   false,
		}
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return res, TableComic.ReadPaginationByQuery(query, &res)
}
