package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/matdexir/ping-ping/db"
	"github.com/matdexir/ping-ping/models"

	"github.com/labstack/echo/v4"
)

func CreateSponsoredPost(c echo.Context) error {
	var sp models.SponsoredPost
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

	insertedPost := models.InsertedPost{
		ID: id, Post: &sp}

	return c.JSON(http.StatusOK, insertedPost)
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

	posts := []models.QueryItems{}

	for row.Next() {
		var post models.SponsoredPost
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

		tmp := models.QueryItems{Title: post.Title, EndAt: post.EndAt}

		posts = append(posts, tmp)
	}

	return c.JSON(http.StatusOK, posts)

}
