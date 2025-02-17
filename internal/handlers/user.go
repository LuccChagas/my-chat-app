package handlers

import (
	"fmt"
	"github.com/LuccChagas/my-chat-app/internal/models"
	"github.com/LuccChagas/my-chat-app/internal/services"
	"github.com/LuccChagas/my-chat-app/utils"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(u *services.UserService) *UserHandler {
	return &UserHandler{
		service: u,
	}
}

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create a new user with the provided details.
// @Tags User
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User Request"
// @Success 201 {object} models.UserResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/register [post]
func (h *UserHandler) CreateUserHandler(c echo.Context) error {
	var user models.UserRequest
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := utils.Validate(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("An error ocurred while validating user data: %s", err.Error()))
	}

	response, err := h.service.CreateUser(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

// GetAllUsersHandler godoc
// @Summary Get all users
// @Description Retrieve a list of all registered users.
// @Tags User
// @Produce json
// @Success 200 {array} models.UserResponse
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/all [get]
func (h *UserHandler) GetAllUsersHandler(c echo.Context) error {
	response, err := h.service.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

// GetUserHandler godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID.
// @Tags User
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/{id} [get]
func (h *UserHandler) GetUserHandler(c echo.Context) error {
	ID := c.QueryParam("id")
	if len(ID) < 36 {
		return c.JSON(http.StatusBadRequest, "No user ID provided")
	}

	parsedID, err := uuid.Parse(ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error parsing user ID - Get User")
	}

	response, err := h.service.GetUser(c.Request().Context(), parsedID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

// UserLoginHandler godoc
// @Summary User login
// @Description Authenticates a user and creates a session.
// @Tags User
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nickname formData string true "User Nickname"
// @Param password formData string true "User Password"
// @Success 302 {string} string "Redirect to /chat"
// @Failure 400 {string} string "Invalid credentials or Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /user/auth [post]

func (h *UserHandler) UserLoginHandler(c echo.Context) error {
	login := models.UserLoginRequest{
		Nickname: c.FormValue("nickname"),
		Password: c.FormValue("password"),
	}

	response, err := h.service.GetUserByUsername(c.Request().Context(), login.Nickname)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if !utils.CheckPasswordHash(login.Password, response.Password) {
		return c.JSON(http.StatusBadRequest, "Invalid credentials")
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Invalid session")
	}

	sess.Values["nickname"] = response.NickName

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	if err = sess.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/chat")
}
