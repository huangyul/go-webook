package errno

type Errno struct {
	HTTP    int    `json:"http"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(message string) *Errno {
	err.Message = message
	return err
}

func Decode(err error) *Errno {
	if err == nil {
		return ErrOK
	}
	switch v := err.(type) {
	case *Errno:
		return v
	default:
		return ErrInternalServer
	}
}
