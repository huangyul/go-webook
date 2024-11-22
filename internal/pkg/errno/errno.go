package errno

type Errno struct {
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

func (err *Errno) Decode() (int, string) {
	return err.Code, err.Message
}
