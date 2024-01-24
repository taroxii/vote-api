package echoserver

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/taroxii/vote-api/pkg/entity"
	"github.com/taroxii/vote-api/pkg/utils/logger"
	"go.uber.org/zap"
)

type (
	signInRequest struct {
		Username string `json:"username" validate:"required"`
	}
	signInResponse struct {
		Token string `json:"token"`
	}
)

type UserHandler struct {
	userUsecase entity.UserUsecase
}

func NewUserHandler(e *echo.Echo, uuc entity.UserUsecase) {
	handler := &UserHandler{
		userUsecase: uuc,
	}

	e.POST("/sign-in", handler.SignIn)

}

func (handler *UserHandler) SignIn(c echo.Context) error {
	req := new(signInRequest)
	var tokenString *string
	var user *entity.User
	var inheritError error
	if err := c.Bind(req); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, ResponseError{Message: err.Error()})
	}
	ctx := c.Request().Context()
	if tokenString, user, inheritError = handler.userUsecase.SignIn(ctx, req.Username); inheritError != nil {
		return c.JSON(getStatusCode(inheritError), ResponseError{Message: inheritError.Error()})
	}
	logger.Logger.Info("User has been logged in", zap.Any("time", time.Now()), zap.Any("user", *user))

	return c.JSON(http.StatusOK, signInResponse{Token: *tokenString})
}
