package request

type CrudTagRequest struct {
	Name string `json:"name" validate:"required"`
}
