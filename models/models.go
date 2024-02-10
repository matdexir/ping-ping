package models

import (
	"time"

	"github.com/go-playground/validator/v10"
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
	IOS     Platform = "iOS"
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

type QueryItems struct {
	Title string    `json:"title" validate:"required"`
	EndAt time.Time `json:"endAt" validate:"required"`
}

type InsertedPost struct {
	ID   int64
	Post *SponsoredPost
}