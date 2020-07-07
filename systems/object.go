package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"log"
)

type ObjectMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Object struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	common.MouseComponent
}

type ObjectSpawningSystem struct {
	world           *ecs.World
	mouseTracker    ObjectMouseTracker
	entities        []*Object
}

func (self *Object) Update(dt float32) {}

func (self *ObjectSpawningSystem) CreateObject(point engo.Point, spriteSheet *common.Spritesheet) *Object {
	entity := &Object{BasicEntity: ecs.NewBasic()}

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
		Group: 1,
	}
	entity.MouseComponent = common.MouseComponent{Track: false}

	self.entities = append(self.entities, entity)

	return entity
}

// New is the initialisation of the System
func (self *ObjectSpawningSystem) New(w *ecs.World) {
	log.Println("ObjectSpawningSystem was added to the Scene")

	self.world = w

	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *ObjectSpawningSystem) Update(dt float32) {
	if engo.Input.Button("AddObject").JustPressed() {
		log.Println("The gamer pressed F2")
		spriteSheet := common.NewSpritesheetFromFile("textures/stone_32x32.png", assets.SpriteWidth, assets.SpriteHeight)
		object := self.CreateObject(engo.Point{self.mouseTracker.MouseX, self.mouseTracker.MouseY}, spriteSheet)

		for _, system := range self.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&object.BasicEntity, &object.RenderComponent, &object.SpaceComponent)
			case *common.CollisionSystem:
				sys.Add(&object.BasicEntity, &object.CollisionComponent, &object.SpaceComponent)
			case *common.MouseSystem:
				sys.Add(&object.BasicEntity, &object.MouseComponent, &object.SpaceComponent, &object.RenderComponent)
			case *controls.ControlsSystem:
				sys.Add(&object.BasicEntity, &object.MouseComponent)
			}
		}
	}

	//log.Printf("Entities: %+v", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		entity.Update(dt)
	}
}

// Remove is called whenever an Object is removed from the World, in order to remove it from this sytem as well
func (*ObjectSpawningSystem) Remove(ecs.BasicEntity) {}
