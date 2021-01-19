package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shkryob/goforum/model"
	"github.com/shkryob/goforum/utils"
)

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userRegisterRequest true "User info for registration"
// @Success 201 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users [post]
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

// Login godoc
// @Summary Login for existing user
// @Description Login for existing user
// @ID login
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body userLoginRequest true "Credentials to use"
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users/login [post]
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

// CurrentUser godoc
// @Summary Get the current user
// @Description Gets the currently logged-in user
// @ID current-user
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} userResponse
// @Failure 400 {object} utils.Error
// @Failure 401 {object} utils.Error
// @Failure 422 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Security ApiKeyAuth
// @Router /user [get]
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
