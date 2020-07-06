package main

import (
	"bytes"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"image"
	"image/color"
	"math"
	//"log"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800

	assets = []string{
		"textures/chick_32x32.png",
		"tilemap/terrain-medium.tmx",
	}

	// Z-indices
	z_idx_hud       int = 999
	z_idx_creatures int = 0
)

type myScene struct{}

type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
} // TODO: move

// Type uniquely defines your game type
func (*myScene) Type() string { return "gaea" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	if err := engo.Files.Load(assets...); err != nil {
		panic(err)
	}
	engo.Files.LoadReaderData("fonts/arcade_n.ttf", bytes.NewReader(gosmallcaps.TTF))
}

func InitWorld(u engo.Updater) {
	world, _ := u.(*ecs.World)

	resource, err := engo.Files.Resource(assets[1])
	if err != nil {
		panic(err)
	}
	tmxResource := resource.(common.TMXResource)
	levelData := tmxResource.Level

	tiles := make([]*Tile, 0)
	z_idx_max := 0.0
	for idx, tileLayer := range levelData.TileLayers {
		for _, tileElement := range tileLayer.Tiles {
			if tileElement.Image != nil {
				tile := &Tile{BasicEntity: ecs.NewBasic()}
				tile.RenderComponent = common.RenderComponent{
					Drawable: tileElement.Image,
					Scale:    engo.Point{1, 1},
				}
				tile.RenderComponent.SetZIndex(float32(idx))
				z_idx_max = math.Max(z_idx_max, float64(idx))
				tile.SpaceComponent = common.SpaceComponent{
					Position: tileElement.Point,
					Width:    0,
					Height:   0,
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

func InitHUD(u engo.Updater) {
	world, _ := u.(*ecs.World)

	hud := HUD{BasicEntity: ecs.NewBasic()}
	ww, wh := engo.WindowWidth(), engo.WindowHeight()
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, wh - (wh / 2)},
		Width:    ww / 2,
		Height:   wh / 2,
	}
	hudImage := image.NewUniform(color.RGBA{205, 205, 205, 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, int(ww/2), int(wh/2))
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud.RenderComponent = common.RenderComponent{
		Drawable: hudTexture,
		Scale:    engo.Point{1, 1},
		Repeat:   common.Repeat,
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(float32(z_idx_hud))

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)

	// Basic systems and controls
	world.AddSystem(&common.RenderSystem{})
	common.SetBackground(color.White)
	world.AddSystem(&common.AnimationSystem{})
	world.AddSystem(&common.MouseSystem{})
	kbs := common.NewKeyboardScroller(
		scrollSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	world.AddSystem(&common.EdgeScroller{scrollSpeed, 20})
	world.AddSystem(&common.MouseZoomer{-0.125})

	engo.Input.RegisterButton("AddCreature", engo.KeyF1)

	// World
	InitWorld(u)

	// HUD
	InitHUD(u)

	// Creatures
	world.AddSystem(&systems.CreatureSpawningSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title:          "Gaea",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}
