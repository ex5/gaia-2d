package messages

import (
	"github.com/EngoEngine/engo"
	"time"
)

const SpacialRequestMessageType string = "SpacialRequestMessage"
const SpacialResponseMessageType string = "SpacialResponseMessage"
const DisplayDebugAABBMessageType string = "DisplayDebugAABBMessage"

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

type DisplayDebugAABBMessage struct {
	Aabbs       []engo.AABB
	Aabbers     []engo.AABBer
	RemoveAfter time.Duration
	Color       string
}

func (SpacialRequestMessage) Type() string {
	return SpacialRequestMessageType
}

func (SpacialResponseMessage) Type() string {
	return SpacialResponseMessageType
}

func (DisplayDebugAABBMessage) Type() string {
	return DisplayDebugAABBMessageType
}
