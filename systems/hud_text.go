package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/config"
	"gogame/messages"
	"log"
	"time"
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
	ecs.BasicEntity

	Position      *engo.Point
	Name          string
	Height        int
	GetTextStatus func() []string

	hideAfter  time.Duration
	shownSince time.Time
	hidden     bool
	fnt        *common.Font

	text []*Text
}

// HUDTextSystem prints the text to our HUD based on the current state of the game
type HUDTextSystem struct {
	world *ecs.World
	fnt   *common.Font

	dtFullSeconds float32
	entities      map[string]*HUDTextEntity
}

func (self *HUDTextEntity) SetHidden(hide bool) {
	if hide {
		self.hideAfter = 0
	} else {
		self.shownSince = time.Now()
	}
	self.hidden = hide
	for i := 0; i < len(self.text); i++ {
		self.text[i].RenderComponent.Hidden = self.hidden
	}
	log.Println("SetHidden", hide, self)
}

func (h *HUDTextSystem) SetText(entityName string, display func() []string, hideAfter time.Duration) {
	entity := h.entities[entityName]
	entity.GetTextStatus = display
	entity.SetHidden(false)
	entity.hideAfter = hideAfter

	// Force update
	entity.Update()
}

func (self *HUDTextEntity) Update() {
	if self.GetTextStatus == nil {
		return
	}
	lines := self.GetTextStatus()
	for i := 0; i < len(lines) && i < self.Height; i++ {
		self.text[i].RenderComponent.Drawable = common.Text{
			Font: self.fnt,
			Text: lines[i],
		}
	}
}

func (h *HUDTextSystem) InitText(entityName string) {
	entity := h.entities[entityName]
	entity.fnt = h.fnt

	for i := 0; i < entity.Height; i++ {
		text := Text{BasicEntity: ecs.NewBasic()}
		text.SetShader(common.TextHUDShader)
		text.RenderComponent.Drawable = common.Text{
			Font: h.fnt,
			Text: "",
		}
		text.RenderComponent.SetZIndex(config.HUDLayer)
		text.SpaceComponent = common.SpaceComponent{
			Position: engo.Point{
				X: entity.Position.X,
				Y: entity.Position.Y + float32(i*config.LineHeight),
			},
		}
		for _, system := range h.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&text.BasicEntity, &text.RenderComponent, &text.SpaceComponent)
			}
		}
		entity.text = append(entity.text, &text)
	}
}

// New is called when the system is added to the world.
// Adds text to our HUD that will update based on the state of the game, then
// listens for messages to update the text.
func (h *HUDTextSystem) New(w *ecs.World) {
	h.world = w
	h.fnt = &common.Font{
		URL:  config.FontURL,
		FG:   color.White,
		BG:   color.Black,
		Size: config.FontSize,
	}
	h.fnt.CreatePreloaded()

	// Initialise all known text elements of the UI
	h.entities = make(map[string]*HUDTextEntity)
	h.Add("HoverInfo", config.HUDMarginL, engo.WindowHeight()-config.HoverInfoHeight, 8)
	h.Add("EventMessage", config.HUDMarginL, config.HUDMarginT, 1)

	// Messages set the texts of the text UI elements
	engo.Mailbox.Listen(messages.HUDTextUpdateMessageType, h.HandleHUDTextUpdateMessage)
}

func (h *HUDTextSystem) HandleHUDTextUpdateMessage(m engo.Message) {
	msg, ok := m.(messages.HUDTextUpdateMessage)
	if !ok {
		return
	}
	h.SetText(msg.Name, msg.GetText, msg.HideAfter)
}

// Add adds an entity to the system
func (self *HUDTextSystem) Add(name string, x float32, y float32, h int) {
	position := &engo.Point{X: x, Y: y}
	entity := &HUDTextEntity{BasicEntity: ecs.NewBasic(), Position: position, Name: name, Height: h}
	self.entities[entity.Name] = entity

	self.InitText(name)
}

// Update is called each frame to update the system.
func (self *HUDTextSystem) Update(dt float32) {
	for _, e := range self.entities {
		if self.dtFullSeconds > 1 {
			e.Update()
		}

		if e.hideAfter > 0 && !e.hidden {
			now := time.Now()
			if now.Sub(e.shownSince) > e.hideAfter {
				e.SetHidden(true)
			}
		}
	}
	if self.dtFullSeconds > 1 {
		self.dtFullSeconds = 0
	}
	self.dtFullSeconds += dt
}

// Remove takes an enitty out of the system.
func (h *HUDTextSystem) Remove(basic ecs.BasicEntity) {
	// TODO Remove all `Text`s from all relevant systems
}
