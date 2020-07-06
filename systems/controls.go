package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"log"
)

type ControlsSystem struct {
}

func (*ControlsSystem) Update(dt float32) {
	if engo.Input.Button("ExitToDesktop").JustPressed() {
		log.Println("ExitToDesktop")
		engo.Exit()
	}
}

func (*ControlsSystem) Remove(ecs.BasicEntity) {}
