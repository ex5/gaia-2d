package messages

import (
	"gogame/calendar"
)

const TimeSecondPassedMessageType = "TimeSecondPassedMessage"
const TimeSunriseMessageType = "TimeSunriseMessage"
const TimeSunsetMessageType = "TimeSunsetMessage"

type TimeSecondPassedMessage struct {
	Time *calendar.Time
	Dt   float32
}

type TimeSunriseMessage struct {
	Time *calendar.Time
	Dt   float32
}

type TimeSunsetMessage struct {
	Time *calendar.Time
	Dt   float32
}

func (TimeSecondPassedMessage) Type() string {
	return TimeSecondPassedMessageType
}

func (TimeSunriseMessage) Type() string {
	return TimeSunriseMessageType
}

func (TimeSunsetMessage) Type() string {
	return TimeSunsetMessageType
}
