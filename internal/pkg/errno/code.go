package errno

var (
	// Common errors
	ErrOK             = &Errno{Code: 0, Message: "OK"}
	ErrInternalServer = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBadRequest     = &Errno{Code: 10002, Message: "Bad request"}

	// User errors
	ErrEmailAlreadyExist    = &Errno{Code: 20100, Message: "Email already exist"}
	ErrNotFoundUser         = &Errno{Code: 20101, Message: "The user was not found."}
	ErrEmailOrPasswordError = &Errno{Code: 20102, Message: "The email or password error."}
)
