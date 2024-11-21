package errno

import "net/http"

var (
	// Common errors
	ErrOK             = &Errno{HTTP: http.StatusOK, Code: 0, Message: "OK"}
	ErrInternalServer = &Errno{HTTP: http.StatusInternalServerError, Code: 10001, Message: "Internal server error"}
	ErrBadRequest     = &Errno{HTTP: http.StatusBadRequest, Code: 10002, Message: "Bad request"}

	// User errors
	ErrEmailAlreadyExist = &Errno{HTTP: http.StatusBadRequest, Code: 20100, Message: "Email already exist"}
	ErrNotFoundUser      = &Errno{HTTP: http.StatusBadRequest, Code: 20101, Message: "The user was not found."}
)
