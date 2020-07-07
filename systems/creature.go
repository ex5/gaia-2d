package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"log"
)

type CreatureMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Creature struct {
	ecs.BasicEntity
	common.AnimationComponent
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
}

// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way
// suitable for the current implementation.
type System interface {
	// Update is ran every frame, with `dt` being the time
	// in seconds since the last frame
	Update(dt float32)

	// Remove removes an Creature from the System
	Remove(ecs.BasicEntity)
}

type CreatureSpawningSystem struct {
	world           *ecs.World
	mouseTracker    CreatureMouseTracker
	entityActions   []*common.Animation
	entities        []*Creature
}

func (self *Creature) Update(dt float32) {
	//log.Printf("%+v %+v", dt, self.SpaceComponent)
	self.SpaceComponent.Position.X += dt * 100
}

func (self *CreatureSpawningSystem) CreateCreature(point engo.Point, spriteSheet *common.Spritesheet) *Creature {
	entity := &Creature{BasicEntity: ecs.NewBasic()}

	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    float32(assets.SpriteWidth),
		Height:   float32(assets.SpriteHeight),
	}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}
	entity.RenderComponent.SetZIndex(10.0)
	entity.CollisionComponent = common.CollisionComponent{
		Main: 1,
	}
	entity.AnimationComponent = common.NewAnimationComponent(spriteSheet.Drawables(), 0.25)

	entity.AnimationComponent.AddAnimations(self.entityActions)
	entity.AnimationComponent.AddDefaultAnimation(self.entityActions[0])

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
			log.Println("COLLISION")
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
	if engo.Input.Button("AddCreature").JustPressed() {
		log.Println("The gamer pressed F1")
		spriteSheet := common.NewSpritesheetFromFile("textures/chick_32x32.png", assets.SpriteWidth, assets.SpriteHeight)
		creature := self.CreateCreature(engo.Point{self.mouseTracker.MouseX, self.mouseTracker.MouseY}, spriteSheet)

		for _, system := range self.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&creature.BasicEntity, &creature.RenderComponent, &creature.SpaceComponent)
			case *common.AnimationSystem:
				sys.Add(&creature.BasicEntity, &creature.AnimationComponent, &creature.RenderComponent)
			case *common.CollisionSystem:
				sys.Add(&creature.BasicEntity, &creature.CollisionComponent, &creature.SpaceComponent)
			}
		}
	}

	//log.Printf("Entities: %+v", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		entity.Update(dt)
	}
}

// Remove is called whenever an Creature is removed from the World, in order to remove it from this sytem as well
func (*CreatureSpawningSystem) Remove(ecs.BasicEntity) {}
