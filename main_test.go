package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/matdexir/ping-ping/models"
	"github.com/stretchr/testify/require"
)

// A normal post is a post with no missing fields, i.e, all values are defined
func TestMarshallNormalPost(t *testing.T) {

	startTime, _ := time.Parse(time.RFC3339, "2024-02-07T14:00:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-02-08T14:00:00Z")

	normal_post := models.SponsoredPost{
		Title:   "post1",
		StartAt: startTime,
		EndAt:   endTime,
		Conditions: models.Settings{
			AgeStart:       20,
			AgeEnd:         35,
			TargetGender:   []models.Gender{models.MALE, models.FEMALE},
			TargetCountry:  []models.Country{models.Brazil, models.Japan},
			TargetPlatform: []models.Platform{models.ANDROID, models.WEB},
		},
	}

	_, err := json.Marshal(normal_post)
	// t.Log(string(data))
	require.NoError(t, err)

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
      "gender":["M"],
      "country":["JP"],
      "platform":["iOS"]
    }
  }`

	expectedStart, _ := time.Parse(time.RFC3339, "2024-02-07T14:00:00Z")
	expectedEnd, _ := time.Parse(time.RFC3339, "2024-02-08T14:00:00Z")

	var post models.SponsoredPost
	err := json.Unmarshal([]byte(data), &post)
	log.Printf("%v", post)
	require.NoError(t, err)
	require.Equal(t, post.Title, "post1")
	require.Equal(t, post.StartAt, expectedStart)
	require.Equal(t, post.EndAt, expectedEnd)
	require.Equal(t, post.Conditions.AgeStart, uint64(20))
	require.Equal(t, post.Conditions.AgeEnd, uint64(35))
	require.Equal(t, post.Conditions.TargetGender, []models.Gender{models.MALE})
	require.Equal(t, post.Conditions.TargetCountry, []models.Country{models.Japan})
	require.Equal(t, post.Conditions.TargetPlatform, []models.Platform{models.IOS})
	t.Log("Done")
}

func TestUnmarshallMissingFieldPost(t *testing.T) {
	// field title is missing
	data := `
  { 
    "startAt":"2024-02-07T14:00:00Z",
    "endAt":"2024-02-08T14:00:00Z",
    "conditions": {
      "ageStart":20,
      "ageEnd":35,
      "gender":["M"],
      "countries":["JP"],
      "platforms":["iOS"]
    }
  }`
	var post models.SponsoredPost
	err := json.Unmarshal([]byte(data), &post)
	// t.Log(post)
	require.NoError(t, err)

	err = post.Validate()
	require.Error(t, err)

}

func TestUnmarshallAgeOutOfBounds(t *testing.T) {
	t.Parallel()
	testData := []string{`
    { 
      "title":"post1",
      "startAt":"2024-02-07T14:00:00Z",
      "endAt":"2024-02-08T14:00:00Z",
      "conditions": {
        "ageStart":0,
        "ageEnd":35,
        "gender":["M"],
        "country":["JP"],
        "platform":["iOS"]
      }
    }`, `
    { 
      "title":"post1",
      "startAt":"2024-02-07T14:00:00Z",
      "endAt":"2024-02-08T14:00:00Z",
      "conditions": {
        "ageStart":1,
        "ageEnd":126,
        "gender":["M"],
        "country":["JP"],
        "platform":["iOS"]
      }
    }`, `
    { 
      "title":"post1",
      "startAt":"2024-02-07T14:00:00Z",
      "endAt":"2024-02-08T14:00:00Z",
      "conditions": {
        "ageStart":45,
        "ageEnd":23,
        "gender":["M"],
        "country":["JP"],
        "platform":["iOS"]
      }
    }`,
	}

	var post *models.SponsoredPost
	for idx, td := range testData {
		t.Run(fmt.Sprintf("%v", idx), func(t *testing.T) {
			err := json.Unmarshal([]byte(td), &post)
			// t.Log(post)
			require.NoError(t, err)
			err = post.Validate()
			// t.Log(err)
			require.Error(t, err)

		})
	}

}
