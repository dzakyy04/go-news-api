package request

type TagRequest struct {
	Name string `json:"name" validate:"required"`
}
