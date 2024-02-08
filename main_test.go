package main

import (
	"encoding/json"
	"testing"
	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

// A normal post is a post with no missing fields, i.e, all values are defined
func TestMarshallNormalPost(t *testing.T) {

	startTime, _ := time.Parse(time.RFC3339, "2024-02-07T14:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-02-08T14:00:00Z")

	normal_post := SponsoredPost{
		Title:   "post1",
		StartAt: startTime,
		EndAt:   endTime,
		Conditions: Settings{
			AgeStart:        20,
			AgeEnd:          35,
			TargetGender:    []Gender{M, F},
			TargetCountries: []Country{JP, FR},
			TargetPlatforms: []Platform{iOS, ANDROID},
		},
	}

	data, err := json.Marshal(normal_post)
	t.Log(string(data))
	require.NoError(t, err)

}

func TestUnmarshallInvalidPost(t *testing.T) {
	// startTime, _ := time.Parse(time.RFC3339, "2024-02-07T14:00:00Z")
	data := `
  { 
    "startAt":"2024-02-07T14:00:00Z",
    "endAt":"2024-02-08T14:00:00Z",
    "conditions": {
      "ageStart":20,
      "ageEnd":35,
      "gender":["M","F"],
      "countries":["JP","FR"],
      "platforms":["iOS","android"]
    }
  }`
	var post SponsoredPost
	err := json.Unmarshal([]byte(data), &post)
	t.Log(post)
	require.NoError(t, err)

	err = post.Validate()
	// err = validator.New().Struct(post)
	require.Error(t, err)

}

func TestUnmarshallNormalPost(t *testing.T) {
	data := `
  { 
    "title":"post1",
    "startAt":"2024-02-07T14:00:00Z",
    "endAt":"2024-02-08T14:00:00Z",
    "conditions": {
      "ageStart":20,
      "ageEnd":35,
      "gender":["M","F"],
      "countries":["JP","FR"],
      "platforms":["iOS","android"]
    }
  }`

	expectedStart, _ := time.Parse(time.RFC3339, "2024-02-07T14:00:00Z")
	expectedEnd, _ := time.Parse(time.RFC3339, "2024-02-08T14:00:00Z")

	var post SponsoredPost
	json.Unmarshal([]byte(data), &post)
	require.Equal(t, post.Title, "post1")
	require.Equal(t, post.StartAt, expectedStart)
	require.Equal(t, post.EndAt, expectedEnd)
	require.Equal(t, post.Conditions.AgeStart, uint8(20))
	require.Equal(t, post.Conditions.AgeEnd, uint8(35))
	require.Equal(t, post.Conditions.TargetGender, []Gender{M, F})
	require.Equal(t, post.Conditions.TargetCountries, []Country{JP, FR})
	require.Equal(t, post.Conditions.TargetPlatforms, []Platform{iOS, ANDROID})
	t.Log("Done")
}
