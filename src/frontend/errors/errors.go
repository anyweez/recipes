package errors


type FrontendError struct {
	ErrorText	string
	HttpCode	int
	ErrorCode	int
}
	
func (f *FrontendError) Error() string {
	return f.ErrorText
}

var WARNING_NOT_LOGGED_IN = FrontendError{
	ErrorText: "User not logged in.",
	HttpCode: 400,
	ErrorCode: 100, 
}

var INVALID_POST_DATA = FrontendError{
	ErrorText: "Invalid parameters sent for POST request.",
	HttpCode: 400,
	ErrorCode: 101,
}

var CORRUPTED_SESSION = FrontendError{
	ErrorText: "Your session has been corrupted. Please log in again.",
	HttpCode: 400,
	ErrorCode: 102,
}

var HANDLER_NOT_IMPLEMENTED = FrontendError{
	ErrorText: "This API function hasn't been implemented yet.",
	HttpCode: 501,
	ErrorCode: 200,
}
