package assets

type errNotFound struct {
	url string
}

func (e errNotFound) Error() string {
	return "document not found: " + e.url
}

type internalError struct {
	cause error
}

func (i internalError) Error() string {
	return "internal error: " + i.cause.Error()
}

func (i internalError) Unwrap() error {
	return i.cause
}


