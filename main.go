package main

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Gender string

const (
	M Gender = "M"
	F        = "F"
)

type Country string

const (
	JP Country = "JP"
	TW         = "TW"
	US         = "US"
	BR         = "BR"
	SA         = "SA"
	FR         = "FR"
)

type Platform string

const (
	ANDROID Platform = "android"
	iOS              = "iOS"
	WEB              = "web"
)

type Settings struct {
	AgeStart        uint8      `json:"ageStart" validate:"ltecsfield=AgeEnd,gte=1,lte=125"`
	AgeEnd          uint8      `json:"ageEnd" validate:"gtecsfield=AgeStart,gte=1,lte=125"`
	TargetGender    []Gender   `json:"gender"`
	TargetCountries []Country  `json:"countries"`
	TargetPlatforms []Platform `json:"platforms"`
}

func (s *Settings) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SponsoredPost struct {
	Title      string    `json:"title" validate:"required"`
	StartAt    time.Time `json:"startAt" validate:"required,ltecsfield=EndAt"`
	EndAt      time.Time `json:"endAt" validate:"required,gtecsfield=StartAt"`
	Conditions Settings  `json:"conditions,omitempty"`
}

func (sp *SponsoredPost) Validate() error {
	validate := validator.New()
	return validate.Struct(sp)
}

type QueryItems struct {
	Title string    `json:"title" validate:"required"`
	EndAt time.Time `json:"endAt" validate:"required"`
}

// func getSponsoredPost(c echo.Context) error {

// 	offset := c.QueryParam("offset")
// 	limit := c.QueryParam("limit")
// 	age := c.QueryParam("age")
// 	gender := c.QueryParam("gender")
// 	country := c.QueryParam("country")
// 	platform := c.QueryParam("platform")

// 	return c.String(http.StatusOK, "")
// }

// func createSponsoredPost(c echo.Context) error {
// 	return c.String(http.StatusOK, "")
// }

func main() {
	e := echo.New()
	e.GET("/api/v1/ad", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world.")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
