package routers

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/LuccChagas/my-chat-app/internal/handlers"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/LuccChagas/my-chat-app/docs/app"
)

type Router struct {
	User handlers.UserHandlerInterface
	Ws   handlers.WsHandlerInterface
}

func NewRouter(
	user handlers.UserHandlerInterface,
	ws handlers.WsHandlerInterface,

) *Router {
	return &Router{
		User: user,
		Ws:   ws,
	}
}

type Template struct {
	templates *template.Template
}

// Render implementa o método Render da interface echo.Renderer.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// @title My Chat App API
// @version 1.0
// @description API documentation for My Chat App, a real-time chat application.
// @contact.name API Support
// @contact.url https://www.linkedin.com/in/luccas-machado-ab5897105/
// @contact.email luccaa.chagas23@gmail.com
// @host localhost:1323
// @BasePath /
// @schemes http

func (router *Router) Serve() {
	authKeyBase64 := os.Getenv("SESSION_AUTH_KEY")
	encKeyBase64 := os.Getenv("SESSION_ENC_KEY")

	// Decodificar as chaves de Base64 para []byte reais
	authKey, err := base64.StdEncoding.DecodeString(authKeyBase64)
	if err != nil {
		log.Fatalf("Erro on decode SESSION_AUTH_KEY: %v", err)
	}

	encKey, err := base64.StdEncoding.DecodeString(encKeyBase64)
	if err != nil {
		log.Fatalf("Erro on decode SESSION_ENC_KEY: %v", err)
	}

	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore(authKey, encKey)))
	router.Endpoints(e)

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(":1323"); err != nil && !errors.Is(http.ErrServerClosed, err) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
