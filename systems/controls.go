package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"log"
)

type ControlMessage struct {
	Action string
}

const ControlMessageType string = "ControlMessage"

func (ControlMessage) Type() string {
  return ControlMessageType
}

type ControlsSystem struct {
}

func (*ControlsSystem) Update(dt float32) {
	if engo.Input.Button("ExitToDesktop").JustPressed() {
		log.Println("ExitToDesktop pressed")
		engo.Mailbox.Dispatch(ControlMessage{
			Action: "exit",
		})
	}
}

func (*ControlsSystem) Remove(ecs.BasicEntity) {}
