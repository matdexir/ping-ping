package main

import (
	"fmt"
	// "reflect"
)

type Gender int

const (
	M Gender = iota
	F
)

func (g Gender) String() string {
	return [...]string{"M", "F"}[g]
}

type Country int

const (
	JP Country = iota
	TW
	US
	BR
	SA
	FR
)

func (c Country) String() string {
	return [...]string{"JP", "TW", "US", "BR", "SA", "FR"}[c]
}

type Platform int

const (
	ANDROID Platform = iota
	IOS
	WEB
)

type Hint int

const (
	P Hint = iota
	C
	G
)

func (p Platform) String() string {
	return [...]string{"android", "iOS", "web"}[p]
}

type Serializable interface {
	Platform | Gender | Country
}

func Serialize(items []Serializable) string {
	serialized := ""
	for idx, item := range items {
		serialized += item.String()

		if idx != len(items)-1 {
			serialized += ","
		}
	}
	return serialized
}

func Deserialize(dbField string, hint Hint) []Serializable {
	split, err := dbField.Split(',')
	if err != nil {
		return []Serializable
	}
	return split
}

func main() {

	ans := Serialize([]Platform{ANDROID, IOS, WEB})
	fmt.Println(ans)
	ans = Serialize([]Gender{M, F})
	fmt.Println(ans)

	fmt.Println()
}
