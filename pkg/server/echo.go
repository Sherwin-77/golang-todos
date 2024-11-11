package server

import "github.com/labstack/echo/v4"

type Server struct {
	*echo.Echo
}

func NewServer() *Server {
	e := echo.New()
	e.HideBanner = true
	return &Server{e}
}
