package services

import "fmt"

// Policy controls service-level behavior, such as read-only mode.
type Policy struct {
	ReadOnly bool
}

func NewPolicy(readOnly bool) Policy {
	return Policy{
		ReadOnly: readOnly,
	}
}

// IsWriteAllowed returns an error when the server is in read-only mode.
func (p Policy) IsWriteAllowed() error {
	if p.ReadOnly {
		return fmt.Errorf("write operation is disabled in read-only mode")
	}

	return nil
}
