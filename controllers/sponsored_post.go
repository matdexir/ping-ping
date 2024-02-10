package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateSponsoredPost(c echo.Context) error {
	return c.String(http.StatusOK, "New sponsored post created")
}

func GetSponsoredPost(c echo.Context) error {
	return c.String(http.StatusOK, "Here are your posts")
}
