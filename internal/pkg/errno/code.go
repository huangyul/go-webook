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
	ErrPhoneAlreadyExist    = &Errno{Code: 20103, Message: "Phone already exist"}

	// Code errors

	ErrCodeSendTooFrequent = &Errno{Code: 20103, Message: "Code send too frequent."}
	ErrCodeNotExist        = &Errno{Code: 20104, Message: "Code not exist."}
	ErrCodeVerifyFailed    = &Errno{Code: 20105, Message: "Code verify failed."}
)
