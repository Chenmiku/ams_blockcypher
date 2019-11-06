package business

import (
	"fmt"
)

type OverrideMode string

func (v *OverrideMode) Inherit(v2 OverrideMode) {
	switch *v {
	case AlwayOn:
	case AlwayOff:
	default:
		// inherit
		*v = v2
	}
}

func (v *OverrideMode) Check() {
	switch *v {
	case AlwayOn:
	case AlwayOff:
	case Inherit:
	default:
		*v = Inherit
	}
}

const (
	AlwayOn  = OverrideMode("alway_on")
	AlwayOff = OverrideMode("alway_off")
	Inherit  = OverrideMode("inherit")
)

type CounterConfig struct {
	RecordTransaction OverrideMode `bson:"record_transaction" json:"record_transaction"`
}

func (c CounterConfig) String() string {
	return fmt.Sprintf("counter:rec=%s", c.RecordTransaction)
}

var defaultCounterConfig = &CounterConfig{
	RecordTransaction: AlwayOff,
}

func (c *CounterConfig) Check() {
	c.Inherit(nil)
}

func (c *CounterConfig) Inherit(cc *CounterConfig) {
	if cc == nil {
		cc = defaultCounterConfig
	}
	c.RecordTransaction.Inherit(cc.RecordTransaction)
}
