package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation      = &Errno{Code: 20001, Message: "Validation failed."}
	ErrValidationTopic = &Errno{Code: 20002, Message: "Validation failed topic can not be empty."}
	ErrValidationDelay = &Errno{Code: 20003, Message: "Validation failed delay should be range from 1 to (2^31 - 1)."}
	ErrValidationBody  = &Errno{Code: 20004, Message: "Validation failed body can not be empty."}
	ErrDatabase        = &Errno{Code: 20005, Message: "Database error."}
	ErrToken           = &Errno{Code: 20006, Message: "Error occurred while signing the JSON web token."}

	// user errors
	ErrEncrypt           = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}
)
