package messages

import (
	"github.com/EngoEngine/engo"
)

const SpacialRequestMessageType string = "SpacialRequestMessage"
const SpacialResponseMessageType string = "SpacialResponseMessage"

type SpacialRequestMessage struct {
	EntityID uint64
	EventID  uint64
	Aabb     engo.AABB
	Filter   func(engo.AABBer) bool
}

type SpacialResponseMessage struct {
	EntityID      uint64
	EventID       uint64
	BasicEntityID int
	Aabb          engo.AABB
	Filter        func(engo.AABBer) bool
	Result        []engo.AABBer
}

func (SpacialRequestMessage) Type() string {
	return SpacialRequestMessageType
}

func (SpacialResponseMessage) Type() string {
	return SpacialResponseMessageType
}
