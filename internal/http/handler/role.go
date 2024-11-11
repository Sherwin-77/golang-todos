package handler

import "github.com/sherwin-77/golang-todos/internal/service"

type RoleHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}
