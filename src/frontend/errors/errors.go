package errors


type FrontendError struct {
	ErrorText	string
	HttpCode	int
	ErrorCode	int
}

var WARNING_NOT_LOGGED_IN = FrontendError{
	ErrorText: "User not logged in.",
	HttpCode: 400,
	ErrorCode: 100, 
}
