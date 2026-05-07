package lifecycle

import "fmt"

type DocStatus int

const (
	Draft     DocStatus = 0
	Submitted DocStatus = 1
	Cancelled DocStatus = 2
)

func (s DocStatus) Int() int {
	return int(s)
}

func (s DocStatus) String() string {
	switch s {
	case Draft:
		return "Draft"
	case Submitted:
		return "Submitted"
	case Cancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

func IsValid(value int) bool {
	switch DocStatus(value) {
	case Draft, Submitted, Cancelled:
		return true
	default:
		return false
	}
}

func Parse(value int) (DocStatus, error) {
	status := DocStatus(value)

	if !IsValid(value) {
		return Draft, fmt.Errorf("invalid docstatus: %d", value)
	}

	return status, nil
}

func CanSubmit(value int) bool {
	return DocStatus(value) == Draft
}

func CanCancel(value int) bool {
	return DocStatus(value) == Submitted
}

func CanAmend(value int) bool {
	return DocStatus(value) == Cancelled
}

func CanEdit(value int, isSubmittable bool) bool {
	status := DocStatus(value)

	if !isSubmittable {
		return status != Cancelled
	}

	return status == Draft
}
