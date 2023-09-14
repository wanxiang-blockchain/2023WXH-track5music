package response

var (
	// ErrSuccess common errors
	ErrSuccess             = newError(0, "ok")
	ErrBadRequest          = newError(400, "Bad Request")
	ErrUnauthorized        = newError(401, "Unauthorized")
	ErrNotFound            = newError(404, "Not Found")
	ErrInternalServerError = newError(500, "Internal Server Error")

	// ErrUsernameAlreadyUse more biz errors
	ErrUsernameAlreadyUse = newError(1001, "The username is already in use.")
)
