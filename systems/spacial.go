package systems

import (
	//"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"gogame/config"
	"gogame/messages"
	"log"
	"time"
	//"math/rand"
)

type SpacialSystem struct {
	world    *ecs.World
	entities *engo.Quadtree
}

func (self *SpacialSystem) New(world *ecs.World) {
	self.world = world
	self.entities = engo.NewQuadtree(
		engo.AABB{
			Min: engo.Point{0, 0},
			Max: engo.Point{float32(config.SpriteWidth * 256), float32(config.SpriteHeight * 256)}},
		true, 16)

	engo.Mailbox.Listen(messages.SpacialRequestMessageType, self.HandleSpacialRequestMessage)
}

func (self *SpacialSystem) Add(e engo.AABBer) {
	self.entities.Insert(e)
}

func (self *SpacialSystem) Update(dt float32) {}

func (self *SpacialSystem) Remove(e ecs.BasicEntity) {
	// TODO
}

func (self *SpacialSystem) Query(aabb engo.AABB, filter func(aabb engo.AABBer) bool) []engo.AABBer {
	engo.Mailbox.Dispatch(messages.DisplayDebugAABBMessage{
		Aabbs:       []engo.AABB{aabb},
		RemoveAfter: 3 * time.Second,
		Color:       "red",
	})
	return self.entities.Retrieve(aabb, filter)
}

func (self *SpacialSystem) HandleSpacialRequestMessage(m engo.Message) {
	log.Printf("[SpacialSystem] %+v", m)
	msg, ok := m.(messages.SpacialRequestMessage)
	if !ok {
		return
	}
	foundEntities := self.Query(msg.Aabb, msg.Filter)
	engo.Mailbox.Dispatch(messages.DisplayDebugAABBMessage{
		Aabbers:     foundEntities,
		RemoveAfter: 3 * time.Second,
		Color:       "green",
	})
	engo.Mailbox.Dispatch(messages.SpacialResponseMessage{
		Aabb:     msg.Aabb,
		EntityID: msg.EntityID,
		Filter:   msg.Filter,
		Result:   foundEntities,
	})
}
