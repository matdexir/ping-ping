package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/matdexir/ping-ping/db"

	"github.com/labstack/echo/v4"
)

type Gender = string

const (
	M Gender = "M"
	F Gender = "F"
	A Gender = "ALL"
)

type Country = string

const (
	JP  Country = "JP"
	TW  Country = "TW"
	US  Country = "US"
	BR  Country = "BR"
	SA  Country = "SA"
	FR  Country = "FR"
	ALL Country = "ALL"
)

type Platform = string

const (
	ANDROID Platform = "android"
	iOS     Platform = "iOS"
	WEB     Platform = "web"
)

type SponsoredPost struct {
	Title      string    `json:"title" validate:"required"`
	StartAt    time.Time `json:"startAt" validate:"required,ltecsfield=EndAt"`
	EndAt      time.Time `json:"endAt" validate:"required,gtecsfield=StartAt"`
	Conditions Settings  `json:"conditions,omitempty"`
}

type Settings struct {
	AgeStart       uint64   `json:"ageStart" validate:"ltecsfield=AgeEnd,gte=1,lte=125"`
	AgeEnd         uint64   `json:"ageEnd" validate:"gtecsfield=AgeStart,gte=1,lte=125"`
	TargetGender   Gender   `json:"gender"`
	TargetCountry  Country  `json:"countries"`
	TargetPlatform Platform `json:"platforms"`
}

func (s *Settings) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

func (sp *SponsoredPost) Validate() error {
	validate := validator.New()
	return validate.Struct(sp)
}

type InsertedPost struct {
	ID   int64
	Post *SponsoredPost
}

func CreateSponsoredPost(c echo.Context) error {
	var sp SponsoredPost
	if err := c.Bind(sp); err != nil {
		return err
	}

	db, _ := db.CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO posts VALUES(NULL, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Database.Exec(sqlStatement, &sp.Title, &sp.EndAt, &sp.Conditions.AgeStart, &sp.Conditions.AgeEnd, &sp.Conditions.TargetGender, &sp.Conditions.TargetCountry, &sp.Conditions.TargetPlatform)
	if err != nil {
		return c.String(http.StatusBadRequest, "Unable to insert into database")
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return c.String(http.StatusBadRequest, "Insert failed")
	}

	insertedPost := InsertedPost{
		ID: id, Post: &sp}

	return c.JSON(http.StatusOK, insertedPost)
}

type Post struct {
	title string
	endAt time.Time
}

func GetSponsoredPost(c echo.Context) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	age := c.QueryParam("age")
	country := c.QueryParam("country")
	platform := c.QueryParam("platform")

	if offset == "" || limit == "" {
		return c.String(http.StatusBadRequest, "Offset and/or Limit cannot be empty")
	}

	sqlStatement := `
      SELECT 
        id, title, endAt, ageStart, ageEnd, targetGender, targetCountries, targetPlatforms  
      FROM 
        posts
      ORDER BY 
        id 
      LIMIT 
        ?
      OFFSET
        ?`

	db, _ := db.CreateConnection()
	defer db.Close()

	row, err := db.Database.Query(sqlStatement, limit, offset)
	if err != nil {
		return c.String(http.StatusBadRequest, "Unable to query")
	}

	posts := []Post{}

	for row.Next() {
		var post SponsoredPost
		err := row.Scan(&post.Title, &post.EndAt, &post.Conditions.AgeStart, &post.Conditions.AgeEnd, &post.Conditions.TargetGender, &post.Conditions.TargetCountry, &post.Conditions.TargetPlatform)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to scan")
		}

		if len(country) > 0 && post.Conditions.TargetCountry != country {
			continue
		}

		if len(platform) > 0 && post.Conditions.TargetPlatform != platform {
			continue
		}

		if len(age) > 0 {
			age, err := strconv.ParseUint(age, 10, 8)
			if err != nil {
				log.Fatalf("Unable to convert age")
				return c.String(http.StatusBadRequest, "Age is not a proper integer")
			}
			if post.Conditions.AgeStart > age || post.Conditions.AgeEnd < age {
				continue
			}
		}

		tmp := Post{title: post.Title, endAt: post.EndAt}

		posts = append(posts, tmp)
	}

	return c.JSON(http.StatusOK, posts)

}
