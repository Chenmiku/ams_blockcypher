package contact

import (
	"db/mgo"
	"gopkg.in/mgo.v2/bson"
)

func GetByID(id string) (*Contact, error) {
	var b Contact
	return &b, TableContact.ReadByID(id, &b)
}

func GetParent(id string) (*Contact, error) {
	var b Contact
	return &b, TableContact.C().Pipe([]bson.M{
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

func GetAll(pageSize int, pageNumber int, sortBy string, sortOrder string, contact *[]Contact) (int, error) {
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

	return TableContact.ReadPagingByQuery(query, contact)
}

func Query(q bson.M) ([]*Contact, error) {
	var brs []*Contact
	err := TableContact.UnsafeReadMany(q, &brs)
	return brs, err
}
