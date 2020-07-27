package messages

import (
	"github.com/EngoEngine/ecs"
	//"github.com/EngoEngine/engo/common"
	//"log"
)

type ControlMessage struct {
	Action string
	Data string
	SpriteID int
}

const ControlMessageType string = "ControlMessage"

func (ControlMessage) Type() string {
	return ControlMessageType
}

type InteractionMessage struct {
	Action string
	BasicEntity *ecs.BasicEntity
}

const InteractionMessageType string = "InteractionMessage"

func (InteractionMessage) Type() string {
	return InteractionMessageType
}

type SaveMessage struct {
	Filepath string
}

const SaveMessageType string = "SaveMessage"

func (SaveMessage) Type() string {
	return SaveMessageType
}

// HUDTextMessage updates the HUD text based on messages sent from other systems
type HUDTextMessage struct {
	//ecs.BasicEntity
	//common.SpaceComponent
	//common.MouseComponent
	Placeholder string
	Line1, Line2, Line3, Line4 string
}

// HUDTextMessageType is the type for an HUDTextMessage
const HUDTextMessageType string = "HUDTextMessage"

// Type implements the engo.Message Interface
func (HUDTextMessage) Type() string {
	return HUDTextMessageType
}
