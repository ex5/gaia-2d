package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"math/rand"
	//"log"
	//"math"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.CollisionComponent

	*Matter
}

type Matter struct {
	Type string // plant, stone ?
	Amount float32
	Units string // kcal, pcs
}

func InitWorld(u engo.Updater) {
	world, _ := u.(*ecs.World)

	assets.FullSpriteSheet = common.NewSpritesheetFromFile("tilemap/terrain-v7.png", 32, 32)
	mapSizeX, mapSizeY := 50, 50

	tiles := make([]*Tile, 0)
	z_idx_max := 0.0
	for i := 0; i < mapSizeX; i ++ {
		for j := 0; j < mapSizeY; j ++ {
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			tile.RenderComponent = common.RenderComponent{
				Drawable: assets.FullSpriteSheet.Cell(1127),
				Scale:    engo.Point{1, 1},
			}
			position := engo.Point{float32(i) * 32.0, float32(j) * 32.0}
			tile.RenderComponent.SetZIndex(float32(0))
			tile.SpaceComponent = common.SpaceComponent{
				Position: position,
				Width:    float32(assets.SpriteWidth),
				Height:   float32(assets.SpriteHeight),
			}
			tiles = append(tiles, tile)

			// Add vegeration tile
			if rand.Int() % 3 == 0 {
				vtile := &Tile{BasicEntity: ecs.NewBasic()}
				// TODO how to distinguish walkable terrain?
				// TODO how to store additional objects on top of terrain tiles?
				vtile.Matter = &Matter{Type: "plant", Amount: 100, Units: "kcal"}
				vtile.RenderComponent = common.RenderComponent{
					Drawable: assets.FullSpriteSheet.Cell(1706),
					Scale:    engo.Point{1, 1},
				}
				vtile.RenderComponent.SetZIndex(float32(2))
				vtile.SpaceComponent = common.SpaceComponent{
					Position: position,
					Width:    float32(assets.SpriteWidth),
					Height:   float32(assets.SpriteHeight),
				}
				tiles = append(tiles, vtile)
			}

		}
	}
	z_idx_max = 3.0
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)
	// add the tiles to the RenderSystem
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		case *common.CollisionSystem:
			for _, v := range tiles {
				sys.Add(&v.BasicEntity, &v.CollisionComponent, &v.SpaceComponent)
			}
		}
	}
}
