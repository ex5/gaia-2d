package data

import (
	"fmt"
	"github.com/EngoEngine/engo"
	"gogame/calendar"
	"gogame/messages"
	"gogame/util"
	"log"
	"strings"
	"time"
)

type Activity uint8
type Want uint8
type Need struct {
	Duration time.Duration
	Want     Want
}

const (
	Idle Activity = iota
	LookingAround
	Eating
	Wandering
	Sleeping
)

const (
	Food Want = iota
	Sleep
)

func (a Activity) String() string {
	return [...]string{"idle", "looking around", "eating", "wandering"}[a]
}

func (w Want) String() string {
	return [...]string{"food", "sleep"}[w]
}

type Creature struct {
	*Tile `deepcopier:"skip"`

	// Species properties, immutable
	ID            int     `json:"id"`
	ObjectID      int     `json:"object_id"`
	EatingSpeed   float32 `json:"eating_speed"`
	Eats          []int   `json:"eats"` // Resource IDs
	MaxFood       float32 `json:"max_food"`
	MaxSleep      float32 `json:"max_sleep"`
	MinFood       float32 `json:"min_food"`
	MinSleep      float32 `json:"min_sleep"`
	MovementSpeed float32 `json:"movement_speed"`
	Species       string  `json:"species"`

	// Live properties, mutable
	Activity       Activity `json:"activity"`
	Food           float32  `json:"food"`
	IsAlive        bool     `json:"is_alive"`
	MovementTarget *Tile    `json:"movement_target"`
	Name           string   `json:"name"`
	Needs          []*Need  `json:"needs"`
	Sleep          float32  `json:"sleep"`
	Target         *Tile    `json:"target"`

	LastEventID uint64 `json:"last_event_id"`
}

type Creatures struct {
	Creatures []*Creature `json:"creatures"`
}

func (self *Creature) FindFood(x engo.AABBer) bool {
	if tile, ok := x.(*Tile); ok {
		fmt.Println("Checking", tile, self.Eats, "resource:", tile.Resource)
		if tile.Resource != nil && util.ContainsInt(self.Eats, tile.Resource.ID) && tile.AccessibleResource.Amount > 0 {
			return true
		}
	}
	return false
}

func (self *Creature) IsHungry() bool {
	return self.Food < self.MinFood
}

func (self *Creature) IsTired() bool {
	return self.Sleep < self.MinSleep
}

func (self *Creature) IsSatiated() bool {
	return self.Food >= self.MaxFood
}

func (self *Creature) IsFullyRested() bool {
	return self.Sleep >= self.MaxSleep
}

func (self *Creature) BecomeIdle() {
	self.Target = nil
	self.Activity = Idle
}

func (self *Creature) DecideToWander() bool {
	return util.Roll(0.3, 0.9) > 0.5
}

func (self *Creature) Update(dt float32) {
	// Handle movement
	// TODO write the SpeedComponent and remove movement logic from here completely.
	if self.MovementTarget != nil {
		if self.TooFar(self.MovementTarget, 0.1) {
			v := self.Direction(self.MovementTarget)
			log.Println("Moving", v)
			self.Tile.SpaceComponent.Position.Add(*v.MultiplyScalar(dt))
			// TODO speed and smooth movement
		} else {
			log.Println(self, "reached the movement target", self.MovementTarget)
			self.BecomeIdle()
			self.Target = self.MovementTarget
			self.MovementTarget = nil
		}
	}
}

func (self *Creature) UpdateActivity(currentTime *calendar.Time) {
	// Handle durations of needs
	for _, n := range self.Needs {
		log.Println(n, n.Duration, time.Duration(int64(time.Second)))
		n.Duration += time.Duration(int64(time.Second))
	}
	if self.Activity != Eating {
		// Expend them calories TODO moving increases, sleeping decreases
		self.Food -= self.EatingSpeed
	}
	// Handle hunger
	if self.IsHungry() && !self.HasNeedFor(Food) {
		self.AddNeedFor(Food)
		log.Println(self, "needs food!", self.Needs)
	}
	if len(self.Needs) > 0 && self.Needs[0].Want == Food && self.Activity != LookingAround {
		self.Activity = LookingAround
		log.Println(self, "looks around", self.Needs)
		engo.Mailbox.Dispatch(messages.SpacialRequestMessage{
			Aabb:     self.SurroundingAreaAABB(2),
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
			log.Println("Moving", v)
			self.Tile.SpaceComponent.Position.Add(*v)
			// TODO speed and smooth movement
		} else {
			if self.Activity == LookingAround {
				if self.FindFood(self.MovementTarget) {
					log.Println(self, "got to the food!")
					self.Activity = Eating
					self.Target = self.MovementTarget
				} else {
					self.BecomeIdle()
				}
			}
			self.MovementTarget = nil
		}
	}

	// Handle eating
	if self.Activity == Eating && self.Target != nil {
		if !self.IsSatiated() {
			target := *self.Target
			if target.AccessibleResource.Amount > 0 {
				eaten := self.EatingSpeed
				self.Food += eaten
				self.Target.AccessibleResource.Amount -= eaten
				log.Println(self, "eating", eaten, self.Food)
			}
			if target.AccessibleResource.Amount <= 0 {
				// Ate it all
				log.Println("Ate all of", self.Target)
				engo.Mailbox.Dispatch(messages.TileRemoveMessage{
					Entity: self.Target.BasicEntity,
				})
				self.BecomeIdle()
			}
		} else {
			self.BecomeIdle()
			self.RemoveNeedFor(Food)
		}
	}

	// Handle idling
	if self.Activity == Idle && len(self.Needs) == 0 {
		if self.Activity != Wandering && self.DecideToWander() {
			self.Activity = Wandering
			engo.Mailbox.Dispatch(messages.SpacialRequestMessage{
				Aabb:     self.SurroundingAreaAABB(5),
				Filter:   func(aabb engo.AABBer) bool { return true },
				EntityID: self.BasicEntity.ID(),
				EventID:  self.LastEventID + 1,
			})
			self.LastEventID++
		}
	}
}

func (self *Creature) HasNeedFor(want Want) bool {
	for _, n := range self.Needs {
		if n.Want == want {
			return true
		}
	}
	return false
}

func (self *Creature) AddNeedFor(want Want) bool {
	if self.HasNeedFor(want) {
		return false
	}
	n := Need{0, want}
	self.Needs = append(self.Needs, &n)
	return true
}

func (self *Creature) RemoveNeedFor(want Want) bool {
	delete := -1
	for index, n := range self.Needs {
		if n.Want == want {
			delete = index
		}
	}
	if delete >= 0 {
		self.Needs = append(self.Needs[:delete], self.Needs[delete+1:]...)
		return true
	}
	return false
}

func (self *Creature) CurrentNeeds() string {
	var needs []string
	for _, n := range self.Needs {
		needs = append(needs, fmt.Sprintf("%s (%s)", n.Want, util.FormatDuration(n.Duration)))
	}
	return strings.Join(needs, ",")
}

func (self *Creature) CurrentHealth() string {
	return fmt.Sprintf(
		"Food: %d/%d, sleep: %d/%d",
		int(self.Food), int(self.MaxFood), int(self.Sleep), int(self.MaxSleep),
	)
}

func (self *Creature) GetTextStatus() string {
	return fmt.Sprintf(
		"%s, %s \nNeeds %s \n%s \n%s \n%s \n",
		self.Name, self.Species,
		self.CurrentNeeds(),
		self.CurrentHealth(),
		self.Activity,
		self.CurrentPosition(),
	)
}

func (self *Creature) TooFar(tile *Tile, dt float32) bool {
	return self.Tile.SpaceComponent.Position.PointDistance(tile.SpaceComponent.Position) > dt
}

func (self *Creature) Direction(tile *Tile) *engo.Point {
	v := engo.Point{0, 0}
	return v.Add(
		tile.SpaceComponent.Position).Subtract(self.Tile.SpaceComponent.Position)
}
