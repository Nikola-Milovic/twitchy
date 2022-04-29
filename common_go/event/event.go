package event

type BaseEvent struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}
