package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/config"
	"gogame/controls"
	"gogame/data"
	"gogame/messages"
	"gogame/save"
	"gogame/util"
	"log"
	"math/rand"
)

type WorldTilesSystem struct {
	world *ecs.World
	tiles []*data.Tile // TODO entities
}

func NewTile(objectID int, position *engo.Point, layer float32, collisionComponent *common.CollisionComponent) *data.Tile {
	basic := ecs.NewBasic()
	tile := &data.Tile{BasicEntity: &basic, Layer: layer, ObjectID: objectID}
	tile.SpaceComponent = &common.SpaceComponent{
		Position: *position,
		Width:    float32(config.SpriteWidth),
		Height:   float32(config.SpriteHeight),
	}
	tile.MouseComponent = &common.MouseComponent{Track: false}
	tile.CollisionComponent = collisionComponent
	return tile
}

func (self *WorldTilesSystem) Add(tile *data.Tile) {
	if tile.BasicEntity == nil {
		basic := ecs.NewBasic()
		tile.BasicEntity = &basic
	}
	if tile.Object == nil {
		tile.Object = assets.GetObjectById(tile.ObjectID)
	}
	if tile.Resource == nil && tile.Object.ResourceID != 0 {
		tile.Resource = assets.GetResourceByID(tile.Object.ResourceID)
	}
	if tile.RenderComponent == nil {
		tile.RenderComponent = &common.RenderComponent{
			Drawable: tile.Object.Spritesheet.Cell(tile.Object.SpriteID),
			Scale:    engo.Point{1, 1},
		}
		tile.RenderComponent.SetZIndex(tile.Layer)
	}
	if tile.AccessibleResource == nil {
		tile.AccessibleResource = &data.AccessibleResource{tile.Object.ResourceID, tile.Object.Amount}
	}
	if tile.AnimationComponent == nil {
		if tile.Object.Animations != nil && len(tile.Object.Animations) > 1 {
			log.Printf("Adding animations to %+v\n", tile)
			animationC := common.NewAnimationComponent(tile.Object.Spritesheet.Drawables(), 0.25)
			tile.AnimationComponent = &animationC
			tile.AnimationComponent.AddAnimations(tile.Object.Animations)
			tile.AnimationComponent.AddDefaultAnimation(tile.Object.Animations[0])
		}
	}
	self.tiles = append(self.tiles, tile)

	// Add the tile to the various systems
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(tile.BasicEntity, tile.RenderComponent, tile.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(tile.BasicEntity, tile.CollisionComponent, tile.SpaceComponent)
		case *common.AnimationSystem:
			if tile.AnimationComponent != nil {
				sys.Add(tile.BasicEntity, tile.AnimationComponent, tile.RenderComponent)
			}
		case *controls.ControlsSystem:
			sys.Add(tile.BasicEntity, tile.MouseComponent, tile.SpaceComponent, tile.RenderComponent)
		case *SpacialSystem:
			sys.Add(tile)
		}
	}
}

func (self *WorldTilesSystem) ReplaceObject(tile *data.Tile, objectID int) {
	tile.ObjectID = objectID
	tile.Object = assets.GetObjectById(objectID)
	if tile.Resource == nil && tile.Object.ResourceID != 0 {
		tile.Resource = assets.GetResourceByID(tile.Object.ResourceID)
	}
	tile.RenderComponent.Drawable = tile.Object.Spritesheet.Cell(tile.Object.SpriteID)
	tile.AccessibleResource = &data.AccessibleResource{tile.Object.ResourceID, tile.Object.Amount}
	// TODO update animations if any
}

func (self *WorldTilesSystem) New(world *ecs.World) {
	self.world = world
	self.tiles = make([]*data.Tile, 0)

	engo.Mailbox.Listen(messages.InteractionMessageType, self.HandleInteractMessage)
	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
	engo.Mailbox.Listen(messages.TileRemoveMessageType, self.HandleTileRemoveMessage)
	engo.Mailbox.Listen(messages.TileReplaceMessageType, self.HandleTileReplaceMessage)
}

func (self *WorldTilesSystem) Generate() {
	mapSizeX, mapSizeY := 50, 50
	groundID := 4 // grassland, default ground
	// ground doesn't collide with anything
	collisionC := &common.CollisionComponent{Main: 0, Group: 0}
	for i := 0; i < mapSizeX; i++ {
		for j := 0; j < mapSizeY; j++ {
			position := util.ToPoint(i, j)
			tile := NewTile(groundID, position, 0, collisionC)
			self.Add(tile)

			// Add a random vegetation
			if rand.Int()%3 == 0 {
				engo.Mailbox.Dispatch(messages.NewPlantMessage{
					Point:   position,
					PlantID: 1,
				})
			}

		}
	}
	common.CameraBounds = engo.AABB{
		Min: engo.Point{0, 0},
		Max: engo.Point{
			float32(mapSizeX * config.SpriteWidth),
			float32(mapSizeY * config.SpriteHeight),
		},
	}
	common.MaxZoom = 1.5
}

func (self *WorldTilesSystem) GetEntityByID(basicEntityID uint64) *data.Tile {
	for _, e := range self.tiles {
		if e.BasicEntity.ID() == basicEntityID {
			return e
		}
	}
	return nil
}

func (self *WorldTilesSystem) HandleInteractMessage(m engo.Message) {
	log.Printf("World: %+v", m)
	msg, ok := m.(messages.InteractionMessage)
	if !ok {
		return
	}
	if msg.Action == "mouse_hover" && msg.BasicEntity != nil {
		entity := self.GetEntityByID(msg.BasicEntity.ID())
		log.Printf("World: %+v", entity)
		if entity != nil {
			engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
				Name:    "HoverInfo",
				GetText: entity.GetTextStatus,
			})
			// FIXME must be a better way than this?
			engo.Mailbox.Dispatch(messages.CreatureHoveredMessage{
				EntityID: entity.BasicEntity.ID(),
			})
			engo.Mailbox.Dispatch(messages.PlantHoveredMessage{
				EntityID: entity.BasicEntity.ID(),
			})
		}
	}
}

func (self *WorldTilesSystem) HandleControlMessage(m engo.Message) {
	log.Printf("%+v", m)
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	if msg.Action == "add_object" {
		// TODO disallow placing on top of an existing overlaying Tile
		for _, system := range self.world.Systems() {
			controlsSystem, ok := system.(*controls.ControlsSystem)
			if ok {
				x, y := util.ToGridPosition(controlsSystem.MouseTracker.MouseX, controlsSystem.MouseTracker.MouseY)
				tile := NewTile(6, &engo.Point{x, y}, 4, &common.CollisionComponent{Main: 0, Group: 1})
				self.Add(tile)
				break
			}
		}
	} else if msg.Action == "WorldGenerate" {
		// TODO the game should be paused first
		self.Generate()
	}
}

func (self *WorldTilesSystem) HandleTileRemoveMessage(m engo.Message) {
	//log.Printf("%+v", m)
	msg, ok := m.(messages.TileRemoveMessage)
	if !ok {
		return
	}
	self.world.RemoveEntity(*msg.Entity)
}

func (self *WorldTilesSystem) HandleTileReplaceMessage(m engo.Message) {
	//log.Printf("%+v", m)
	msg, ok := m.(messages.TileReplaceMessage)
	if !ok {
		return
	}
	tile := self.GetEntityByID(msg.Entity.ID())
	self.ReplaceObject(tile, msg.ObjectID)
}

func (self *WorldTilesSystem) Update(dt float32) {}

func (self *WorldTilesSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.tiles {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	if delete >= 0 {
		self.tiles = append(self.tiles[:delete], self.tiles[delete+1:]...)
	}
}

func (self *WorldTilesSystem) UpdateSave(saveFile *save.SaveFile) {
	for _, t := range self.tiles {
		entityID := t.BasicEntity.ID()
		if _, ok := saveFile.SeenEntityIDs[entityID]; !ok {
			saveFile.Tiles = append(saveFile.Tiles, t)
			saveFile.SeenEntityIDs[entityID] = struct{}{}
		}
	}
}

func (self *WorldTilesSystem) LoadSave(saveFile *save.SaveFile) {
	log.Printf("[WorldTilesSystem] Tiles in the save file: %d\n", len(saveFile.Tiles))
	for _, t := range saveFile.Tiles {
		self.Add(t)
	}
}
