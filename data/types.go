package data

import (
	"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/config"
	"gogame/messages"
	"gogame/util"
	"log"
)

type Spritesheet struct {
	ID         int                 `json:"id"`
	FilePath   string              `json:"filepath"`
	Animations []*common.Animation `json:"animations"`
}

type Spritesheets struct {
	Spritesheets []*Spritesheet              `json:"spritesheets"`
	Loaded       map[int]*common.Spritesheet `json:"-"`
}

// Description of a resource
type Resource struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type Resources struct {
	Resources []*Resource `json:"resources"`
}

// Mutable resource, e.g. plant matter available at a specific tile
type AccessibleResource struct {
	ResourceID int     `json:"resource_id"`
	Amount     float32 `json:"amount"`
}

type Object struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // tile|creature: to distinguish during [un]marshalling
	SpriteID      int     `json:"sprite_id"`
	SpritesheetID int     `json:"spritesheet_id"`
	Name          string  `json:"name"`
	ResourceID    int     `json:"resource_id"`
	Amount        float32 `json:"amount"`

	// Runtime only fields
	Spritesheet *common.Spritesheet `json:"-"`
	Animations  []*common.Animation `json:"-"`
}

type Objects struct {
	Objects []*Object `json:"objects"`
}

type Tile struct {
	*ecs.BasicEntity           `json:"-"` // FIXME? marshalled into an empty object
	*common.RenderComponent    `json:"-"` // FIXME cannot unmarshal .. color.Color
	*common.AnimationComponent `json:"-"`

	SpaceComponent     *common.SpaceComponent
	CollisionComponent *common.CollisionComponent
	MouseComponent     *common.MouseComponent

	Layer              float32
	ObjectID           int
	AccessibleResource *AccessibleResource
	Object             *Object   `json:"-"`
	Resource           *Resource `json:"-"`
}

type Creature struct {
	*Tile `deepcopier:"skip"`

	ID             int      `json:"id"`
	ObjectID       int      `json:"object_id"`
	IsAlive        bool     `json:"is_alive"`
	Name           string   `json:"name"`
	Subspecies     string   `json:"subspecies"`
	Food           float32  `json:"food"`
	Eats           []int    `json:"eats"` // Resource IDs
	MinFood        float32  `json:"min_food"`
	MaxFood        float32  `json:"max_food"`
	Sleep          float32  `json:"sleep"`
	MinSleep       float32  `json:"min_sleep"`
	MaxSleep       float32  `json:"max_sleep"`
	MovementSpeed  float32  `json:"movement_speed"`
	EatingSpeed    float32  `json:"eating_speed"`
	Needs          []string `json:"needs"`
	Activity       string   `json:"activity"` // idle, eating, looking_for_food, sleeping
	LastEventID    uint64   `json:"last_event_id"`
	MovementTarget *Tile    `json:"movement_target"`
	Target         *Tile    `json:"target"`
}

type Creatures struct {
	Creatures []*Creature `json:"creatures"`
}

func (self *Creature) FindFood(x engo.AABBer) bool {
	if tile, ok := x.(*Tile); ok {
		fmt.Println("Checking", tile, self.Eats, tile.Resource.ID)
		if util.ContainsInt(self.Eats, tile.Resource.ID) && tile.AccessibleResource.Amount > 0 {
			return true
		}
	}
	return false
}

func (self *Creature) IsHungry() bool {
	return self.Food < self.MinFood
}

func (self *Creature) IsSatiated() bool {
	return self.Food >= self.MaxFood
}

func (self *Creature) BecomeIdle() {
	self.Target = nil
	self.Activity = ""
}

func (self *Creature) DecideToWander() bool {
	return util.Roll(0.3, 0.9) > 0.5
}

type SaveFile struct {
	Tiles     []*Tile     `json:"tiles"`
	Creatures []*Creature `json:"creatures"`
}

func (self *Creature) Update(dt float32) {
	//log.Printf("%+v %+v", dt, self.SpaceComponent)
	// Expend them calories TODO moving increases, sleeping decreases
	if self.Activity != "eating" {
		self.Food -= self.EatingSpeed * dt
	}
	// Handle hunger
	log.Println(self.IsHungry(), self.Food, self.MaxFood, self.MinFood, self.EatingSpeed, self.Needs, self.HasNeedFor("food"))
	if self.IsHungry() && !self.HasNeedFor("food") {
		self.AddNeedFor("food")
		log.Println(self, "needs food!", self.Needs, len(self.Needs) > 1 && self.Needs[0] == "food", self.Needs)
	}
	if len(self.Needs) > 0 && self.Needs[0] == "food" && self.Activity != "looking_for_food" {
		self.Activity = "looking_for_food"
		log.Println(self, "looks for food", self.Needs)
		engo.Mailbox.Dispatch(messages.SpacialRequestMessage{
			Aabb:     self.Tile.SpaceComponent.AABB(),
			Filter:   self.FindFood,
			EntityID: self.BasicEntity.ID(),
			EventID:  self.LastEventID + 1,
		})
		self.LastEventID++
	}
	// Handle movement
	// TODO write the SpeedComponent and remove movement logic from here completely.
	if self.MovementTarget != nil {
		if self.TooFar(self.MovementTarget, 0.1) {
			v := self.Direction(self.MovementTarget)
			log.Println("Need to move", v)
			self.Tile.SpaceComponent.Position.Add(*v.MultiplyScalar(dt))
			// self.Tile.SpaceComponent.Position.Add(*v.MultiplyScalar(dt))
			// TODO speed and smooth movement
		} else {
			if self.Activity == "looking_for_food" {
				if self.FindFood(self.MovementTarget) {
					self.Activity = "eating"
					self.Target = self.MovementTarget
				} else {
					self.BecomeIdle()
				}
			}
			self.MovementTarget = nil
		}
	}

	// Handle eating
	if self.Activity == "eating" && self.Target != nil {
		if !self.IsSatiated() {
			target := *self.Target
			if target.AccessibleResource.Amount > 0 {
				eaten := self.EatingSpeed * dt
				self.Food += eaten
				self.Target.AccessibleResource.Amount -= eaten
				log.Println(self, "eating", dt, eaten, self.Food)
			}
			if target.AccessibleResource.Amount <= 0 {
				// Ate it all
				log.Println("Ate all of", self.Target)
				engo.Mailbox.Dispatch(messages.TileRemoveMessage{
					Entity: self.Target.BasicEntity,
				})
				self.BecomeIdle()
				// panic("eat all")
			}
		} else {
			self.BecomeIdle()
			self.RemoveNeedFor("food")
			//panic("not hungry")
		}
	}

	// Handle idling
	if self.Activity == "" && len(self.Needs) == 0 {
		if self.Activity != "wandering" && self.DecideToWander() {
			self.Activity = "wandering"
			engo.Mailbox.Dispatch(messages.SpacialRequestMessage{
				Aabb:     self.SurroundingAreaAABB(100),
				Filter:   func(aabb engo.AABBer) bool { return true },
				EntityID: self.BasicEntity.ID(),
				EventID:  self.LastEventID + 1,
			})
			self.LastEventID++
		}
	}
}

func (self *Tile) AABB() engo.AABB {
	return self.SpaceComponent.AABB()
}

func (self *Tile) SurroundingAreaAABB(radius float32) engo.AABB {
	return engo.AABB{
		Min: engo.Point{
			X: self.SpaceComponent.Position.X - radius,
			Y: self.SpaceComponent.Position.Y - radius,
		},
		Max: engo.Point{
			X: self.SpaceComponent.Position.X + float32(config.SpriteWidth) + radius,
			Y: self.SpaceComponent.Position.Y + float32(config.SpriteHeight) + radius,
		},
	}
}

func (self *Creature) AddNeedFor(smthg string) bool {
	for _, n := range self.Needs {
		if n == smthg {
			return false
		}
	}
	self.Needs = append(self.Needs, smthg)
	return true
}

func (self *Creature) HasNeedFor(smthg string) bool {
	return util.ContainsStr(self.Needs, smthg)
}

func (self *Creature) RemoveNeedFor(smthg string) bool {
	delete := -1
	for index, n := range self.Needs {
		if n == smthg {
			delete = index
		}
	}
	if delete >= 0 {
		self.Needs = append(self.Needs[:delete], self.Needs[delete+1:]...)
		return true
	}
	return false
}

func (self *Creature) TooFar(tile *Tile, dt float32) bool {
	return self.Tile.SpaceComponent.Position.PointDistance(tile.SpaceComponent.Position) > dt
}

func (self *Creature) Direction(tile *Tile) *engo.Point {
	v := engo.Point{0, 0}
	return v.Add(
		tile.SpaceComponent.Position).Subtract(self.Tile.SpaceComponent.Position)
}
