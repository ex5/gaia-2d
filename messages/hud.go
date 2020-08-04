package messages

import (
	"time"
)

const HUDTextUpdateMessageType string = "HUDTextUpdateMessage"

type HUDTextUpdateMessage struct {
	Name      string
	HideAfter time.Duration
	GetText   func() string
}

func (HUDTextUpdateMessage) Type() string {
	return HUDTextUpdateMessageType
}
