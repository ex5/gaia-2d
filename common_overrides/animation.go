/* FIXME This is almost exact copy of common.AnimationSystem.
The differences are: a global adjustable speed and an ability to stop `Update`s temporarily.
The speed is used to speed up animation playback when game speed is changed, and stopping `Update`s
allows pausing animations playback when the game is paused.

The only method omitted from the copypasta is AddByInterface because it drags along
the whole pile of interfaces outside the scope of this feature.
*/
package common_overrides

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"gogame/messages"
	"log"
)

// AnimationComponent tracks animations of an entity it is part of.
// This component should be created using NewAnimationComponent.
type AnimationComponent struct {
	Drawables        []common.Drawable            // Renderables
	Animations       map[string]*common.Animation // All possible animations
	CurrentAnimation *common.Animation            // The current animation
	Rate             float32                      // How often frames should increment, in seconds.
	index            int                          // What frame in the is being used
	change           float32                      // The time since the last incrementation
	def              *common.Animation            // The default animation to play when nothing else is playing
}

// NewAnimationComponent creates an AnimationComponent containing all given
// drawables. Animations will be played using the given rate.
func NewAnimationComponent(drawables []common.Drawable, rate float32) AnimationComponent {
	return AnimationComponent{
		Animations: make(map[string]*common.Animation),
		Drawables:  drawables,
		Rate:       rate,
	}
}

// SelectAnimationByName sets the current animation. The name must be
// registered.
func (ac *AnimationComponent) SelectAnimationByName(name string) {
	ac.CurrentAnimation = ac.Animations[name]
	ac.index = 0
}

// SelectAnimationByAction sets the current animation.
// An nil action value selects the default animation.
func (ac *AnimationComponent) SelectAnimationByAction(action *common.Animation) {
	ac.CurrentAnimation = action
	ac.index = 0
}

// AddDefaultAnimation adds an animation which is used when no other animation is playing.
func (ac *AnimationComponent) AddDefaultAnimation(action *common.Animation) {
	ac.AddAnimation(action)
	ac.def = action
}

// AddAnimation registers an animation under its name, making it available
// through SelectAnimationByName.
func (ac *AnimationComponent) AddAnimation(action *common.Animation) {
	ac.Animations[action.Name] = action
}

// AddAnimations registers all given animations.
func (ac *AnimationComponent) AddAnimations(actions []*common.Animation) {
	for _, action := range actions {
		ac.AddAnimation(action)
	}
}

// Cell returns the drawable for the current frame.
func (ac *AnimationComponent) Cell() common.Drawable {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation. Selecting zeroth drawable. If this is incorrect, add an action to the animation.")
		return ac.Drawables[0]
	}
	idx := ac.CurrentAnimation.Frames[ac.index]

	return ac.Drawables[idx]
}

// NextFrame advances the current animation by one frame.
func (ac *AnimationComponent) NextFrame() {
	if len(ac.CurrentAnimation.Frames) == 0 {
		log.Println("No frame data for this animation")
		return
	}

	ac.index++
	ac.change = 0
	if ac.index >= len(ac.CurrentAnimation.Frames) {
		ac.index = 0

		if !ac.CurrentAnimation.Loop {
			ac.CurrentAnimation = nil
			return
		}
	}
}

// AnimationSystem tracks AnimationComponents, advancing their current animation.
type AnimationSystem struct {
	speed         float32
	previousSpeed float32

	entities map[uint64]animationEntity
}

type animationEntity struct {
	*AnimationComponent
	*common.RenderComponent
}

func (self *AnimationSystem) New(w *ecs.World) {
	self.speed = 1.0
	self.previousSpeed = 1.0
	log.Println("[AnimationSystem New]", self.speed)

	engo.Mailbox.Listen(messages.ControlMessageType, self.HandleControlMessage)
}

// Add starts tracking the given entity.
func (a *AnimationSystem) Add(basic *ecs.BasicEntity, anim *AnimationComponent, render *common.RenderComponent) {
	if a.entities == nil {
		a.entities = make(map[uint64]animationEntity)
	}
	a.entities[basic.ID()] = animationEntity{anim, render}
}

// Remove stops tracking the given entity.
func (a *AnimationSystem) Remove(basic ecs.BasicEntity) {
	if a.entities != nil {
		delete(a.entities, basic.ID())
	}
}

// Update advances the animations of all tracked entities.
func (a *AnimationSystem) Update(dt float32) {
	if a.speed == 0 {
		return
	}
	for _, e := range a.entities {
		if e.AnimationComponent.CurrentAnimation == nil {
			if e.AnimationComponent.def == nil {
				continue
			}
			e.AnimationComponent.SelectAnimationByAction(e.AnimationComponent.def)
		}

		e.AnimationComponent.change += dt
		if e.AnimationComponent.change*a.speed >= e.AnimationComponent.Rate {
			e.RenderComponent.Drawable = e.AnimationComponent.Cell()
			e.AnimationComponent.NextFrame()
		}
	}
}

func (self *AnimationSystem) HandleControlMessage(m engo.Message) {
	msg, ok := m.(messages.ControlMessage)
	if !ok {
		return
	}
	log.Printf("[AnimationSystem] %+v", m)
	switch msg.Action {
	case "TogglePause":
		if self.speed > 0 {
			self.previousSpeed = self.speed
			self.speed = 0
		} else {
			self.speed = self.previousSpeed
			self.previousSpeed = 0
		}
	}
}
