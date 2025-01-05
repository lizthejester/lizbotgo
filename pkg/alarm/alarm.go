package alarm

import "time"

type Alarm struct {
	ChannelID string
	Deadline  *time.Time
	Content   string
}
