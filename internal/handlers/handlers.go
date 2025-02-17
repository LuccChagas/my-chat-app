package handlers

import "github.com/labstack/echo/v4"

type UserHandlerInterface interface {
	CreateUserHandler(c echo.Context) error
	GetAllUsersHandler(c echo.Context) error
	GetUserHandler(c echo.Context) error
	UserLoginHandler(c echo.Context) error
}

type WsHandlerInterface interface {
	WsHandler(echo.Context) error
}
