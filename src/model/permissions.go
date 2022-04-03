package model

type Role struct {
	ID          Snowflake `json:"id"`
	Name        string    `json:"name"`
	Color       int       `json:"color"`
	Hoist       bool      `json:"hoist"`
	Position    uint8     `json:"position"`
	Permissions string    `json:"permissions"`
	Managed     bool      `json:"managed"`
	Mentionable bool      `json:"mentionable"`
	Tags        RoleTags  `json:"tags,omitempty"`
}

type RoleTags struct {
	BotID             Snowflake `json:"bot_id,omitempty"`
	IntegrationID     Snowflake `json:"integration_id,omitempty"`
	PremiumSubscriber string    `json:"premium_subscriber,omitempty"`
}
