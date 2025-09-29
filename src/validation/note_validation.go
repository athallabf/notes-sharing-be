package validation

type CreateNote struct {
	Title   string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"omitempty"`
}

type UpdateNote struct {
	Title   string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"omitempty"`
}
