package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"gogame/calendar"
	"gogame/messages"
)

type TimeSystem struct {
	world         *ecs.World
	dtFullSeconds float32

	Speed uint8
	Time  *calendar.Time
	Month calendar.Month
}

func (self *TimeSystem) New(w *ecs.World) {
	self.world = w
	self.Time = &calendar.Time{}
	engo.Mailbox.Dispatch(messages.TimeSecondPassedMessage{
		Time: self.Time,
	})
}

func (*TimeSystem) Add() {}

func (self *TimeSystem) Update(dt float32) {
	if self.dtFullSeconds > 1 {
		self.dtFullSeconds = 0
		self.Time.AddSecond()
		engo.Mailbox.Dispatch(messages.TimeSecondPassedMessage{
			Time: self.Time,
			Dt:   self.dtFullSeconds,
		})
	}
	self.dtFullSeconds += dt
	// TODO might be the good place to implement speed of in-game time
	// TODO and the PauseSystem
}

func (*TimeSystem) Remove(ecs.BasicEntity) {}
