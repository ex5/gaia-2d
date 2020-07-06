package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"

    //"log"
    "fmt"
)

type AnimalMouseTracker struct {
    ecs.BasicEntity
    common.MouseComponent
}

type Animal struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// System is an interface which implements an ECS-System. A System
// should iterate over its Entities on `Update`, in any way
// suitable for the current implementation.
type System interface {
	// Update is ran every frame, with `dt` being the time
	// in seconds since the last frame
	Update(dt float32)

	// Remove removes an Entity from the System
	Remove(ecs.BasicEntity)
}

type AnimalSpawningSystem struct {
    world *ecs.World
    mouseTracker AnimalMouseTracker
}

// New is the initialisation of the System
func (self *AnimalSpawningSystem) New(w *ecs.World) {
	fmt.Println("AnimalSpawningSystem was added to the Scene")

    self.world = w

    self.mouseTracker.BasicEntity = ecs.NewBasic()
	self.mouseTracker.MouseComponent = common.MouseComponent{Track: true}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&self.mouseTracker.BasicEntity, &self.mouseTracker.MouseComponent, nil, nil)
		}
	}
}

// Update is ran every frame, with `dt` being the time
// in seconds since the last frame
func (self *AnimalSpawningSystem) Update(dt float32) {
    if engo.Input.Button("AddAnimal").JustPressed()  {
		fmt.Println("The gamer pressed F1")
        animal := Animal{BasicEntity: ecs.NewBasic()}

		animal.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{self.mouseTracker.MouseX, self.mouseTracker.MouseY},
			Width:    30,
			Height:   64,
		}

		texture, err := common.LoadedSprite("textures/chick_24x24.png")
		if err != nil {
			panic("Unable to load texture: " + err.Error())
		}

		animal.RenderComponent = common.RenderComponent{
			Drawable: texture,
			Scale:    engo.Point{X: 1, Y: 1},
		}
		animal.RenderComponent.SetZIndex(10.0)

		for _, system := range self.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&animal.BasicEntity, &animal.RenderComponent, &animal.SpaceComponent)
			}
		}
	}
}

// Remove is called whenever an Entity is removed from the World, in order to remove it from this sytem as well
func (*AnimalSpawningSystem) Remove(ecs.BasicEntity) {}
