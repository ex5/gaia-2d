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
}

func (self *ControlsSystem) Add(basic *ecs.BasicEntity, mouse *common.MouseComponent) {
	self.entities = append(self.entities, controlEntity{basic, mouse})
}

func (self *ControlsSystem) Update(dt float32) {
	if engo.Input.Button("ExitToDesktop").JustPressed() {
		log.Println("ExitToDesktop pressed")
		engo.Mailbox.Dispatch(messages.ControlMessage{
			Action: "exit",
		})
	}

	//log.Printf("Entities: %+v", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		if entity.MouseComponent.Enter {
			engo.SetCursor(engo.CursorHand)
		} else if entity.MouseComponent.Leave {
			engo.SetCursor(engo.CursorNone)
		}
	}
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
