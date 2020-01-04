package main

type SlackBlockIface interface{}

// SlackBlocks is a collection of SlackBlock
type SlackBlocks []SlackBlock

// SlackMessage represents a slack UI message response
type SlackMessage struct {
	Blocks          SlackBlocks `json:"blocks"`
	ResponseType    string      `json:"response_type,omitempty"`
	ReplaceOriginal bool        `json:"replace_original"`
}

// NewSlackMessage creates a new SlackMessage
func NewSlackMessage() *SlackMessage {
	return &SlackMessage{
		Blocks:          make(SlackBlocks, 0),
		ReplaceOriginal: false,
	}
}

// AddBlock adds a new slack UI block to a message
func (m *SlackMessage) AddBlock(block *SlackBlock) {
	m.Blocks = append(m.Blocks, *block)
}

// AddDivider adds a slack UI divider into the message
func (m *SlackMessage) AddDivider() {
	m.Blocks = append(m.Blocks, SlackBlock{
		Type: "divider",
	})
}

// SlackBlock represents a Slack UI block message
type SlackBlock struct {
	Type         string           `json:"type"`
	ResponseType string           `json:"response_type,omitempty"`
	Accessory    *SlackAccessory  `json:"accessory,omitempty"`
	Text         *SlackText       `json:"text,omitempty"`
	Fields       []SlackText      `json:"fields,omitempty"`
	Elements     []SlackAccessory `json:"elements,omitempty"`
}

// SlackText is a slack UI text block
type SlackText struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji,omitempty"`
}

// SlackAccessory is a slack UI accessory block
type SlackAccessory struct {
	Type     string     `json:"type"`
	Text     *SlackText `json:"text,omitempty"`
	ImageURL string     `json:"image_url,omitempty"`
	AltText  string     `json:"alt_text,omitempty"`
	Value    string     `json:"value,omitempty"`
	ActionTS string     `json:"action_ts,omitempty"`
}

// SlackBlockAction is an action block sent in response to a slack interaction
type SlackBlockAction struct {
	Type        string            `json:"type"`
	Team        map[string]string `json:"team"`
	User        map[string]string `json:"user"`
	APIAppID    string            `json:"api_app_id"`
	Token       string            `json:"token"`
	TriggerID   string            `json:"trigger_id"`
	ResponseURL string            `json:"response_url"`
	Actions     []SlackAccessory  `json:"actions,omitempty"`
}
