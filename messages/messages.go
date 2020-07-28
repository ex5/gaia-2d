package messages

import (
	"github.com/EngoEngine/ecs"
	//"github.com/EngoEngine/engo/common"
	//"log"
)

const ControlMessageType string = "ControlMessage"
const InteractionMessageType string = "InteractionMessage"
const SaveMessageType string = "SaveMessage"

type ControlMessage struct {
	Action   string
	Data     string
	SpriteID int
}

type InteractionMessage struct {
	Action      string
	BasicEntity *ecs.BasicEntity
}

type SaveMessage struct {
	Filepath string
}

func (ControlMessage) Type() string {
	return ControlMessageType
}

func (InteractionMessage) Type() string {
	return InteractionMessageType
}

func (SaveMessage) Type() string {
	return SaveMessageType
}
