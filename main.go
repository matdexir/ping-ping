package main

import (
	// "errors"
	// "fmt"
	"github.com/go-playground/validator/v10"
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

func main() {
	println("hello world")
}
