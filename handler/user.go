package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
)

func (handler *Handler) SignUp(context echo.Context) error {
	var u model.User
	req := &userRegisterRequest{}
	if err := req.bind(context, &u); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	if err := handler.userStore.Create(&u); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	return utils.ResponseByContentType(context, http.StatusCreated, newUserResponse(&u))
}

func (handler *Handler) Login(context echo.Context) error {
	req := &userLoginRequest{}
	if err := req.bind(context); err != nil {
		return utils.ResponseByContentType(context, http.StatusUnprocessableEntity, utils.NewError(err))
	}
	u, err := handler.userStore.GetByEmail(req.User.Email)
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return utils.ResponseByContentType(context, http.StatusForbidden, utils.AccessForbidden())
	}
	if !u.CheckPassword(req.User.Password) {
		return utils.ResponseByContentType(context, http.StatusForbidden, utils.AccessForbidden())
	}
	return utils.ResponseByContentType(context, http.StatusOK, newUserResponse(u))
}

func userIDFromToken(context echo.Context) uint {
	id, ok := context.Get("user").(uint)
	if !ok {
		return 0
	}
	return id
}

func (handler *Handler) CurrentUser(context echo.Context) error {
	u, err := handler.userStore.GetByID(userIDFromToken(context))
	if err != nil {
		return utils.ResponseByContentType(context, http.StatusInternalServerError, utils.NewError(err))
	}
	if u == nil {
		return utils.ResponseByContentType(context, http.StatusNotFound, utils.NotFound())
	}
	return utils.ResponseByContentType(context, http.StatusOK, newUserResponse(u))
}
