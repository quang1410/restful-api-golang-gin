package handler

import (
	"galvin/lession05-exercise-user-management/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) GetAllUsers(ctx *gin.Context) {

}

func (h *UserHandler) CreateUser(ctx *gin.Context) {

}

func (h *UserHandler) GetUserByUUID(ctx *gin.Context) {

}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {

}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {

}
