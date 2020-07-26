package systems

import (
	//"fmt"
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	//"image"
	//"image/color"
)

var (
	z_idx_hud       int = 999
)

type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func InitHUD(u engo.Updater) {
	/*
	world, _ := u.(*ecs.World)

	hud := HUD{BasicEntity: ecs.NewBasic()}
	ww, wh := engo.WindowWidth(), engo.WindowHeight()
	hud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, wh - (wh / 2)},
		Width:    ww / 2,
		Height:   wh / 2,
	}
	hudImage := image.NewUniform(color.RGBA{205, 205, 205, 255})
	hudNRGBA := common.ImageToNRGBA(hudImage, int(ww/2), int(wh/2))
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)

	hud.RenderComponent = common.RenderComponent{
		Drawable: hudTexture,
		Scale:    engo.Point{1, 1},
		Repeat:   common.Repeat,
	}
	hud.RenderComponent.SetShader(common.HUDShader)
	hud.RenderComponent.SetZIndex(float32(z_idx_hud))

	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&hud.BasicEntity, &hud.RenderComponent, &hud.SpaceComponent)
		}
	}
	*/
}
