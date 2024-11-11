package handler

import "github.com/sherwin-77/go-echo-template/internal/service"

type RoleHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}
