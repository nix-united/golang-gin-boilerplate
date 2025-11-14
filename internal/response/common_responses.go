package response

type MessageResponse struct {
	Message string `json:"message"`
}

func NewMessageResponse(message string) MessageResponse {
	return MessageResponse{Message: message}
}

type CollectionResponse[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}

func NewCollectionResponse[T any](items []T, total, offset, limit int) CollectionResponse[T] {
	return CollectionResponse[T]{
		Data: items,
		Meta: Meta{
			Count:  len(items),
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
}

type Meta struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ErrorResponseCode string

const (
	CodeBadRequest          ErrorResponseCode = "bad_request"
	CodeNotFound            ErrorResponseCode = "not_found"
	CodeAlreadyExists       ErrorResponseCode = "already_exists"
	CodeAccessDenied        ErrorResponseCode = "access_denied"
	CodeInternalServerError ErrorResponseCode = "internal_server_error"
)

type ErrorResponse struct {
	Code    ErrorResponseCode `json:"code"`
	Message string            `json:"message"`
}

func NewErrorResponse(code ErrorResponseCode, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Message: message,
	}
}
