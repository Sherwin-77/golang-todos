package dto

type RoleRequest struct {
	Name      string `json:"name" validate:"required"`
	AuthLevel int    `json:"auth_level" validate:"required"`
}

type UpdateRoleRequest struct {
	RoleRequest
	ID string `param:"id" validate:"required,uuid"`
}

type ChangeRoleRequest struct {
	UserID string                  `param:"id" validate:"required,uuid"`
	Items  []ChangeRoleRequestItem `json:"items" validate:"required"`
}

type ChangeRoleRequestItem struct {
	ID     string `json:"id" validate:"required,uuid"`
	Action string `json:"action" validate:"required,oneof=add remove"`
}
