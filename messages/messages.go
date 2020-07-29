package messages

import (
	"github.com/EngoEngine/ecs"
	//"github.com/EngoEngine/engo/common"
	//"log"
)

const ControlMessageType string = "ControlMessage"
const InteractionMessageType string = "InteractionMessage"
const SaveMessageType string = "SaveMessage"
const LoadMessageType string = "LoadMessage"

type ControlMessage struct {
	Action        string
	Data          string
	ObjectID      int
	CreatureID    int
}

type InteractionMessage struct {
	Action      string
	BasicEntity *ecs.BasicEntity
}

type SaveMessage struct {
	Filepath string
}

type LoadMessage struct {
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

func (LoadMessage) Type() string {
	return LoadMessageType
}
