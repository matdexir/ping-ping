package main

import (
	"errors"
	"fmt"
	"time"
)

type Gender int

const (
	Male Gender = iota
	Female
)

type Country int

const (
	JP Country = iota
	TW
	US
	BR
	SA
	FR
)

type Platform int

const (
	android Platform = iota
	iOS
	web
)

type Conditions struct {
	age      uint8
	gender   Gender
	country  Country
	platform Platform
}

func (c *Conditions) SetAge(age uint8) error {
	if age < 1 || age > 100 {
		return errors.New("Value must be between 1 and 100")
	}
	c.age = age
	return nil
}

func (c *Conditions) SetGender(gender Gender) error {
	c.gender = gender
	return nil
}

func (c *Conditions) SetCountry(country Country) error {
	c.country = country
	return nil
}

func (c *Conditions) SetPlatform(platform Platform) error {
	c.platform = platform
	return nil
}

type SponsoredPost struct {
	title      string
	startAt    time.Time
	endAt      time.Time
	conditions Conditions
}

func main() {
	fmt.Println("Hello world.")
}
