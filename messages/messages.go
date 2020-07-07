package messages

import (
	//"log"
)

type ControlMessage struct {
	Action string
}

const ControlMessageType string = "ControlMessage"

func (ControlMessage) Type() string {
  return ControlMessageType
}
