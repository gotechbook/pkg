package errors

type Status struct {
	Code     int32
	Reason   string
	Message  string
	Metadata map[string]string
}

func (s *Status) GetCode() int32 {
	return s.Code
}

func (s *Status) GetReason() string {
	return s.Reason
}

func (s *Status) GetMessage() string {
	return s.Message
}

func (s *Status) GetMetadata() map[string]string {
	return s.Metadata
}
