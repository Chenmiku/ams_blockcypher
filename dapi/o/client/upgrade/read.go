package upgrade

import ()
import "gopkg.in/mgo.v2/bson"

func GetByID(id string) (*Upgrade, error) {
	var u Upgrade
	return &u, TableUpgrade.ReadByID(id, &u)
}

func GetLastest() (*Upgrade, error) {

	var u Upgrade

	pipeline := []bson.M{
		bson.M{
			"$sort": bson.M{
				"mtime": -1,
			},
		},
	}

	return &u, TableUpgrade.C().Pipe(pipeline).One(&u)
}
