package messages

type ClientMessage struct {
	Type string `json:"type"`
	Data string `json:"data"`

	// Very sketchy way to handle dynamic data
	AdditionalData map[string]interface{} `json:"additionalData"`
}
