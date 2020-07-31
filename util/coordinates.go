package util

import (
	"github.com/EngoEngine/engo"
	"gogame/config"
)

func ToGridPosition(x float32, y float32) (float32, float32) {
	return float32(int(x/float32(config.SpriteWidth)) * config.SpriteWidth),
		float32(int(y/float32(config.SpriteHeight)) * config.SpriteHeight)
}

func ToPoint(i int, j int) *engo.Point {
	return &engo.Point{float32(i * config.SpriteWidth), float32(j * config.SpriteHeight)}
}
