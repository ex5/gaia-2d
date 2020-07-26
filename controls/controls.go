package controls

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/messages"
	"log"
)

type controlEntity struct {
	*ecs.BasicEntity
	*common.MouseComponent
}

type ControlsSystem struct {
	entities []controlEntity
	hoveredEntity *controlEntity
}

func (self *ControlsSystem) Add(basic *ecs.BasicEntity, mouse *common.MouseComponent) {
	self.entities = append(self.entities, controlEntity{basic, mouse})
}

func (self *ControlsSystem) Update(dt float32) {
	if engo.Input.Button("ExitToDesktop").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "exit",
		})
	}
	if engo.Input.Button("AddCreature").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "add_creature",
			Data: "textures/chick_32x32.png",
		})
	}
	if engo.Input.Button("AddObject").JustPressed() {
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "add_object",
			AtlasID: 1664,
		})
	}

	var newHoveredEntity *controlEntity
	for _, entity := range self.entities {
		if entity.MouseComponent.Hovered {
			newHoveredEntity = &entity
		}
	}
	if newHoveredEntity == nil && self.hoveredEntity != nil {
		log.Printf("Not hovering anything\n")
		engo.SetCursor(engo.CursorNone)
		engo.Mailbox.Dispatch(messages.InteractionMessage{
			Action: "mouse_hover",
			BasicEntity: nil,
		})
	} else if newHoveredEntity != nil && self.hoveredEntity == nil {
		log.Printf("Hovering over an entity: %+v #%d\n", newHoveredEntity, newHoveredEntity.ID())
		engo.Mailbox.Dispatch(messages.InteractionMessage{
			Action: "mouse_hover",
			BasicEntity: newHoveredEntity.BasicEntity,
		})
		engo.SetCursor(engo.CursorHand)
	}
	self.hoveredEntity = newHoveredEntity
}


func (self *ControlsSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range self.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}
