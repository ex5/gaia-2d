package systems

import (
	"encoding/json"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"gogame/messages"
	"gogame/util"
	"log"
	"math/rand"
	"os"
	"time"
	//"math"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	*common.MouseComponent
	*common.Spritesheet

	Layer float32
	*assets.Object
	*assets.AccessibleResource
	*assets.Resource
}

type WorldTilesSystem struct {
	world *ecs.World
	tiles []*Tile // TODO entities
}

func (self *WorldTilesSystem) Add(spritesheet *common.Spritesheet, spriteID int, position *engo.Point, layer float32, collisionComponent common.CollisionComponent) *Tile {
	tile := &Tile{BasicEntity: ecs.NewBasic(), Spritesheet: spritesheet, Layer: layer}
	tile.RenderComponent = common.RenderComponent{
		Drawable: tile.Spritesheet.Cell(spriteID),
		Scale:    engo.Point{1, 1},
	}
	tile.RenderComponent.SetZIndex(tile.Layer)
	tile.SpaceComponent = common.SpaceComponent{
		Position: *position,
		Width:    float32(assets.SpriteWidth),
		Height:   float32(assets.SpriteHeight),
	}
	tile.MouseComponent = &common.MouseComponent{Track: false}
	tile.CollisionComponent = collisionComponent
	self.tiles = append(self.tiles, tile)

	// Add the tile to the various systems
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&tile.BasicEntity, &tile.CollisionComponent, &tile.SpaceComponent)
		case *controls.ControlsSystem:
			sys.Add(&tile.BasicEntity, tile.MouseComponent, &tile.SpaceComponent, &tile.RenderComponent)
		}
	}
	return tile
}

func (self *WorldTilesSystem) CreateFromSpriteSource(point *engo.Point, spriteSrc string, collisionComponent common.CollisionComponent) *Tile {
	spritesheet := assets.GetSpritesheet(spriteSrc)
	return self.Add(spritesheet, 0, point, 10, collisionComponent)
}

func (self *WorldTilesSystem) New(world *ecs.World) {
	self.world = world
	self.tiles = make([]*Tile, 0)

	engo.Mailbox.Listen(messages.InteractionMessageType, self.HandleInteractMessage)
	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
	engo.Mailbox.Listen(messages.SaveMessageType, self.HandleSaveMessage)
}

func (self *WorldTilesSystem) Generate() {
	mapSizeX, mapSizeY := 50, 50
	ground := assets.GetObjectById(3) // grassland, default ground
	// ground doesn't collide with anything
	collisionC := common.CollisionComponent{Main: 0, Group: 0}
	for i := 0; i < mapSizeX; i++ {
		for j := 0; j < mapSizeY; j++ {
			position := util.ToPoint(i, j)
			tile := self.Add(assets.FullSpriteSheet, ground.SpriteID, position, 0, collisionC)
			tile.Object = ground

			// Add a random vegetation tile
			if rand.Int()%3 == 0 {
				plant := assets.GetRandomObjectOfType("plant")
				vtile := self.Add(assets.FullSpriteSheet, plant.SpriteID, position, 1, collisionC)
				vtile.AccessibleResource = &assets.AccessibleResource{plant.ResourceID, plant.Amount}
				vtile.Resource = assets.GetResourceByID(plant.ResourceID)
				vtile.Object = plant
			}

		}
	}
	z_idx_max := 3.0
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)
}

func (self *WorldTilesSystem) GetEntityByID(basicEntityID uint64) *Tile {
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
			lines := []string{
				fmt.Sprintf("#%d", entity.BasicEntity.ID()),
				fmt.Sprintf("%v", entity.Object),
				fmt.Sprintf("%v", entity.AccessibleResource),
				fmt.Sprintf("%v", entity.Resource),
			}
			engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
				Name:  "HoverInfo",
				Lines: lines,
			})
		}
	}
}

func (self *WorldTilesSystem) HandleSaveMessage(m engo.Message) {
	log.Printf("World.HandleSaveMessage: %+v", m)
	msg, ok := m.(messages.SaveMessage)
	if !ok {
		return
	}
	self.Save(msg.Filepath)
}

func (self *WorldTilesSystem) HandleControlMessage(m engo.Message) {
	log.Printf("%+v", m)
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	if msg.Action == "add_object" {
		for _, system := range self.world.Systems() {
			controlsSystem, ok := system.(*controls.ControlsSystem)
			if ok {
				x, y := util.ToGridPosition(controlsSystem.MouseTracker.MouseX, controlsSystem.MouseTracker.MouseY)
				self.Add(assets.FullSpriteSheet, msg.SpriteID, &engo.Point{x, y}, 4, common.CollisionComponent{Main: 0, Group: 1})
				break
			}
		}
	} else if msg.Action == "SaveGame" {
		// TODO the game should be paused first
		self.Save(msg.Data)
	} else if msg.Action == "WorldLoadSaveFile" {
		// TODO the game should be paused first
		self.LoadFromSaveFile(msg.Data)
	} else if msg.Action == "WorldGenerate" {
		// TODO the game should be paused first
		self.Generate()
	}
}

func (*WorldTilesSystem) Update(dt float32) {}

func (self *WorldTilesSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.tiles {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	// Also remove from whichever other systems this system might have added the entity to
	for _, system := range self.world.Systems() {
		system.Remove(e)
	}
	if delete >= 0 {
		self.tiles = append(self.tiles[:delete], self.tiles[delete+1:]...)
	}
}

// Saving-related functionality
func (self *Tile) ToSavedTile() *assets.SavedTile {
	var objectID int
	if self.Object != nil {
		objectID = self.Object.ID
	}
	return &assets.SavedTile{objectID, self.AccessibleResource, &self.SpaceComponent.Position, self.Layer}
}

func (self *WorldTilesSystem) ToSavedTiles() *assets.SavedTiles {
	result := &assets.SavedTiles{}
	for _, v := range self.tiles {
		result.Tiles = append(result.Tiles, v.ToSavedTile())
	}
	return result
}

func (self *WorldTilesSystem) Save(filepath string) {
	log.Println("Saving the world tiles")
	f1, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f1)
	err = enc.Encode(self.ToSavedTiles())
	if err != nil {
		panic(err)
	}
	f1.Close()

	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:      "EventMessage",
		HideAfter: 3 * time.Second,
		Lines: []string{
			fmt.Sprintf("Saved to %s", filepath),
		},
	})
}

func (self *WorldTilesSystem) LoadFromSaveFile(filepath string) {
	f2, err := os.Open(filepath)
	dec := json.NewDecoder(f2)
	savedTiles := &assets.SavedTiles{}
	err = dec.Decode(&savedTiles)
	if err != nil {
		panic(err)
	}
	f2.Close()

	for _, savedTile := range savedTiles.Tiles {
		log.Println(savedTile, savedTile.ObjectID, savedTile.Position, savedTile.Layer)
		object := assets.GetObjectById(savedTile.ObjectID)
		log.Println(object)
		vtile := self.Add(assets.FullSpriteSheet, object.SpriteID, savedTile.Position, savedTile.Layer, common.CollisionComponent{Main: 0, Group: 0}) // TODO load collision status/group etc
		vtile.AccessibleResource = savedTile.AccessibleResource
		vtile.Resource = assets.GetResourceByID(object.ResourceID)
		vtile.Object = object
	}

	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:      "EventMessage",
		HideAfter: 3 * time.Second,
		Lines: []string{
			fmt.Sprintf("Loaded %s", filepath),
		},
	})
	z_idx_max := 3.0
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)
}
