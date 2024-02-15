package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	// "strconv"

	"github.com/matdexir/ping-ping/db"
	"github.com/matdexir/ping-ping/memcached"
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

	sqlRawString := `INSERT INTO posts VALUES(NULL, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.Database.Prepare(sqlRawString)

	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	res, err := stmt.Exec(&sp.Title, &sp.StartAt, &sp.EndAt, &sp.Conditions.AgeStart, &sp.Conditions.AgeEnd, &targetGender, &targetCountry, &targetPlatform)

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
	params := c.QueryParams()

	if offset == "" || limit == "" {
		return c.String(http.StatusBadRequest, "Offset and/or Limit cannot be empty")
	}

	mc, _ := memcached.NewMemcached()
	defer mc.Close()
	items, err := mc.GetPosts(params.Encode())
	if err == nil {
		return c.JSON(http.StatusOK, items.Items)
	}

	sqlRawString := `
      SELECT 
        id, title, endAt, ageStart, ageEnd, targetGender, targetCountry, targetPlatform  
      FROM 
        posts`

	var whereClause string
	var args []interface{}
	addAnd := false

	if age != "" {
		whereClause += `ageStart <= ? AND ageEnd >= ? `
		addAnd = true
		args = append(args, age, age)
	}

	if country != "" {
		if addAnd {
			whereClause += `AND `
			// addAnd = false
		}
		whereClause += `targetCountry LIKE '%'||?||'%' `
		addAnd = true
		args = append(args, country)
	}

	if platform != "" {
		if addAnd {
			whereClause += `AND `
			// addAnd = false
		}
		whereClause += `targetPlatform LIKE '%'||?||'%' `
		addAnd = true
		args = append(args, platform)
	}

	if gender != "" {
		if addAnd {
			whereClause += `AND `
			// addAnd = false
		}
		whereClause += `targetGender LIKE '%'||?||'%' `
		// addAnd = true
		args = append(args, gender)
	}

	if whereClause != "" {
		sqlRawString += ` WHERE ` + whereClause
	}

	orderStatement := `
      ORDER BY 
        id 
      LIMIT 
        ?
      OFFSET
        ?`

	sqlRawString += orderStatement
	args = append(args, limit, offset)

	fmt.Println(sqlRawString)
	db, _ := db.CreateConnection()
	defer db.Close()

	stmt, err := db.Database.Prepare(sqlRawString)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}

	row, err := stmt.Query(args...)

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

		tmp := models.QueryItems{Title: post.Title, EndAt: post.EndAt}
		log.Printf("%v\n", tmp)

		posts = append(posts, tmp)
	}

	_ = mc.SetPosts(models.QueryCache{Parameters: params.Encode(), Items: posts})

	return c.JSON(http.StatusOK, posts)

}
