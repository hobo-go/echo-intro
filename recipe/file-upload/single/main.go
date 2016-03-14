package main

import (
	"fmt"
	"io"
	"os"

	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

func upload() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Read form fields
		name := c.Form("name")
		email := c.Form("email")

		//-----------
		// Read file
		//-----------

		// Source
		file, err := req.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded successfully with fields name=%s and email=%s.</p>", file.Filename, name, email))
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static("public"))

	e.Post("/upload", upload())

	e.Run(standard.New(":1323"))
}
