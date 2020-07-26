package util

import (
	"gogame/assets"
)


func ToGridPosition(x float32, y float32) (float32, float32) {
	return float32(int(x / float32(assets.SpriteWidth)) * assets.SpriteWidth),
		float32(int(y / float32(assets.SpriteHeight)) * assets.SpriteHeight)
}
