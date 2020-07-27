package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"gogame/messages"
	"gogame/util"
	"log"
)

type ObjectMouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

type Object struct {
	ecs.BasicEntity
	*common.Spritesheet
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	common.MouseComponent
}

type ObjectSpawningSystem struct {
	world        *ecs.World
	mouseTracker ObjectMouseTracker
	spritesheets map[string]*common.Spritesheet
	entities     []*Object
}

func (self *Object) Update(dt float32) {}

func (self *ObjectSpawningSystem) GetOrLoadSpritesheet(spriteSrc string) *common.Spritesheet {
	_, exists := self.spritesheets[spriteSrc]
	if !exists {
		log.Printf("Loading a new sprite source: %s\n", spriteSrc)
		self.spritesheets[spriteSrc] = common.NewSpritesheetFromFile(spriteSrc, assets.SpriteWidth, assets.SpriteHeight)
	} else {
		log.Printf("%s already loaded\n", spriteSrc)
	}
	return self.spritesheets[spriteSrc]
}

func (self *ObjectSpawningSystem) CreateObjectFromSpriteSource(point engo.Point, spriteSrc string, collisionMain bool) *Object {
	spriteSheet := self.GetOrLoadSpritesheet(spriteSrc)
	entity := &Object{BasicEntity: ecs.NewBasic(), Spritesheet: spriteSheet}
	entity.RenderComponent = common.RenderComponent{
		Drawable: spriteSheet.Cell(0),
		Scale:    engo.Point{1, 1},
	}

	return self.AddObjectEntity(point, entity, collisionMain)
}

func (self *ObjectSpawningSystem) CreateObjectFromTextureAtlas(point engo.Point, spriteID int, collisionMain bool) *Object {
	entity := &Object{BasicEntity: ecs.NewBasic(), Spritesheet: assets.FullSpriteSheet}
	entity.RenderComponent = common.RenderComponent{
		Drawable: assets.FullSpriteSheet.Cell(spriteID),
		Scale:    engo.Point{1, 1},
	}

	return self.AddObjectEntity(point, entity, collisionMain)
}

func (self *ObjectSpawningSystem) AddObjectEntity(point engo.Point, entity *Object, collisionMain bool) *Object {
	entity.SpaceComponent = common.SpaceComponent{
		Position: point,
		Width:    float32(assets.SpriteWidth),
		Height:   float32(assets.SpriteHeight),
	}
	entity.RenderComponent.SetZIndex(10.0)
	if collisionMain {
		entity.CollisionComponent = common.CollisionComponent{
			Main: 1,
		}
	} else {
		entity.CollisionComponent = common.CollisionComponent{
			Group: 1,
		}
	}
	entity.MouseComponent = common.MouseComponent{Track: false}

	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&entity.BasicEntity, &entity.RenderComponent, &entity.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&entity.BasicEntity, &entity.CollisionComponent, &entity.SpaceComponent)
		//case *common.MouseSystem:
		//	sys.Add(&entity.BasicEntity, &entity.MouseComponent, &entity.SpaceComponent, &entity.RenderComponent)
		case *controls.ControlsSystem:
			sys.Add(&entity.BasicEntity, &entity.MouseComponent, &entity.SpaceComponent, &entity.RenderComponent)
		}
	}

	self.entities = append(self.entities, entity)

	return entity
}

func (self *ObjectSpawningSystem) HandleControlMessage(m engo.Message) {
	log.Printf("%+v", m)
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	if msg.Action == "add_object" {
		x, y := util.ToGridPosition(self.mouseTracker.MouseX, self.mouseTracker.MouseY)
		self.CreateObjectFromTextureAtlas(engo.Point{x, y}, msg.SpriteID, false)
	}
}

func (self *ObjectSpawningSystem) HandleInteractMessage(m engo.Message) {
	log.Printf("Objects: %+v", m)
	msg, ok := m.(messages.InteractionMessage)
	if !ok {
		return
	}
	if msg.Action == "mouse_hover" && msg.BasicEntity != nil {
		entity := self.GetEntityByID(msg.BasicEntity.ID())
		log.Printf("Objects: %+v", entity)
		if entity != nil {
			engo.Mailbox.Dispatch(messages.HUDTextMessage{
				Line1: fmt.Sprintf("#%d", entity.ID()),
				Line2: fmt.Sprintf("sprite: %s", entity.Spritesheet),
				Line3: "<Object>",
				Line4: "",
			})
		}
	}
}

func (self *ObjectSpawningSystem) GetEntityByID(basicEntityID uint64) *Object {
	for _, e := range self.entities {
		if e.ID() == basicEntityID {
			return e
		}
	}
	return nil
}

// New is the initialisation of the System
func (self *ObjectSpawningSystem) New(w *ecs.World) {
	log.Println("ObjectSpawningSystem was added to the Scene")

	self.world = w
	self.spritesheets = make(map[string]*common.Spritesheet)

	self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}

	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
	engo.Mailbox.Listen(messages.InteractionMessageType, self.HandleInteractMessage)
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *ObjectSpawningSystem) Update(dt float32) {
	//log.Printf("Entities: %+v", self.entities)
	for _, entity := range self.entities {
		//log.Printf("Entity: %d, %+v", i, entity)
		entity.Update(dt)
	}
}

// Remove is called whenever an Object is removed from the World, in order to remove it from this sytem as well
func (*ObjectSpawningSystem) Remove(ecs.BasicEntity) {}
