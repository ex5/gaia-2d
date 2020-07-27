package util

import (
	"github.com/EngoEngine/engo"
	"gogame/assets"
)

func ToGridPosition(x float32, y float32) (float32, float32) {
	return float32(int(x/float32(assets.SpriteWidth)) * assets.SpriteWidth),
		float32(int(y/float32(assets.SpriteHeight)) * assets.SpriteHeight)
}

func ToPoint(i int, j int) *engo.Point {
	return &engo.Point{float32(i * assets.SpriteWidth), float32(j * assets.SpriteHeight)}
}
