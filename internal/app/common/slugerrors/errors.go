package slugerrors

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown       = ErrorType{"unknown"}
	ErrorTypeAuthorization = ErrorType{"authorization"}
	ErrorTypeBadRequest    = ErrorType{"bad-request"}
	ErrorTypeNotFound      = ErrorType{"not-found"}
)

type SlugError struct {
	message   string
	slug      string
	errorType ErrorType
}

func (s SlugError) Error() string {
	return s.message
}

func (s SlugError) Slug() string {
	return s.slug
}

func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

func NewSlugError(errMsg string, slug string) SlugError {
	return SlugError{
		message:   errMsg,
		slug:      slug,
		errorType: ErrorTypeUnknown,
	}
}

func NewAuthorizationError(errMsg string, slug string) SlugError {
	return SlugError{
		message:   errMsg,
		slug:      slug,
		errorType: ErrorTypeAuthorization,
	}
}

func NewBadRequestError(errMsg string, slug string) SlugError {
	return SlugError{
		message:   errMsg,
		slug:      slug,
		errorType: ErrorTypeBadRequest,
	}
}

func NewNotFoundError(errMsg string, slug string) SlugError {
	return SlugError{
		message:   errMsg,
		slug:      slug,
		errorType: ErrorTypeNotFound,
	}
}
