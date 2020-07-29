package save

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"gogame/data"
	"gogame/messages"
	"gogame/systems"
	"log"
	"os"
	"fmt"
	"time"
	"encoding/json"
)

func HandleSaveMessage(world *ecs.World, filepath string) {
	log.Println("[SaveGame] preparing the save file")
	// TODO the game should be paused first

	// All systems that save anything should do it here
	saveFile := &data.SaveFile{}
	for _, system := range world.Systems() {
		if sys, ok := system.(*systems.CreatureSpawningSystem); ok {
			sys.UpdateSave(saveFile)
		}
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
		Lines: []string{
			fmt.Sprintf("Saved to %s", filepath),
		},
	})
}

func HandleLoadMessage(world *ecs.World, filepath string) {
	log.Printf("[SaveGame] loading from a save file '%s'", filepath)
	// TODO the game should be paused first
	f2, err := os.Open(filepath)
	dec := json.NewDecoder(f2)
	saveFile := &data.SaveFile{}
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
		Lines: []string{
			fmt.Sprintf("Loaded %s", filepath),
		},
	})
}
