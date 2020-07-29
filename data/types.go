package data

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
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
	*Tile

	ID             int       `json:"id"`
	ObjectID       int       `json:"object_id"`
	IsAlive        bool      `json:"is_alive"`
	Name           string    `json:"name"`
	Subspecies     string    `json:"subspecies"`
	Food           float32   `json:food`
	MinFood        float32   `json:"min_food"`
	MaxFood        float32   `json:"max_food"`
	Sleep          float32   `json:"sleep"`
	MinSleep       float32   `json:"min_sleep"`
	MaxSleep       float32   `json:"max_sleep"`
	Speed          float32   `json:"speed"`
	Needs          []*string `json:"needs"`
	MovementTarget *Tile     `json:"movement_target"`
}

type Creatures struct {
	Creatures []*Creature `json:"creatures"`
}

type SaveFile struct {
	Tiles     []*Tile     `json:"tiles"`
	Creatures []*Creature `json:"creatures"`
}

func (self Creature) Update(dt float32) {
	//log.Printf("%+v %+v", dt, self.SpaceComponent)
	self.Tile.SpaceComponent.Position.X += dt * 100
}
