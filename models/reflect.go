package main

import (
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

type Platform int

const (
	ANDROID Platform = iota
	IOS
	WEB
)

func (p Platform) String() string {
	return [...]string{"android", "iOS", "web"}[p]
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

func Parse[T SerializableItem](item string, hint Hint) (T, bool) {
	mp := getMap(hint)
	value, ok := mp[strings.ToLower(item)]
	if !ok {
		return T{}, false
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
		enums = append(enums, parsed)
	}
	return enums, nil
}

func main() {
	platforms := []Platform{ANDROID, IOS, WEB}
	serializedPlatforms := Serialize(platforms)
	fmt.Println("Serialized platforms:", serializedPlatforms) // Output: android, iOS, web

	genders := []Gender{Male, Female}
	serializedGenders := Serialize(genders)
	fmt.Println("Serialized genders:", serializedGenders) // Output: Male, Female

	gender, ok := Parse[Gender]("M", GenderHint)
	if ok {
		fmt.Println("Parsed gender:", gender) // Output: Male
	} else {
		fmt.Println("Error parsing gender")
	}
}
