package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/messages"
	"log"
	"time"
)

type DebugShape struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent

	removeAfter time.Duration
	shownSince  time.Time
}

type DebugSystem struct {
	world    *ecs.World
	entities []*DebugShape
}

func (self *DebugSystem) New(w *ecs.World) {
	self.world = w

	// TODO listen for events
	engo.Mailbox.Listen(messages.DisplayDebugAABBMessageType, self.HandleDisplayDebugAABBMessage)
}

func (self *DebugSystem) AddAABBer(aabber engo.AABBer, colorName string, removeAfter time.Duration) {
	self.AddAABB(aabber.AABB(), colorName, removeAfter)
}

func (self *DebugSystem) AddAABB(aabb engo.AABB, colorName string, removeAfter time.Duration) {
	entity := &DebugShape{
		BasicEntity: ecs.NewBasic(),
		removeAfter: removeAfter,
		shownSince: time.Now(),
	}
	entity.SpaceComponent = &common.SpaceComponent{
		Position: aabb.Min,
		Width:    aabb.Max.X - aabb.Min.X,
		Height:   aabb.Max.Y - aabb.Min.Y,
	}
	entity.RenderComponent = &common.RenderComponent{
		Drawable:    common.Rectangle{
			BorderWidth: 1,
			BorderColor: getColor(colorName),
		},
		Color:       color.RGBA{0, 0, 0, 0},
	}
	entity.RenderComponent.SetZIndex(10)
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&entity.BasicEntity, entity.RenderComponent, entity.SpaceComponent)
		}
	}
	self.entities = append(self.entities, entity)
}

func (self *DebugSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.entities {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}

func (self *DebugSystem) Update(dt float32) {
	for _, e := range self.entities {
		if e.removeAfter > 0 {
			now := time.Now()
			if now.Sub(e.shownSince) > e.removeAfter {
				engo.Mailbox.Dispatch(messages.TileRemoveMessage{
					Entity: &e.BasicEntity,
				})
			}
		}
	}
}

func getColor(name string) color.Color {
	switch name {
	case "blue":
		return color.RGBA{0, 0, 255, 255}
	case "green":
		return color.RGBA{255, 0, 0, 255}
	case "red":
		return color.RGBA{0, 255, 0, 255}
	}
	return color.RGBA{255, 255, 255, 255}
}

func (self *DebugSystem) HandleDisplayDebugAABBMessage(m engo.Message) {
	log.Printf("[DebugSystem] %+v", m)
	msg, ok := m.(messages.DisplayDebugAABBMessage)
	if !ok {
		return
	}
	for _, a := range msg.Aabbs {
		self.AddAABB(a, msg.Color, msg.RemoveAfter)
	}
	for _, a := range msg.Aabbers {
		self.AddAABBer(a, msg.Color, msg.RemoveAfter)
	}
}
