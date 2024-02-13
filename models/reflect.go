package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Gender int

const (
	Male Gender = iota
	Female
)

func (g Gender) String() string {
	return [...]string{"Male", "Female"}[g]
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
	return [...]string{"Japan", "Taiwan", "USA", "Brazil", "SouthAfrica", "France"}[c]
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
			"iOS":     IOS,
			"web":     WEB,
		}
	case CountryHint:
		return map[string]SerializableItem{
			"japan":       Japan,
			"taiwan":      Taiwan,
			"usa":         USA,
			"brazil":      Brazil,
			"southafrica": SouthAfrica,
			"france":      France,
		}
	case GenderHint:
		return map[string]SerializableItem{
			"m": Male,
			"f": Female,
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
	splitField := strings.Split(dbField, ",")

	for _, split := range splitField {
		parsed, ok := Parse[T](split, hint)
		if !ok {
			return nil, fmt.Errorf("unable to deserialize field: %s", split)
		}
		enums = append(enums, parsed.(T))
	}
	return enums, nil
}

func main() {
	platforms := []Platform{ANDROID, IOS, WEB}
	serializedPlatforms := Serialize(platforms)
	marshaledPlatforms, _ := json.Marshal(platforms)
	fmt.Println("Serialized platforms:", serializedPlatforms) // Output: android, iOS, web
	fmt.Println("Marshaled platforms:", string(marshaledPlatforms))

	genders := []Gender{Male, Female}
	serializedGenders := Serialize(genders)
	fmt.Println("Serialized genders:", serializedGenders) // Output: Male, Female

	gender, ok := Parse[Gender]("M", GenderHint)
	if ok {
		fmt.Println("Parsed gender:", gender) // Output: Male
	} else {
		fmt.Println("Error parsing gender")
	}

	country, ok := Parse[Country]("Japan", CountryHint)
	if ok {
		fmt.Println("Parsed country:", country) // Output: Japan
	} else {
		fmt.Println("Error parsing gender")
	}
}
