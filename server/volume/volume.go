package volume

type Status int

const (
	Available Status = iota
	Unavailable
	Attaching
	Attached
	Detaching
)

type Volume struct {
	ID     string
	Region string
	Tag    string
	Status Status
}
