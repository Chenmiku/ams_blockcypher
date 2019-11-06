package business

import (
	"fmt"
)

type KioskConfig struct {
	MaxTicket     int        `bson:"max_ticket" json:"max_ticket"`
	TimeWorkInDay []TimeWork `bson:"time_work_in_day" json:"time_work_in_day"`
}

type TimeWork struct {
	TimeStart float32 `bson:"time_start" json:"time_start"`
	TimeEnd   float32 `bson:"time_end" json:"time_end"`
}

func (k KioskConfig) String() string {
	return fmt.Sprintf(
		"MaxTicket:so_ve=%d",
		k.MaxTicket)
}

var defaultKioskConfig = &KioskConfig{
	MaxTicket: 40,
	TimeWorkInDay: []TimeWork{TimeWork{
		TimeStart: 8.5,
		TimeEnd:   10.5,
	}, TimeWork{
		TimeStart: 13,
		TimeEnd:   16,
	},
	},
}
var defaultTimeWork = TimeWork{}

func (c *KioskConfig) Check() {
	c.Inherit(nil)
}

func (c *KioskConfig) Inherit(cc *KioskConfig) {
	if cc == nil {
		cc = defaultKioskConfig
	}
	if c.MaxTicket <= 0 {
		c.MaxTicket = cc.MaxTicket
	}
	if len(c.TimeWorkInDay) <= 0 {
		c.TimeWorkInDay = cc.TimeWorkInDay
	}
}
