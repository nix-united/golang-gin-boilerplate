package response

import "basic_server/server/model"

type CreatePostResponse struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostResponse struct {
	Id      uint   `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CollectionResponse struct {
	Collection interface{} `json:"collection"`
	Meta       Meta `json:"meta"`
}

type Meta struct {
	Amount int `json:"amount"`
}

func CreatePostsCollectionResponse(posts []model.Post) CollectionResponse {
	collection := make([]GetPostResponse, 0)

	for _, post := range posts {
		collection = append(collection, GetPostResponse{
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
		})
	}
	return CollectionResponse{Collection: collection, Meta: Meta{Amount: len(collection)}}
}
