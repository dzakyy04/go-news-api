package request

type CreateArticleRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Slug        string `json:"slug" validate:"required,min=3,max=100"`
	Thumbnail   string `json:"thumbnail"`
	Content     string `json:"content" validate:"required"`
	CategoryID  uint   `json:"category_id" form:"category_id" validate:"required,category_exists"`
	AuthorID    uint   `json:"author_id" form:"author_id" validate:"required,author_exists"`
}