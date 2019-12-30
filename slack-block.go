package main

// SlackBlocks is a collection of SlackBlock
type SlackBlocks []SlackBlock

// SlackMessage represents a slack UI message response
type SlackMessage struct {
	Blocks SlackBlocks `json:"blocks"`
}

// NewSlackMessage creates a new SlackMessage
func NewSlackMessage() *SlackMessage {
	return &SlackMessage{
		Blocks: make(SlackBlocks, 0),
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
}
