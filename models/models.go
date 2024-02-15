package models

import (
	"encoding/json"
	"fmt"
	// "log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type Gender int

const (
	MALE Gender = iota
	FEMALE
)

func (g Gender) String() string {
	return [...]string{"M", "F"}[g]
}

func (g Gender) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *Gender) UnmarshalJSON(data []byte) error {
	var genderStr string
	if err := json.Unmarshal(data, &genderStr); err != nil {
		return err
	}
	parsed, ok := Parse[Gender](genderStr, GenderHint)
	if !ok {
		return fmt.Errorf("Unable to unmarshal Gender JSON")
	}
	*g = parsed.(Gender)
	return nil
}

type Country int

const (
	Japan Country = iota
	Taiwan
	USA
	Brazil
	SouthAfrica
	France
)

func (c Country) String() string {
	return [...]string{"JP", "TW", "US", "BR", "SA", "FR"}[c]
}

func (c Country) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Country) UnmarshalJSON(data []byte) error {
	var countryStr string
	if err := json.Unmarshal(data, &countryStr); err != nil {
		return err
	}
	parsed, ok := Parse[Country](countryStr, CountryHint)
	if !ok {
		return fmt.Errorf("Unable to unmarshal Country JSON")
	}
	*c = parsed.(Country)
	return nil
}

type Platform int

const (
	ANDROID Platform = iota
	IOS
	WEB
)

func (p Platform) String() string {
	return [...]string{"android", "iOS", "web"}[p]
}

func (p Platform) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p *Platform) UnmarshalJSON(data []byte) error {
	var platformStr string

	if err := json.Unmarshal(data, &platformStr); err != nil {
		return err
	}

	parsed, ok := Parse[Platform](platformStr, PlatformHint)

	if !ok {
		return fmt.Errorf("Unable to parse the Platform json.")
	}

	*p = parsed.(Platform)

	return nil
}

type SerializableItem interface {
	// Platform | Gender | Country
	String() string
}

type Hint int

const (
	PlatformHint Hint = iota
	CountryHint
	GenderHint
)

func getMap(hint Hint) map[string]SerializableItem {
	switch hint {
	case PlatformHint:
		return map[string]SerializableItem{
			"android": ANDROID,
			"ios":     IOS,
			"web":     WEB,
		}
	case CountryHint:
		return map[string]SerializableItem{
			"jp": Japan,
			"tw": Taiwan,
			"us": USA,
			"br": Brazil,
			"sa": SouthAfrica,
			"fr": France,
		}
	case GenderHint:
		return map[string]SerializableItem{
			"m": MALE,
			"f": FEMALE,
		}
	default:
		panic("invalid hint")
	}
}

// should I use generic?
// SerializableItem needs to be a concrete type
// however I am not too sure how to make it fit the requirement of both String()
func Parse[T SerializableItem](item string, hint Hint) (interface{}, bool) {
	mp := getMap(hint)
	value, ok := mp[strings.ToLower(item)]
	// log.Printf("%v: %v\n", item, mp)
	if !ok {
		return nil, false
	}
	return value.(T), true
}

func Serialize[T SerializableItem](items []T) string {
	serialized := ""
	for idx, item := range items {
		serialized += item.String()
		if idx != len(items)-1 {
			serialized += ", " // Add space after comma for readability
		}
	}
	return serialized
}

func Deserialize[T SerializableItem](dbField string, hint Hint) ([]T, error) {
	enums := []T{}
	splitField := strings.Split(dbField, ", ")

	for _, split := range splitField {
		parsed, ok := Parse[T](split, hint)
		if !ok {
			return nil, fmt.Errorf("unable to deserialize field: %s", split)
		}
		enums = append(enums, parsed.(T))
	}
	return enums, nil
}

type SponsoredPost struct {
	Title      string    `json:"title" validate:"required"`
	StartAt    time.Time `json:"startAt" validate:"required,ltecsfield=EndAt"`
	EndAt      time.Time `json:"endAt" validate:"required,gtecsfield=StartAt"`
	Conditions Settings  `json:"conditions,omitempty" validate:"omitempty"`
}

type Settings struct {
	AgeStart       uint64     `json:"ageStart" validate:"gte=1,lte=125,ltecsfield=AgeEnd"`
	AgeEnd         uint64     `json:"ageEnd" validate:"gte=1,lte=125,gtecsfield=AgeStart"`
	TargetGender   []Gender   `json:"gender"`
	TargetCountry  []Country  `json:"country"`
	TargetPlatform []Platform `json:"platform"`
}

func (s *Settings) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(s)
}

func (sp *SponsoredPost) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
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
