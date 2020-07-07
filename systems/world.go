package systems

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"log"
	"math"
)

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func InitWorld(u engo.Updater) {
	world, _ := u.(*ecs.World)

	resource, err := engo.Files.Resource(assets.PreloadList[1])
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	tiles := make([]*Tile, 0)
	z_idx_max := 0.0
	for idx, tileLayer := range levelData.TileLayers {
		log.Printf("%+v", tileLayer.Tiles[0])
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement.Image,
					Scale:    engo.Point{1, 1},
				}
				// TODO how to distinguish walkable terrain?
				// TODO how to store additional objects on top of terrain tiles?
				tile.RenderComponent.SetZIndex(float32(idx))
				z_idx_max = math.Max(z_idx_max, float64(idx))
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    float32(assets.SpriteWidth),
					Height:   float32(assets.SpriteHeight),
				}
				tiles = append(tiles, tile)
			} else {
				fmt.Printf("image is nil!\n")
			}
		}
	}
	fmt.Printf("Max Z index of the terrain: %d\n", z_idx_max)
	// add the tiles to the RenderSystem
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range tiles {
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
}
