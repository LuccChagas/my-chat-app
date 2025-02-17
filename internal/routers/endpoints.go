package routers

import (
	"github.com/LuccChagas/my-chat-app/internal/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func (router *Router) Endpoints(e *echo.Echo) {

	// users routes
	user := e.Group("/user")
	user.POST("/register", router.User.CreateUserHandler)
	user.GET("/:id", router.User.GetUserHandler)
	user.GET("/all", router.User.GetAllUsersHandler)
	user.POST("/auth", router.User.UserLoginHandler)

	// websocket route
	e.GET("/ws", router.Ws.WsHandler, middleware.AuthMiddleware)

	// html Templates
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login.html", nil)
	})

	e.GET("/chat", func(c echo.Context) error {
		return c.Render(http.StatusOK, "chat.html", nil)
	}, middleware.AuthMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
