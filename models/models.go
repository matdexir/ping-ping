package models

import (
	"fmt"
	"reflect"
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

// type SponsoredPostGuard struct {
// 	Title      string   `json:"title" validate:"required"`
// 	StartAt    string   `json:"startAt" validate:"required"`
// 	EndAt      string   `json:"endAt" validate:"required"`
// 	Conditions Settings `json:"conditions,omitempty"`
// }

func Serialize(q interface{}) {
	t := reflect.TypeOf(q)
	v := reflect.ValueOf(q)

	fmt.Println("Type:", t)
	fmt.Println("Value:", v)
}

type SponsoredPost struct {
	Title      string    `json:"title" validate:"required"`
	StartAt    time.Time `json:"startAt" validate:"required"`
	EndAt      time.Time `json:"endAt" validate:"required"`
	Conditions Settings  `json:"conditions,omitempty"`
}

type Settings struct {
	AgeStart       uint64   `json:"ageStart" validate:"gte=1,lte=125,ltecsfield=AgeEnd"`
	AgeEnd         uint64   `json:"ageEnd" validate:"gte=1,lte=125,gtecsfield=AgeStart"`
	TargetGender   Gender   `json:"gender"`
	TargetCountry  Country  `json:"country"`
	TargetPlatform Platform `json:"platform"`
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
