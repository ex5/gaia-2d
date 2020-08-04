package calendar

import (
	"fmt"
)

const modulo uint8 = 60
const dayModulo uint8 = 24
const monthModulo uint8 = 30

type Month uint8

const (
	// Spring
	Lightwake Month = iota
	Greencrest
	Blossomreach
	// Summer
	Solarcrest
	Æstival
	Amplenuts
	// Autumn
	Withercrown
	Crimsongrasp
	Stormreach
	// Winter
	Icewane
	Nightsoltice
	Whitereign
)

func (m Month) String() string {
	return [...]string{
		"Lightwake",
		"Greencrest",
		"Blossomreach",
		"Solarcrest",
		"Growrich",
		"Amplenuts",
		"Withercrown",
		"Crimsongrasp",
		"Stormreach",
		"Icewane",
		"Nightcrown",
		"Whitereign",
	}[m]
}

type Time struct {
	SecondsSinceBeginningOfTime uint64

	Second uint8
	Minute uint8
	Hour   uint8
	Day    uint8
	Month  Month
	Year   uint8
}

func (self *Time) AddSecond() {
	self.SecondsSinceBeginningOfTime++

	self.Second++
	if self.Second == modulo {
		self.Second = 0
		self.Minute++
	}
	if self.Minute == modulo {
		self.Minute = 0
		self.Hour++
	}
	if self.Hour == dayModulo {
		self.Hour = 0
		self.Day++
	}
	if self.Day == monthModulo {
		self.Day = 0
		self.Month++
	}
	if self.Month == Whitereign {
		self.Month = 0
		self.Year++
	}
}

func (self *Time) GetTextStatus() string {
	return fmt.Sprintf(
		"Year %d, day %d of %s\n%02d:%02d", self.Year, self.Day, self.Month,
		self.Hour, self.Minute,
	)
}
