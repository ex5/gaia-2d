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
	//"math"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	*common.MouseComponent

	Layer float32
	*assets.Object
	*assets.AccessibleResource
	*assets.Resource
}

type WorldTilesSystem struct {
	world *ecs.World
	tiles []*Tile // TODO entities
}

func (self *WorldTilesSystem) Add(spriteID int, position *engo.Point, layer float32) *Tile {
	tile := &Tile{BasicEntity: ecs.NewBasic()}
	tile.RenderComponent = common.RenderComponent{
		Drawable: assets.FullSpriteSheet.Cell(spriteID),
		Scale:    engo.Point{1, 1},
	}
	tile.Layer = layer
	tile.RenderComponent.SetZIndex(tile.Layer)
	tile.SpaceComponent = common.SpaceComponent{
		Position: *position,
		Width:    float32(assets.SpriteWidth),
		Height:   float32(assets.SpriteHeight),
	}
	tile.MouseComponent = &common.MouseComponent{Track: false}
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

func (self *WorldTilesSystem) New(world *ecs.World) {
	self.world = world

	assets.InitAssets()

	self.tiles = make([]*Tile, 0)
	//self.Generate()
	self.LoadFromSaveFile("quick.save")

	engo.Mailbox.Listen(messages.InteractionMessageType, self.HandleInteractMessage)
	engo.Mailbox.Listen(messages.SaveMessageType, self.HandleSaveMessage)
}

func (self *WorldTilesSystem) Generate() {
	mapSizeX, mapSizeY := 50, 50
	ground := assets.GetObjectById(3) // grassland, default ground
	for i := 0; i < mapSizeX; i++ {
		for j := 0; j < mapSizeY; j++ {
			position := util.ToPoint(i, j)
			tile := self.Add(ground.SpriteID, position, 0)
			tile.Object = ground

			// Add a random vegetation tile
			if rand.Int()%3 == 0 {
				plant := assets.GetRandomObjectOfType("plant")
				vtile := self.Add(plant.SpriteID, position, 1)
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
			engo.Mailbox.Dispatch(messages.HUDTextMessage{
				Line1: fmt.Sprintf("#%d", entity.BasicEntity.ID()),
				Line2: fmt.Sprintf("%v", entity.Object),
				Line3: fmt.Sprintf("%v", entity.AccessibleResource),
				Line4: fmt.Sprintf("%v", entity.Resource),
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

func (*WorldTilesSystem) Update(dt float32) {}

func (*WorldTilesSystem) Remove(ecs.BasicEntity) {}

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
		vtile := self.Add(object.SpriteID, savedTile.Position, savedTile.Layer)
		vtile.AccessibleResource = savedTile.AccessibleResource
		vtile.Resource = assets.GetResourceByID(object.ResourceID)
		vtile.Object = object
	}
	z_idx_max := 3.0
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)
}
