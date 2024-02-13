package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/matdexir/ping-ping/db"
	"github.com/matdexir/ping-ping/models"

	"github.com/labstack/echo/v4"
)

func GetRAWJSON(body io.ReadCloser) map[string]interface{} {
	jsonBody := make(map[string]interface{})
	_ = json.NewDecoder(body).Decode(&jsonBody)

	return jsonBody
}

func CreateSponsoredPost(c echo.Context) error {
	sp := new(models.SponsoredPost)
	if err := c.Bind(&sp); err != nil {
		log.Printf("RAW: %v\n", GetRAWJSON(c.Request().Body))
		return c.String(http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	log.Printf("Post is: %+v\n", sp)

	db, _ := db.CreateConnection()
	defer db.Close()

	targetCountry := models.Serialize[models.Country](sp.Conditions.TargetCountry)
	targetPlatform := models.Serialize[models.Platform](sp.Conditions.TargetPlatform)
	targetGender := models.Serialize[models.Gender](sp.Conditions.TargetGender)

	sqlStatement := `INSERT INTO posts VALUES(NULL, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.Database.Exec(sqlStatement, &sp.Title, &sp.StartAt, &sp.EndAt, &sp.Conditions.AgeStart, &sp.Conditions.AgeEnd, &targetGender, &targetCountry, &targetPlatform)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return c.String(http.StatusBadRequest, "Insert failed")
	}

	insertedPost := models.InsertedPost{
		ID: id, Post: sp}

	return c.JSON(http.StatusOK, insertedPost)
}

func GetSponsoredPost(c echo.Context) error {
	offset := c.QueryParam("offset")
	limit := c.QueryParam("limit")
	age := c.QueryParam("age")
	country := c.QueryParam("country")
	platform := c.QueryParam("platform")
	gender := c.QueryParam("gender")

	if offset == "" || limit == "" {
		return c.String(http.StatusBadRequest, "Offset and/or Limit cannot be empty")
	}

	sqlStatement := `
      SELECT 
        id, title, endAt, ageStart, ageEnd, targetGender, targetCountry, targetPlatform  
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
		return c.String(http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	posts := []models.QueryItems{}

	for row.Next() {
		post := new(models.SponsoredPost)
		var id uint64
		var endTime string
		var targetCountry string
		var targetPlatform string
		var targetGender string

		err := row.Scan(&id, &post.Title, &endTime, &post.Conditions.AgeStart, &post.Conditions.AgeEnd, &targetGender, &targetCountry, &targetPlatform)

		if err != nil {
			log.Printf("Unable to row scan: %v\n", err)
			continue
		}

		post.EndAt, err = time.Parse(time.DateTime, endTime[:len(endTime)-6])
		if err != nil {
			log.Printf("Unable to parse time: %v\n", err)
			continue
		}

		post.Conditions.TargetCountry, err = models.Deserialize[models.Country](targetCountry, models.CountryHint)
		if err != nil {
			log.Printf("Unable to deserialize country: %v\n", targetCountry)
			continue
		}

		post.Conditions.TargetGender, err = models.Deserialize[models.Gender](targetGender, models.GenderHint)
		if err != nil {
			log.Printf("Unable to deserialize gender: %v\n", targetGender)
			continue
		}

		post.Conditions.TargetPlatform, err = models.Deserialize[models.Platform](targetPlatform, models.PlatformHint)
		if err != nil {
			log.Printf("Unable to deserialize gender: %v\n", targetPlatform)
			continue
		}

		// ok := slices.Contains(post.Conditions.TargetCountry, country)
		ok := false
		for _, c := range post.Conditions.TargetCountry {
			if c.String() == country {
				ok = true
				break
			}
		}
		if len(country) > 0 && !ok {
			continue
		}

		ok = false
		for _, c := range post.Conditions.TargetPlatform {
			if c.String() == platform {
				ok = true
				break
			}
		}

		if len(platform) > 0 && !ok {
			continue
		}

		ok = false
		for _, c := range post.Conditions.TargetGender {
			if c.String() == gender {
				ok = true
				break
			}
		}

		if len(gender) > 0 && !ok {
			continue
		}

		if len(age) > 0 {
			age, err := strconv.ParseUint(age, 10, 8)
			if err != nil {
				log.Printf("Unable to convert age: %v\n", err)
				continue
			}
			if post.Conditions.AgeStart > age || post.Conditions.AgeEnd < age {
				continue
			}
		}

		tmp := models.QueryItems{Title: post.Title, EndAt: post.EndAt}
		log.Printf("%v\n", tmp)

		posts = append(posts, tmp)
	}

	return c.JSON(http.StatusOK, posts)

}
