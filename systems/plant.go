package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/ulule/deepcopier"
	"gogame/life/plants"
	"gogame/messages"
	"gogame/save"
	"gogame/shaders"
	"log"
)

type PlantSpawningSystem struct {
	world    *ecs.World
	shader   common.Shader
	entities []*plants.Plant
}

func NewPlant(plantID int, position *engo.Point) *plants.Plant {
	plant := plants.GetPlantByID(plantID)
	tile := NewTile(plant.ObjectID, position, 2, &common.CollisionComponent{Main: 0, Group: 0})
	entity := &plants.Plant{ID: plantID, Tile: tile}
	// Initialise plant's stats from its initial record
	deepcopier.Copy(plant).To(entity)
	return entity
}

func (self *PlantSpawningSystem) Add(entity *plants.Plant) {
	self.entities = append(self.entities, entity)

	// Add the entity to the various systems
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *WorldTilesSystem:
			sys.Add(entity.Tile)
			entity.Tile.RenderComponent.SetShader(self.shader)
		}
	}
}

// New is the initialisation of the System
func (self *PlantSpawningSystem) New(w *ecs.World) {
	log.Println("PlantSpawningSystem was added to the Scene")

	self.world = w
	self.shader = shaders.WindShader

	engo.Mailbox.Listen(messages.NewPlantMessageType, self.HandleNewPlantMessage)
	engo.Mailbox.Listen(messages.PlantHoveredMessageType, self.HandlePlantHoveredMessage)
	engo.Mailbox.Listen(messages.TimeSecondPassedMessageType, self.HandleTimeSecondPassedMessage)
	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
}

func (self *PlantSpawningSystem) changeShader(wave float32, speed float32) {
	self.shader = shaders.WindShader
	// TODO can be changed elsewhere
	shader, ok := self.shader.(*shaders.BasicShader)
	if !ok {
		panic("not a shader we've expected")
	}
	shader.Wave.Y = wave
	shader.Speed = speed
}

func (self *PlantSpawningSystem) getShader() (float32, float32) {
	self.shader = shaders.WindShader
	// TODO can be changed elsewhere
	shader, ok := self.shader.(*shaders.BasicShader)
	if !ok {
		panic("not a shader we've expected")
	}
	return shader.Wave.Y, shader.Speed
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *PlantSpawningSystem) Update(dt float32) {}

// Remove is called whenever an Plant is removed from the World, in order to remove it from this sytem as well
func (self *PlantSpawningSystem) Remove(e ecs.BasicEntity) {
	delete := -1
	for index, entity := range self.entities {
		if entity.BasicEntity.ID() == e.ID() {
			delete = index
		}
	}
	if delete >= 0 {
		self.entities = append(self.entities[:delete], self.entities[delete+1:]...)
	}
}

func (self *PlantSpawningSystem) Get(entityID uint64) *plants.Plant {
	for _, e := range self.entities {
		if e.BasicEntity.ID() == entityID {
			return e
		}
	}
	return nil
}

func (self *PlantSpawningSystem) HandleNewPlantMessage(m engo.Message) {
	log.Printf("%+v", m)
	msg, ok := m.(messages.NewPlantMessage)
	if !ok {
		return
	}
	e := NewPlant(msg.PlantID, msg.Point)
	self.Add(e)
}

func (self *PlantSpawningSystem) HandlePlantHoveredMessage(m engo.Message) {
	msg, ok := m.(messages.PlantHoveredMessage)
	if !ok {
		return
	}
	entity := self.Get(msg.EntityID)
	if entity == nil {
		return
	}
	engo.Mailbox.Dispatch(messages.HUDTextUpdateMessage{
		Name:    "HoverInfo",
		GetText: entity.GetTextStatus,
	})
}

func (self *PlantSpawningSystem) HandleTimeSecondPassedMessage(m engo.Message) {
	msg, ok := m.(messages.TimeSecondPassedMessage)
	if !ok {
		return
	}
	for _, e := range self.entities {
		e.Update(msg.Time)
	}
}

func (self *PlantSpawningSystem) HandleControlMessage(m engo.Message) {
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	log.Printf("[HUD] %+v", m)
	switch msg.Action {
	}
}

func (self *PlantSpawningSystem) UpdateSave(saveFile *save.SaveFile) {
	for _, e := range self.entities {
		entityID := e.BasicEntity.ID()
		if _, ok := saveFile.SeenEntityIDs[entityID]; !ok {
			saveFile.Plants = append(saveFile.Plants, e)
			saveFile.SeenEntityIDs[e.BasicEntity.ID()] = struct{}{}
		}
	}
}

func (self *PlantSpawningSystem) LoadSave(saveFile *save.SaveFile) {
	log.Printf("[PlantSpawningSystem] Plants in the save file: %d\n", len(saveFile.Plants))
	for _, c := range saveFile.Plants {
		self.Add(c)
	}
}
