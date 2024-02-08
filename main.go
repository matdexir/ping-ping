package main

import (
	// "errors"
	// "fmt"
	"time"
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
	AgeStart        uint8      `json:"ageStart"`
	AgeEnd          uint8      `json:"ageEnd"`
	TargetGender    []Gender   `json:"gender"`
	TargetCountries []Country  `json:"countries"`
	TargetPlatforms []Platform `json:"platforms"`
}

type SponsoredPost struct {
	Title      string    `json:"title"`
	StartAt    time.Time `json:"startAt"`
	EndAt      time.Time `json:"endAt"`
	Conditions Settings  `json:"conditions"`
}

func main() {
	println("hello world")
}
