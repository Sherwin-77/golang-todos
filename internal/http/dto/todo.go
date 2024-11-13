package dto

type TodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

type UpdateTodoRequest struct {
	ID          string `param:"id" validate:"required,uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
