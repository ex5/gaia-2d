package messages

import (
	//"log"
)

type ControlMessage struct {
	Action string
	Data string
}

const ControlMessageType string = "ControlMessage"

func (ControlMessage) Type() string {
  return ControlMessageType
}

type InteractionMessage struct {
	Action string
	BasicEntityID uint64
}

const InteractionMessageType string = "InteractionMessage"

func (InteractionMessage) Type() string {
	return InteractionMessageType
}
