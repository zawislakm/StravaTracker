package server

import (
	"app/cmd/web/handlers"
	"net/http"

	"fmt"
	"log"
	"time"

	"app/cmd/web"
	"github.com/a-h/templ"
	"github.com/coder/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/static/*", echo.WrapHandler(fileServer))

	webHandler := handlers.NewHandler(s.db)
	e.GET("/", echo.WrapHandler(http.HandlerFunc(webHandler.HandleIndex)))
	e.GET("/table", echo.WrapHandler(http.HandlerFunc(webHandler.HandleTable)))
	e.GET("/health", echo.WrapHandler(http.HandlerFunc(webHandler.HandleHealth)))
	e.GET("/websocket", s.websocketHandler)
	e.GET("/events", s.eventsHandler)
	return e
}

func (s *Server) eventsHandler(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	c.Response().Header().Set("X-Accel-Buffering", "no") // used so nginx proxy pass events

	fmt.Fprintf(c.Response(), "event: Keep-alive\ndata: First connection\n\n")
	c.Response().Flush()

	sendDateUpdate(c, s.lastUpdate)
	for {
		select {
		case <-c.Request().Context().Done():
			fmt.Println("Client connection closed with reason:", c.Request().Context().Err())

			return nil
		case newData := <-s.newActivitiesChan:
			fmt.Printf("New activities: %v\n", newData)
			if newData {
				fmt.Fprintf(c.Response(), "event: Table\ndata: New data athletes data to fetch\n\n")
			}
			sendDateUpdate(c, s.lastUpdate)
		case <-time.After(10 * time.Second):
			fmt.Fprintf(c.Response(), "event: Keep-alive\ndata: connected\n\n")
			c.Response().Flush()
		}
	}
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) websocketHandler(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
	return nil
}
