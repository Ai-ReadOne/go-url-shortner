package controller

// This is a custom error type that implements the BadRequester interface.
// It is used to return a 400 Bad Request error response to the client.
type badRequest struct {
	error
}

func (badRequest) BadRequest() {}

// This is a custom error type that implements the NotFounder interface.
// It is used to return a 404 Bad Request error response to the client.
type notFound struct {
	error
}

func (notFound) NotFound() {}
