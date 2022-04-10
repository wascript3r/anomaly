package request

type Validate interface {
	RawRequest(s interface{}) error
	GetDateTimeFormat() string
}
