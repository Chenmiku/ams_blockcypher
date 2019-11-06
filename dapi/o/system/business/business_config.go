package business

import (
	"db/mgo"
	"fmt"
)

type BusinessConfig struct {
	mgo.BaseModel `bson:",inline"`
	BranchID      string         `bson:"branch_id" json:"branch_id"`
	General       GeneralConfig  `bson:"general" json:"general"`
	Service       ServiceConfig  `bson:"service" json:"service"`
	Priority      PriorityConfig `bson:"priority" json:"priority"`
	Counter       CounterConfig  `bson:"counter" json:"counter"`
}

func (c BusinessConfig) String() string {
	return fmt.Sprintf("business=[%s][%s][%s][%s]", c.General, c.Service, c.Priority, c.Counter)
}

func (c *BusinessConfig) Check() {
	c.General.Check()
	c.Service.Check()
	c.Priority.Check()
	c.Counter.Check()
}
