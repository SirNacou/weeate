package domain

import (
	"database/sql/driver"
	"fmt"
)

type PollStrategy string

const (
	OrderMultiple  PollStrategy = "ORDER_PERSONAL_CHOICE"
	OrderConsensus PollStrategy = "ORDER_CONSENSUS_ITEM"
)

// Scan tells GORM how to read a string from the DB
// into your PollStrategy type.
func (ps *PollStrategy) Scan(value any) error {
	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if ok {
			str = string(bytes)
		} else {
			return fmt.Errorf("failed to scan type: %T", value)
		}
	}
	*ps = PollStrategy(str)
	return nil
}

// Value tells GORM how to write your PollStrategy type
// to the DB as a plain string.
func (ps PollStrategy) Value() (driver.Value, error) {
	return string(ps), nil
}
