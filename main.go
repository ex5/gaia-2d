package main

import (
	"encoding/json"
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/controls"
	"gogame/messages"
	"gogame/save"
	"gogame/systems"
	"image/color"
	"log"
	"os"
	"time"
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
	assets.InitAssets()
}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (self *myScene) Setup(u engo.Updater) {
	log.Println("[myScene] Setup")
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
	engo.Input.RegisterButton("NewWorld", engo.KeyF4)
	engo.Input.RegisterButton("QuickSave", engo.KeyF5)
	engo.Input.RegisterButton("QuickLoad", engo.KeyF6)
	engo.Input.RegisterButton("ExitToDesktop", engo.KeyEscape)

	// Visual debug
	world.AddSystem(&systems.DebugSystem{})

	// Spacial (quadtree, pathfinding etc)
	world.AddSystem(&systems.SpacialSystem{})

	// Controls
	world.AddSystem(&controls.ControlsSystem{})

	// World
	world.AddSystem(&systems.WorldTilesSystem{})

	// HUD
	systems.InitHUD(u)
	world.AddSystem(&systems.HUDTextSystem{})

	// In-game time
	world.AddSystem(&systems.TimeSystem{})

	// Creatures and plants
	world.AddSystem(&systems.CreatureSpawningSystem{})
	world.AddSystem(&systems.PlantSpawningSystem{})

	engo.Mailbox.Listen(messages.SaveMessageType, func(m engo.Message) {
		log.Printf("%+v", m)
		msg, ok := m.(messages.SaveMessage)
		if !ok {
			return
		}
		HandleSaveMessage(world, msg.Filepath)
	})
	engo.Mailbox.Listen(messages.LoadMessageType, func(m engo.Message) {
		log.Printf("%+v", m)
		msg, ok := m.(messages.LoadMessage)
		if !ok {
			return
		}
		HandleLoadMessage(world, msg.Filepath)
	})
	engo.Mailbox.Listen(messages.ControlMessageType, func(m engo.Message) {
		log.Printf("%+v", m)
		msg, ok := m.(messages.ControlMessage)
		if !ok {
			return
		}
		if msg.Action == "exit" {
			self.Exit()
		}
		if msg.Action == "ReloadWorld" {
			// Set new scene, forcing to recreate the world
			newScene := &myScene{}
			engo.SetScene(newScene, true)
			if msg.Data != "" {
				engo.Mailbox.Dispatch(messages.LoadMessage{
					Filepath: msg.Data,
				})
			} else {
				engo.Mailbox.Dispatch(messages.ControlMessage{
					Action: "WorldGenerate",
				})
			}
		}
	})
}

func HandleSaveMessage(world *ecs.World, filepath string) {
	log.Println("[SaveGame] preparing the save file")
	// TODO the game should be paused first

	// All systems that save anything should do it here
	saveFile := &save.SaveFile{}
	saveFile.SeenEntityIDs = make(map[uint64]struct{})
	// Collect data from all systems in the fixed order to avoid writing duplicate Tiles
	for _, system := range world.Systems() {
		if sys, ok := system.(*systems.CreatureSpawningSystem); ok {
			sys.UpdateSave(saveFile)
		}
	}
	for _, system := range world.Systems() {
		if sys, ok := system.(*systems.PlantSpawningSystem); ok {
			sys.UpdateSave(saveFile)
		}
	}
	for _, system := range world.Systems() {
		if sys, ok := system.(*systems.WorldTilesSystem); ok {
			sys.UpdateSave(saveFile)
		}
	}

	log.Printf("[SaveGame] writing the save file '%s'", filepath)
	f1, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f1)
	err = enc.Encode(saveFile)
	if err != nil {
		panic(err)
	}
	f1.Close()
	log.Printf(".. Done.\n")

	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:      "EventMessage",
		HideAfter: 3 * time.Second,
		GetText: func() []string {
			return []string{
				fmt.Sprintf("Saved to %s", filepath),
			}
		},
	})
}

func HandleLoadMessage(world *ecs.World, filepath string) {
	log.Printf("[SaveGame] loading from a save file '%s'", filepath)
	// TODO the game should be paused first
	f2, err := os.Open(filepath)
	dec := json.NewDecoder(f2)
	saveFile := &save.SaveFile{}
	err = dec.Decode(&saveFile)
	if err != nil {
		panic(err)
	}
	f2.Close()

	// All systems that save anything should do it here
	for _, system := range world.Systems() {
		if sys, ok := system.(*systems.CreatureSpawningSystem); ok {
			sys.LoadSave(saveFile)
		}
		if sys, ok := system.(*systems.WorldTilesSystem); ok {
			sys.LoadSave(saveFile)
		}
	}

	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:      "EventMessage",
		HideAfter: 3 * time.Second,
		GetText: func() []string {
			return []string{
				fmt.Sprintf("Loaded %s", filepath),
			}
		},
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
		Title:               "Gaea",
		Width:               worldWidth,
		Height:              worldHeight,
		StandardInputs:      true,
		OverrideCloseAction: true,
	}
	engo.Run(opts, &myScene{})
}
