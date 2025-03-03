package vectorclock

type Occurred uint8

func (oc Occurred) String() string {
	if oc == BEFORE {
		return "Before"
	}
	if oc == AFTER {
		return "After"
	}
	return "Concurrently"
}

const (
	BEFORE Occurred = iota
	AFTER
	CONCURRENTLY
)
