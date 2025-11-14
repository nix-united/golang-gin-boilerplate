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

func NewCollectionResponse[T any](items []T) CollectionResponse[T] {
	return CollectionResponse[T]{
		Data: items,
		Meta: Meta{Amount: len(items)},
	}
}

type Meta struct {
	Amount int `json:"amount"`
}
