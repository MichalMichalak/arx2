package svc

type Status string

const (
	Created Status = "CREATED"
	Running Status = "RUNNING"
	Closing Status = "CLOSING"
	Closed  Status = "CLOSED"
)

var validStatuses = map[Status]struct{}{Created: {}, Running: {}}

func (s Status) Valid() bool {
	_, valid := validStatuses[s]
	return valid
}

func (s Status) IsCreated() bool {
	return s == Created
}
