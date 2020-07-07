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
