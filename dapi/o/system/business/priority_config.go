package business

import (
	"fmt"
)

type PriorityConfig struct {
	PriorityStep             int `bson:"priority_step" json:"priority_step"`
	InternalVipCard          int `bson:"internal_vip_card" json:"internal_vip_card"`
	CustomerVipCard          int `bson:"customer_vip_card" json:"customer_vip_card"`
	PrivilegedCustomer       int `bson:"privileged_customer" json:"privileged_customer"`
	MovedTicket              int `bson:"moved_ticket" json:"moved_ticket"`
	BookedTicket             int `bson:"booked_ticket" json:"booked_ticket"`
	RestoreTicket            int `bson:"restore_ticket" json:"restore_ticket"`
	MinPriorityRestricted    int `bson:"min_priority_restricted" json:"min_priority_restricted"`
	MinPriorityUnorderedCall int `bson:"min_priority_unordered_call" json:"min_priority_unordered_call"`
}

func (c PriorityConfig) String() string {
	return fmt.Sprintf(
		"priority:step=%d;vi=%d;vc=%d;pc=%d;mo=%d;bo=%d;re=%d;mps=%d;mpuc=%d",
		c.PriorityStep,
		c.InternalVipCard, c.CustomerVipCard, c.PrivilegedCustomer,
		c.MovedTicket, c.BookedTicket, c.RestoreTicket, c.MinPriorityRestricted, c.MinPriorityUnorderedCall,
	)
}

var defaultPriorityConfig = &PriorityConfig{
	PriorityStep:             0,
	MinPriorityRestricted:    1 << 10,
	MinPriorityUnorderedCall: 1,
	InternalVipCard:          1,
	CustomerVipCard:          1,
	MovedTicket:              1,
	BookedTicket:             1,
	RestoreTicket:            1,
	PrivilegedCustomer:       1,
}

func (c *PriorityConfig) Check() {
	c.Inherit(nil)
}

func (c *PriorityConfig) Inherit(c2 *PriorityConfig) *PriorityConfig {
	if c2 == nil {
		c2 = defaultPriorityConfig
	}
	if c.PriorityStep < 0 {
		c.PriorityStep = c2.PriorityStep
	}
	if c.MinPriorityRestricted < 1 {
		// no restriction
		c.MinPriorityRestricted = c2.MinPriorityRestricted
	}
	if c.MinPriorityUnorderedCall < 1 {
		// no unordered call
		c.MinPriorityUnorderedCall = c2.MinPriorityUnorderedCall
	}

	// customer
	if c.InternalVipCard < 0 {
		c.InternalVipCard = c2.InternalVipCard
	}
	if c.CustomerVipCard < 0 {
		c.CustomerVipCard = c2.CustomerVipCard
	}
	if c.RestoreTicket < 0 {
		c.RestoreTicket = c2.RestoreTicket
	}
	if c.BookedTicket < 0 {
		c.BookedTicket = c2.BookedTicket
	}
	//

	if c.MovedTicket < 0 {
		c.MovedTicket = c2.MovedTicket
	}
	if c.PrivilegedCustomer < 0 {
		c.PrivilegedCustomer = c2.PrivilegedCustomer
	}
	return c
}
