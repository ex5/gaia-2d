package messages

import (
	"time"
)

const HUDTextUpdateMessageType string = "HUDTextUpdateMessage"

type HUDTextUpdateMessage struct {
	Name      string
	HideAfter time.Duration
	Lines     []string
}

func (HUDTextUpdateMessage) Type() string {
	return HUDTextUpdateMessageType
}
