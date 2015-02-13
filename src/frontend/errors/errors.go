package errors

type FrontendError struct {
	ErrorText string
	HttpCode  int
	ErrorCode int
}

func (f *FrontendError) Error() string {
	return f.ErrorText
}

var WARNING_NOT_LOGGED_IN = FrontendError{
	ErrorText: "User not logged in.",
	HttpCode:  400,
	ErrorCode: 100,
}

var INVALID_POST_DATA = FrontendError{
	ErrorText: "Invalid parameters sent for POST request.",
	HttpCode:  400,
	ErrorCode: 101,
}

var CORRUPTED_SESSION = FrontendError{
	ErrorText: "Your session has been corrupted. Please log in again.",
	HttpCode:  400,
	ErrorCode: 102,
}

var USER_DOESNT_EXIST = FrontendError{
	ErrorText: "The requested user doesn't exist.",
	HttpCode:  500,
	ErrorCode: 103,
}

var MISSING_QUERY_PARAMS = FrontendError{
	ErrorText: "Invalid or missing values for required query parameters",
	HttpCode:  400,
	ErrorCode: 104,
}

var HANDLER_NOT_IMPLEMENTED = FrontendError{
	ErrorText: "This API function hasn't been implemented yet.",
	HttpCode:  501,
	ErrorCode: 200,
}

var NOT_LOGGED_IN = FrontendError{
	ErrorText: "The requested user is not logged in.",
	HttpCode:  401,
	ErrorCode: 201,
}

var COULDNT_COMPLETE_OPERATION = FrontendError{
	ErrorText: "Couldn't complete the requested operation.",
	HttpCode:  401,
	ErrorCode: 301,
}
