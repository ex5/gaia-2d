package messages

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	//"log"
)

const ControlMessageType string = "ControlMessage"
const InteractionMessageType string = "InteractionMessage"
const SaveMessageType string = "SaveMessage"
const LoadMessageType string = "LoadMessage"
const TileRemoveMessageType string = "TileRemoveMessage"
const TileReplaceMessageType string = "TileReplaceMessage"
const CreatureHoveredMessageType string = "CreatureHoveredMessage"
const PlantHoveredMessageType string = "PlantHoveredMessage"
const NewPlantMessageType string = "NewPlantMessage"

type ControlMessage struct {
	Action     string
	Data       string
	ObjectID   int
	CreatureID int
}

type InteractionMessage struct {
	Action      string
	BasicEntity *ecs.BasicEntity
}

type SaveMessage struct {
	Filepath string
}

type LoadMessage struct {
	Filepath string
}

type TileRemoveMessage struct {
	Entity *ecs.BasicEntity
}

type TileReplaceMessage struct {
	Entity   *ecs.BasicEntity
	ObjectID int
}

type CreatureHoveredMessage struct {
	EntityID uint64
}

type PlantHoveredMessage struct {
	EntityID uint64
}

type NewPlantMessage struct {
	PlantID int
	Point   *engo.Point
}

func (ControlMessage) Type() string {
	return ControlMessageType
}

func (InteractionMessage) Type() string {
	return InteractionMessageType
}

func (SaveMessage) Type() string {
	return SaveMessageType
}

func (LoadMessage) Type() string {
	return LoadMessageType
}

func (TileRemoveMessage) Type() string {
	return TileRemoveMessageType
}

func (TileReplaceMessage) Type() string {
	return TileReplaceMessageType
}

func (CreatureHoveredMessage) Type() string {
	return CreatureHoveredMessageType
}

func (PlantHoveredMessage) Type() string {
	return PlantHoveredMessageType
}

func (NewPlantMessage) Type() string {
	return NewPlantMessageType
}
