package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	StatusNotFound      = &Errno{Code: 10001, Message: "The incorrect API route."}
	InternalServerError = &Errno{Code: 10002, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10003, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation      = &Errno{Code: 20001, Message: "Validation failed."}
	ErrValidationTopic = &Errno{Code: 20002, Message: "Validation failed topic can not be empty."}
	ErrValidationDelay = &Errno{Code: 20003, Message: "Validation failed delay should be range from 1 to (2^31 - 1)."}
	ErrValidationBody  = &Errno{Code: 20004, Message: "Validation failed body can not be empty."}
)
