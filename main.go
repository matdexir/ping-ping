package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matdexir/ping-ping/controllers"
	"github.com/matdexir/ping-ping/db"
)

func main() {
	e := echo.New()

	database, err := db.CreateConnection()
	if err != nil {
		os.Exit(1)
	}
	err = database.CreateTable()
	if err != nil {
		os.Exit(1)
	}
	database.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/v1/ad", controllers.CreateSponsoredPost)
	e.GET("/api/v1/ad", controllers.GetSponsoredPost)
	e.Logger.Fatal(e.Start(":8080"))
}
