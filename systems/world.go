package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"gogame/messages"
	"math/rand"
	"log"
	//"math"
)

var (
	tiles []*Tile
)

type Matter struct {
	Type string // plant, stone ?
	Amount float32
}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent
	*common.MouseComponent

	*assets.Object
	*Matter
}

func Add(world *ecs.World, spriteID int, i int, j int, layer float32) *Tile {
	tile := &Tile{BasicEntity: ecs.NewBasic()}
	tile.RenderComponent = common.RenderComponent{
		Drawable: assets.FullSpriteSheet.Cell(spriteID),
		Scale:    engo.Point{1, 1},
	}
	position := engo.Point{float32(i * assets.SpriteWidth), float32(j * assets.SpriteHeight)}
	tile.RenderComponent.SetZIndex(layer)
	tile.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    float32(assets.SpriteWidth),
		Height:   float32(assets.SpriteHeight),
	}
	tile.MouseComponent = &common.MouseComponent{Track: false}
	tiles = append(tiles, tile)

	// Add the tile to the various systems
	for _, system := range world.Systems() {
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

func InitWorld(u engo.Updater) {
	world, _ := u.(*ecs.World)

	assets.InitAssets()
	mapSizeX, mapSizeY := 50, 50

	tiles = make([]*Tile, 0)
	for i := 0; i < mapSizeX; i ++ {
		for j := 0; j < mapSizeY; j ++ {
			Add(world, 1127, i, j, 0)

			// Add a random vegetation tile
			if rand.Int() % 3 == 0 {
				plant := assets.GetRandomObjectOfType("plant")
				vtile := Add(world, plant.ID, i, j, 1)
				vtile.Matter = &Matter{Type: plant.Matter.Type, Amount: plant.Amount}
				vtile.Object = plant
			}

		}
	}
	z_idx_max := 3.0
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)

	engo.Mailbox.Listen(messages.InteractionMessageType, HandleInteractMessage)
}

func GetEntityByID(basicEntityID uint64) *Tile {
	for _, e := range tiles {
		if e.BasicEntity.ID() == basicEntityID {
			return e
		}
	}
	return nil
}

func HandleInteractMessage(m engo.Message) {
	log.Printf("World: %+v", m)
	msg, ok := m.(messages.InteractionMessage)
	if !ok {
		return
	}
	if msg.Action == "mouse_hover" && msg.BasicEntity != nil {
		entity := GetEntityByID(msg.BasicEntity.ID())
		log.Printf("World: %+v", entity)
		if entity != nil {
			engo.Mailbox.Dispatch(messages.HUDTextMessage{
				Line1:          fmt.Sprintf("#%d", entity.BasicEntity.ID()),
				Line3:          fmt.Sprintf("%v", entity.Matter),
				Line2:          fmt.Sprintf("%v", entity.Object),
				Line4:          "<World>",
			})
		}
	}
}
