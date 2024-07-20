package request

type CreateArticleRequest struct {
	Title      string `json:"title" validate:"required,min=3,max=100"`
	Slug       string `json:"slug" validate:"required,min=3,max=100"`
	Thumbnail  string `json:"thumbnail"`
	Content    string `json:"content" validate:"required"`
	CategoryID uint   `json:"category_id" form:"category_id" validate:"required,category_exists"`
	Tags       []string  `json:"tags" validate:"required"`
}

type UpdateArticleRequest struct {
    Title      *string `json:"title" form:"title"`
    Slug    *string `json:"slug" form:"slug"`
	Thumbnail *string `json:"thumbnail" form:"thumbnail"`
    Content    *string `json:"content" form:"content"`
    CategoryID *uint    `json:"category_id" form:"category_id"`
    Tags       []string `json:"tags" form:"tags"`
}
