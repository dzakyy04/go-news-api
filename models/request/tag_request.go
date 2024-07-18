package request

type CreateTagRequest struct {
	Name string `json:"name" validate:"required"`
}
