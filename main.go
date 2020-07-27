package main

import (
	"bytes"
	// "fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"gogame/messages"
	"gogame/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"image/color"
	// "math"
	"log"
)

var (
	scrollSpeed float32 = 700

	worldWidth  int = 800
	worldHeight int = 800
)

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "gaea" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {
	if err := engo.Files.Load(assets.PreloadList...); err != nil {
		panic(err)
	}
	engo.Files.LoadReaderData("fonts/arcade_n.ttf", bytes.NewReader(gosmallcaps.TTF))
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (self *myScene) Setup(u engo.Updater) {
	world, _ := u.(*ecs.World)

	// Basic systems and controls
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.CollisionSystem{Solids: 1})
	world.AddSystem(&common.AnimationSystem{})
	world.AddSystem(&common.MouseSystem{})
	kbs := common.NewKeyboardScroller(
		scrollSpeed,
		engo.DefaultHorizontalAxis,
		engo.DefaultVerticalAxis)
	world.AddSystem(kbs)
	world.AddSystem(&common.EdgeScroller{scrollSpeed, 20})
	world.AddSystem(&common.MouseZoomer{-0.125})

	common.SetBackground(color.Black)

	engo.Input.RegisterButton("AddCreature", engo.KeyF1)
	engo.Input.RegisterButton("AddObject", engo.KeyF2)
	engo.Input.RegisterButton("QuickSave", engo.KeyF5)
	engo.Input.RegisterButton("ExitToDesktop", engo.KeyEscape)

	// Controls
	world.AddSystem(&controls.ControlsSystem{})

	// World
	world.AddSystem(&systems.WorldTilesSystem{})

	// HUD
	systems.InitHUD(u)
	world.AddSystem(&systems.HUDTextSystem{})

	// Creatures
	world.AddSystem(&systems.CreatureSpawningSystem{})

	// Solid inanimate Objects
	world.AddSystem(&systems.ObjectSpawningSystem{})

	engo.Mailbox.Listen(messages.ControlMessageType, func(m engo.Message) {
		log.Printf("%+v", m)
		msg, ok := m.(messages.ControlMessage)
		if !ok {
			return
		}
		if msg.Action == "exit" {
			self.Exit()
		}
	})
}

func (*myScene) Exit() {
	log.Println("Exit event called; we can do whatever we want now")
	// TODO Here if you want you can prompt the user if they're sure they want to close
	log.Println("Manually closing")
	engo.Exit()
}

func main() {
	opts := engo.RunOptions{
		Title:          "Gaea",
		Width:          worldWidth,
		Height:         worldHeight,
		StandardInputs: true,
		OverrideCloseAction: true,
	}
	engo.Run(opts, &myScene{})
}
