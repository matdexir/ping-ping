package main

import (
	"database/sql"
	"fmt"
	// "net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/matdexir/ping-ping/controllers"
	"github.com/matdexir/ping-ping/db"
)

type Gender string

const (
	M Gender = "M"
	F Gender = "F"
)

type Country string

const (
	JP Country = "JP"
	TW Country = "TW"
	US Country = "US"
	BR Country = "BR"
	SA Country = "SA"
	FR Country = "FR"
)

type Platform string

const (
	ANDROID Platform = "android"
	iOS     Platform = "iOS"
	WEB     Platform = "web"
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

func newDB() (*db.PostDB, error) {
	sql, err := sql.Open("sqlite3", "./db/file.db")
	if err != nil {
		fmt.Println("Unable to open database")
		return nil, err
	}
	return &db.PostDB{Database: sql}, nil
}

func main() {
	e := echo.New()

	database, err := newDB()
	if err != nil {
		os.Exit(1)
	}
	err = database.CreateTable()
	if err != nil {
		os.Exit(1)
	}

	e.POST("/api/v1/ad", controllers.CreateSponsoredPost)
	e.GET("/api/v1/ad", controllers.GetSponsoredPost)
	e.Logger.Fatal(e.Start(":8080"))
}
