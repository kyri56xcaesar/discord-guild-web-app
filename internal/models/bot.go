package models

type Bot struct {
	Name           string    `json:"username"`
	CreatedAt      string    `json:"joinedat"`
	Author         string    `json:"author"`
	Status         string    `json:"status"`
	RoomConstraits []string  `json:"rooms"`
	IsSinger       bool      `json:"issinger"`
	Triggers       []Trigger `json:"triggerwords"`
	Lines          []Line    `json:"linewords"`
}

type Trigger struct {
}

type Line struct {
}
