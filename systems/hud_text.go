package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/config"
	"gogame/messages"
	"log"
	"strings"
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

type UIBackground struct {
	ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent

	GetDimensions func() (float32, float32)
	GetPosition   func() *engo.Point
	Color         color.Color
	BorderColor   color.Color

	parent *UIElement
}

// UIElement is an entity for the text system. This keeps track of the position
// size and text associated with that position.
type UIElement struct {
	ecs.BasicEntity

	position    *engo.Point
	Name        string
	GetText     func() string
	GetPosition func() *engo.Point

	hideAfter  time.Duration
	shownSince time.Time
	hidden     bool
	fnt        *common.Font

	bg     *UIBackground
	text   *Text
	height float32
}

// HUDSystem prints the text to our HUD based on the current state of the game
type HUDSystem struct {
	world *ecs.World
	fnt   *common.Font

	entities map[string]*UIElement
}

func (self *UIElement) SetHidden(hide bool) {
	if hide {
		self.hideAfter = 0
	} else {
		self.shownSince = time.Now()
	}
	self.hidden = hide
	self.text.RenderComponent.Hidden = self.hidden
	if self.bg != nil {
		self.bg.RenderComponent.Hidden = self.hidden
	}
	log.Println("SetHidden", hide, self)
}

func (self *HUDSystem) SetText(entityName string, getText func() string, hideAfter time.Duration) {
	entity := self.entities[entityName]
	entity.GetText = getText
	entity.SetHidden(false)
	entity.hideAfter = hideAfter

	// Force update
	entity.Update()
}

func (self *UIElement) Update() {
	if self.GetText == nil {
		return
	}
	newText := self.GetText()
	self.text.RenderComponent.Drawable = common.Text{
		Font: self.fnt,
		Text: newText,
	}
	self.bg.Refresh()
}

func (self *UIElement) Refresh() {
	self.position = self.GetPosition()

	// Initialise text lines, if necessary
	textPosition := engo.Point{
		X: self.position.X + config.HUDTextPadding,
		Y: self.position.Y + config.HUDTextPadding,
	}
	if self.text == nil {
		self.text = &Text{BasicEntity: ecs.NewBasic()}
		self.text.SetShader(common.TextHUDShader)
		self.text.RenderComponent.Drawable = common.Text{
			Font: self.fnt,
			Text: "",
		}
		self.text.RenderComponent.SetZIndex(config.HUDLayer)
		self.text.SpaceComponent = common.SpaceComponent{
			Position: textPosition,
		}
		// Parenting is useless currently: it doesn't imply coordinate transform parenting FIXME 2.0
		// self.BasicEntity.AppendChild(&text.BasicEntity)
	}
	self.text.SpaceComponent.Position = textPosition

	self.bg.Refresh()
}

func (self *UIBackground) Refresh() {
	if self == nil {
		return
	}
	width, height := self.GetDimensions()
	position := *self.GetPosition()
	log.Println("UIBackground Refresh", width, height, position)
	if self.SpaceComponent == nil {
		self.SpaceComponent = &common.SpaceComponent{
			Position: position,
			Width:    width,
			Height:   height,
		}
	}
	if self.RenderComponent == nil {
		self.RenderComponent = &common.RenderComponent{
			Drawable: common.Rectangle{
				BorderWidth: 1,
				BorderColor: self.BorderColor,
			},
			Color: self.Color,
		}
		self.RenderComponent.SetZIndex(config.HUDLayer - 1)
		self.RenderComponent.SetShader(common.HUDShader)
	}
	self.SpaceComponent.Position = position
	self.SpaceComponent.Width = width
	self.SpaceComponent.Height = height
}

// New is called when the system is added to the world.
// Adds text to our HUD that will update based on the state of the game, then
// listens for messages to update the text.
func (self *HUDSystem) New(w *ecs.World) {
	self.world = w
	self.fnt = &common.Font{
		URL:  config.FontURL,
		FG:   color.White,
		Size: config.FontSize,
	}
	self.fnt.CreatePreloaded()

	self.entities = make(map[string]*UIElement)
	// Initialise all known text elements of the UI
	self.NewUIElement("HoverInfo", func() *engo.Point {
		return &engo.Point{config.HUDMarginL, engo.WindowHeight() - config.HoverInfoHeight}
	}, 8, &UIBackground{
		Color:       color.RGBA{0, 0, 0, 150},
		BorderColor: color.RGBA{50, 50, 50, 255},
	})

	self.NewUIElement("EventMessage", func() *engo.Point {
		return &engo.Point{config.HUDMarginL, config.HUDMarginT}
	}, 1, nil)

	self.NewUIElement("CurrentTime", func() *engo.Point {
		return &engo.Point{engo.WindowWidth() - float32(config.FontSize*20), engo.WindowHeight() - config.HoverInfoHeight}
	}, 2, nil)

	self.NewUIElement("Overlay", func() *engo.Point {
		return &engo.Point{0, 0}
	}, 0, &UIBackground{
		GetDimensions: func() (float32, float32) {
			return engo.WindowWidth(), engo.WindowHeight()
		},
		Color:       color.RGBA{0, 0, 0, 100},
		BorderColor: color.RGBA{0, 0, 0, 255},
	})
	self.entities["Overlay"].SetHidden(true)

	// Messages set the texts of the text UI elements
	engo.Mailbox.Listen(messages.HUDTextUpdateMessageType, self.HandleHUDTextUpdateMessage)
	engo.Mailbox.Listen(messages.TimeSecondPassedMessageType, self.HandleTimeSecondPassedMessage)
	engo.Mailbox.Listen("WindowResizeMessage", self.HandleWindowResizeMessage)
	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
}

// Add adds an entity to the system
func (self *HUDSystem) NewUIElement(name string, getPosition func() *engo.Point, h int, bg *UIBackground) {
	entity := &UIElement{
		BasicEntity: ecs.NewBasic(),
		GetPosition: getPosition,
		GetText:     func() string { return "" },
		position:    getPosition(),
		Name:        name,
		fnt:         self.fnt,
		bg:          bg,
	}

	if entity.bg != nil {
		entity.bg.GetPosition = getPosition
		// IF not GetDimensions is given, set dimensions based on the text
		if entity.bg.GetDimensions == nil {
			entity.bg.GetDimensions = func() (float32, float32) {
				text := entity.GetText()
				_, h, b := entity.fnt.TextDimensions(text)
				maxW := 0
				for _, l := range strings.Split(text, "\n") {
					width, _, _ := entity.fnt.TextDimensions(l)
					if width > maxW {
						maxW = width
					}
				}
				return float32(maxW) + 3*config.HUDTextPadding, float32(h * b)
			}
		}
	}

	entity.Refresh()

	self.Add(entity)
}

func (self *HUDSystem) Add(entity *UIElement) {
	if entity.bg != nil {
		for _, system := range self.world.Systems() {
			switch sys := system.(type) {
			case *common.RenderSystem:
				sys.Add(&entity.bg.BasicEntity, entity.bg.RenderComponent, entity.bg.SpaceComponent)
			}
		}
	}
	for _, system := range self.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&entity.text.BasicEntity, &entity.text.RenderComponent, &entity.text.SpaceComponent)
		}
	}

	self.entities[entity.Name] = entity
}

func (self *HUDSystem) HandleHUDTextUpdateMessage(m engo.Message) {
	msg, ok := m.(messages.HUDTextUpdateMessage)
	if !ok {
		return
	}
	self.SetText(msg.Name, msg.GetText, msg.HideAfter)
}

func (self *HUDSystem) HandleTimeSecondPassedMessage(m engo.Message) {
	msg, ok := m.(messages.TimeSecondPassedMessage)
	if !ok {
		return
	}
	// Display current date and time
	self.SetText("CurrentTime", msg.Time.GetTextStatus, 0)

	// Update other UI elements
	for _, e := range self.entities {
		e.Update()
	}
}

func (self *HUDSystem) HandleWindowResizeMessage(m engo.Message) {
	_, ok := m.(engo.WindowResizeMessage)
	if !ok {
		return
	}

	// Resize UI elements
	for _, e := range self.entities {
		e.Refresh()
	}
}

func (self *HUDSystem) HandleControlMessage(m engo.Message) {
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	log.Printf("[HUD] %+v", m)
	switch msg.Action {
	case "TogglePause":
		self.entities["Overlay"].SetHidden(!self.entities["Overlay"].hidden)
	}
}

// Update is called each frame to update the system.
func (self *HUDSystem) Update(dt float32) {}

// Remove takes an enitty out of the system.
func (self *HUDSystem) Remove(basic ecs.BasicEntity) {
	// TODO Remove all `Text`s from all relevant systems
}
