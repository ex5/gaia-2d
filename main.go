package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

    "gogame/systems"
    // "log"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

type myScene struct {}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
    engo.Files.Load("textures/chick_24x24.png")
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(u engo.Updater) {
    world, _ := u.(*ecs.World)

	world.AddSystem(&common.RenderSystem{})
    //world.AddSystem(&common.AnimationSystem{})
    world.AddSystem(&common.MouseSystem{})
    kbs := common.NewKeyboardScroller(
		scrollSpeed,
        engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
    world.AddSystem(&common.EdgeScroller{scrollSpeed, 20})

    engo.Input.RegisterButton("AddAnimal", engo.KeyF1)
    engo.Input.RegisterButton("horizontal", engo.KeyD)
    engo.Input.RegisterButton("vertical", engo.KeyS)

	world.AddSystem(&systems.AnimalSpawningSystem{})
}

func main() {
	opts := engo.RunOptions{
		Title: "Gaea",
		Width:  worldWidth,
		Height: worldHeight,
        StandardInputs: true,
	}
	engo.Run(opts, &myScene{})
}
