package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type echoServer struct {
	e *echo.Echo
}

var server echoServer

func StartApi() {
	server.e = echo.New()
	server.e.HideBanner = true
	echo.NotFoundHandler = func(c echo.Context) error {
		return c.String(http.StatusNotFound, "")
	}

	server.e.Use(middleware.Logger(), middleware.Recover())

	server.e.GET("/", home)
	server.e.POST("/verify", verifyPost)

	image := server.e.Group("/images")
	{
		image.GET("/", imageInstructions)
		image.GET("/:id", getImage)
		image.GET("/list", getImageList)
		image.GET("/upload", imageUploadInstructions)
		image.GET("/delete", imageDeleteInstructions)

		image.POST("/upload", imageUpload, imageAuthPerm)
		image.POST("/delete", imageDelete, imageAuthPerm)
	}

	go func() {
		fmt.Print("Api is now running")
		if err := server.e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			server.e.Logger.Fatal("shutting down the server")
		}
	}()
}

func Shutdown(ctx context.Context) {
	if err := server.e.Shutdown(ctx); err != nil {
		server.e.Logger.Fatal(err)
	}
}
