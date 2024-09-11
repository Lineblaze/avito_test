package errors

type ErrorType string

const (
	ErrorTypeInternal     ErrorType = "INTERNAL"
	ErrorTypeInvalidInput ErrorType = "INVALID_INPUT"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeAuth         ErrorType = "AUTH"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeValidation   ErrorType = "VALIDATION"
)

type Error struct {
	error     string
	slug      string
	errorType ErrorType
}

func (e Error) Error() string {
	return e.error
}

func (e Error) Slug() string {
	return e.slug
}

func (e Error) Type() ErrorType {
	return e.errorType
}

func NewInternalError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeInternal,
	}
}

func NewInvalidInputError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeInvalidInput,
	}
}

func NewNotFoundError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeNotFound,
	}
}

func NewConflictError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeConflict,
	}
}

func NewAuthError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeAuth,
	}
}

func NewForbiddenError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeForbidden,
	}
}

func NewValidationError(error, slug string) Error {
	return Error{
		error:     error,
		slug:      slug,
		errorType: ErrorTypeValidation,
	}
}
