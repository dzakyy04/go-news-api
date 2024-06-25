package request

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=500"`
}
