package save

import (
	"gogame/data"
	"gogame/life/plants"
)

type SaveFile struct {
	Tiles     []*data.Tile     `json:"tiles"`
	Creatures []*data.Creature `json:"creatures"`
	Plants    []*plants.Plant  `json:"plants"`

	SeenEntityIDs map[uint64]struct{} `json:"-"`
}
