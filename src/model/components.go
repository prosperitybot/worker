package model

type ComponentType uint8
type ButtonStyleType uint8

const (
	_ ComponentType = iota
	ActionRowComponent
	ButtonComponent
	SelectMenuComponent
)

const (
	_ ButtonStyleType = iota
	PrimaryButtonStyle
	SecondaryButtonStyle
	SuccessButtonStyle
	DangerButtonStyle
	LinkButtonStyle
)

type DiscordButton interface {
	GetButtonType() ButtonStyleType
}

type Component interface {
	GetComponentType() ComponentType
}
type ActionRow struct {
	Type       ComponentType `json:"type"`
	Components []Component   `json:"components"`
}

func (ar ActionRow) GetComponentType() ComponentType {
	return ar.Type
}

type LinkButton struct {
	Type     ComponentType   `json:"type"`
	Style    ButtonStyleType `json:"style,omitempty"`
	Label    string          `json:"label,omitempty"`
	Emoji    PartialEmoji    `json:"emoji,omitempty"`
	URL      string          `json:"url,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (lb LinkButton) GetComponentType() ComponentType {
	return lb.Type
}

type ActionButton struct {
	Type     ComponentType   `json:"type"`
	Style    ButtonStyleType `json:"style,omitempty"`
	Label    string          `json:"label,omitempty"`
	Emoji    PartialEmoji    `json:"emoji,omitempty"`
	CustomID string          `json:"custom_id,omitempty"`
	Disabled bool            `json:"disabled,omitempty"`
}

func (ab ActionButton) GetComponentType() ComponentType {
	return ab.Type
}
