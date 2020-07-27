package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/assets"
	"gogame/messages"
)

type mouseTracker struct {
	ecs.BasicEntity
	common.MouseComponent
}

// Text is an entity containing text printed to the screen
type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

// HUDTextEntity is an entity for the text system. This keeps track of the position
// size and text associated with that position.
type HUDTextEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.MouseComponent
	Line1, Line2, Line3, Line4 string
}

// HUDTextSystem prints the text to our HUD based on the current state of the game
type HUDTextSystem struct {
	world *ecs.World
	fnt   *common.Font

	text1, text2, text3, text4 *Text

	entities []HUDTextEntity
}

func (h *HUDTextSystem) SetText(text []string) {
	h.text1.RenderComponent.Drawable = common.Text{
		Font: h.fnt,
		Text: text[0],
	}
	h.text2.RenderComponent.Drawable = common.Text{
		Font: h.fnt,
		Text: text[1],
	}
	h.text3.RenderComponent.Drawable = common.Text{
		Font: h.fnt,
		Text: text[2],
	}
	h.text4.RenderComponent.Drawable = common.Text{
		Font: h.fnt,
		Text: text[3],
	}
}

func (h *HUDTextSystem) InitText(lineNo int) *Text {
	text := Text{BasicEntity: ecs.NewBasic()}
	text.SetShader(common.TextHUDShader)
	text.RenderComponent.Drawable = common.Text{
		Font: h.fnt,
		Text: "",
	}
	text.RenderComponent.SetZIndex(assets.HUDLayer)
	text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: assets.HUDMarginL, Y: engo.WindowHeight() - (assets.HoverTipHeight - float32(lineNo*assets.LineHeight))},
	}
	for _, system := range h.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
		}
	}
	return &text
}

// New is called when the system is added to the world.
// Adds text to our HUD that will update based on the state of the game, then
// listens for messages to update the text.
func (h *HUDTextSystem) New(w *ecs.World) {
	h.world = w
	h.fnt = &common.Font{
		URL:  assets.FontURL,
		FG:   color.White,
		BG:   color.Black,
		Size: assets.FontSize,
	}
	h.fnt.CreatePreloaded()

	h.text1 = h.InitText(0)
	h.text2 = h.InitText(1)
	h.text3 = h.InitText(2)
	h.text4 = h.InitText(3)

	engo.Mailbox.Listen(messages.InteractionMessageType, func(m engo.Message) {
		msg, ok := m.(messages.InteractionMessage)
		if !ok {
			return
		}
		if msg.BasicEntity == nil {
			h.text1.RenderComponent.Drawable = common.Text{
				Font: h.fnt,
				Text: "",
			}
			h.text2.RenderComponent.Drawable = common.Text{
				Font: h.fnt,
				Text: "",
			}
			h.text3.RenderComponent.Drawable = common.Text{
				Font: h.fnt,
				Text: "",
			}
			h.text4.RenderComponent.Drawable = common.Text{
				Font: h.fnt,
				Text: "",
			}
		}
	})

	engo.Mailbox.Listen(messages.HUDTextMessageType, func(m engo.Message) {
		msg, ok := m.(messages.HUDTextMessage)
		if !ok {
			return
		}
		h.text1.RenderComponent.Drawable = common.Text{
			Font: h.fnt,
			Text: msg.Line1,
		}
		h.text2.RenderComponent.Drawable = common.Text{
			Font: h.fnt,
			Text: msg.Line2,
		}
		h.text3.RenderComponent.Drawable = common.Text{
			Font: h.fnt,
			Text: msg.Line3,
		}
		h.text4.RenderComponent.Drawable = common.Text{
			Font: h.fnt,
			Text: msg.Line4,
		}
	})
}

// Add adds an entity to the system
func (h *HUDTextSystem) Add(b *ecs.BasicEntity, s *common.SpaceComponent, m *common.MouseComponent, l1, l2, l3, l4 string) {
	h.entities = append(h.entities, HUDTextEntity{b, s, m, l1, l2, l3, l4})
}

// Update is called each frame to update the system.
func (h *HUDTextSystem) Update(dt float32) {}

// Remove takes an enitty out of the system.
func (h *HUDTextSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range h.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		h.entities = append(h.entities[:delete], h.entities[delete+1:]...)
	}
}
