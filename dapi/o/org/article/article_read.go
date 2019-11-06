package article

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Article, error) {
	var b Article
	return &b, TableArticle.ReadByID(id, &b)
}

func GetParent(id string) (*Article, error) {
	var b Article
	return &b, TableArticle.C().Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"parent": id,
			},
		},
	}).One(&b)
}

// func Search(articleIDs []string) ([]*Article, error) {
// 	var article = []*Article{}
// 	var err = TableArticle.ReadManyIn("article_id", articleIDs, &article)
// 	return article, err
// }

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, article *[]Article) (int, error) {
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

	return TableArticle.ReadPagingByQuery(query, article)
}

func GetArticleByCategoryID(categoryID string, pageSize int, pageNumber int, sortBy string, sortOrder string, article *[]Article) (int, error) {
	var where bson.M
	query := mgo.PagingQuery{}
	query.PageSize = pageSize
	query.PageNumber = pageNumber
	query.SortBy = sortBy
	query.SortOrder = sortOrder

	if categoryID != "" {
		where = bson.M{
			"dtime":       0,
			"active":      true,
			"category_id": categoryID,
		}
	} else {
		where = bson.M{
			"dtime":  0,
			"active": true,
		}
	}

	query.Pipeline = append(query.Pipeline,
		bson.M{
			"$match": where,
		})

	return TableArticle.ReadPagingByQuery(query, article)
}
