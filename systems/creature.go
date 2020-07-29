package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/data"
	"gogame/messages"
	"gogame/util"
	"log"
)

type CreatureMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type CreatureSpawningSystem struct {
	world         *ecs.World
	mouseTracker  CreatureMouseTracker
	entityActions []*common.Animation
	entities      []*data.Creature
}

func NewCreature(creatureID int, position *engo.Point) *data.Creature {
	creature := assets.GetCreatureById(creatureID)
	tile := NewTile(creature.ObjectID, position, 4, &common.CollisionComponent{Main: 1, Group: 0})
	entity := &data.Creature{ID: creatureID, Tile: tile}
	return entity
}

func (self *CreatureSpawningSystem) Add(entity *data.Creature) {
	self.entities = append(self.entities, entity)

	// Add the entity to the various systems
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *WorldTilesSystem:
			sys.Add(entity.Tile)
		}
	}
}

// New is the initialisation of the System
func (self *CreatureSpawningSystem) New(w *ecs.World) {
	log.Println("CreatureSpawningSystem was added to the Scene")

	self.world = w

	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			//log.Println("COLLISION")
		}
	})

	engo.Mailbox.Listen(messages.ControlMessageType, func(m engo.Message) {
		log.Printf("%+v", m)
		msg, ok := m.(messages.ControlMessage)
		if !ok {
			return
		}
		if msg.Action == "add_creature" {
			x, y := util.ToGridPosition(self.mouseTracker.MouseX, self.mouseTracker.MouseY)
			c := NewCreature(msg.CreatureID, &engo.Point{x, y})
			self.Add(c)

		}
	})

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *CreatureSpawningSystem) Update(dt float32) {
	//log.Printf("Entities: %+v\n", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		entity.Update(dt)
	}
}

// Remove is called whenever an Creature is removed from the World, in order to remove it from this sytem as well
func (self *CreatureSpawningSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.entities {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	// Also remove from whichever other systems this system might have added the entity to
	for _, system := range self.world.Systems() {
		system.Remove(e)
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}

func (self *CreatureSpawningSystem) UpdateSave(saveFile *data.SaveFile) {
	for _, e := range self.entities {
		if e.Tile.Object.Type == "creature" {
			saveFile.Creatures = append(saveFile.Creatures, e)
		}
	}
}

func (self *CreatureSpawningSystem) LoadSave(saveFile *data.SaveFile) {
	log.Printf("[CreatureSpawningSystem] Creatures in the save file: %d\n", len(saveFile.Creatures))
	for _, c := range saveFile.Creatures {
		self.Add(c)
	}
}
