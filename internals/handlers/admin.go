package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/riad804/go_ecommerce_api/internals/service"
)

type AdminhHandler struct {
	adminService *service.AdminService
	validate     *validator.Validate
}

func NewAdminHandler(adminService *service.AdminService, validate *validator.Validate) *AdminhHandler {
	return &AdminhHandler{
		adminService: adminService,
		validate:     validate,
	}
}
