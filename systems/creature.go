package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/messages"
	"gogame/util"
	"log"
)

type CreatureMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Creature struct {
	Object *Object
	common.AnimationComponent
}

type CreatureSpawningSystem struct {
	world           *ecs.World
	mouseTracker    CreatureMouseTracker
	entityActions   []*common.Animation
	entities        []*Creature
}

func (self *Creature) Update(dt float32) {
	//log.Printf("%+v %+v", dt, self.SpaceComponent)
	self.Object.SpaceComponent.Position.X += dt * 100
}

func (self *CreatureSpawningSystem) CreateCreature(point engo.Point, spriteSrc string) *Creature {
	var entity *Creature
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *ObjectSpawningSystem:
			entity = &Creature{Object: sys.CreateObjectFromSpriteSource(point, spriteSrc, true)}
		}
	}

	entity.AnimationComponent = common.NewAnimationComponent(entity.Object.Spritesheet.Drawables(), 0.25)
	entity.AnimationComponent.AddAnimations(self.entityActions)
	entity.AnimationComponent.AddDefaultAnimation(self.entityActions[0])

	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.AnimationSystem:
			sys.Add(&entity.Object.BasicEntity, &entity.AnimationComponent, &entity.Object.RenderComponent)
		}
	}

	self.entities = append(self.entities, entity)

	return entity
}

// New is the initialisation of the System
func (self *CreatureSpawningSystem) New(w *ecs.World) {
	log.Println("CreatureSpawningSystem was added to the Scene")

	self.world = w

	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	self.entityActions = []*common.Animation{
		&common.Animation{Name: "feed", Frames: []int{4, 4, 5, 6, 7, 8, 9, 9}},
		&common.Animation{Name: "walk_right", Frames: []int{4, 4, 0, 1, 2, 3, 4, 4, 4}, Loop: true},
		&common.Animation{Name: "walk_left", Frames: []int{9, 9, 10, 11, 12, 13, 9, 9}, Loop: true},
		&common.Animation{Name: "walk_down", Frames: []int{14, 15, 16, 17, 18, 19, 20}, Loop: true},
		&common.Animation{Name: "walk_up", Frames: []int{21, 22, 23, 24, 25, 26, 27}, Loop: true},
	}
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
			self.CreateCreature(engo.Point{x, y}, msg.Data)

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
	//log.Printf("Entities: %+v", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		entity.Update(dt)
	}
}

// Remove is called whenever an Creature is removed from the World, in order to remove it from this sytem as well
func (*CreatureSpawningSystem) Remove(ecs.BasicEntity) {}
