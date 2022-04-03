package model

type InteractionType uint8
type ApplicationCommandType uint8
type ApplicationCommandOptionType uint8
type InteractionCallbackType uint8
type InteractionCallbackDataFlags uint8

const (
	_ InteractionType = iota
	InteractionPing
	InteractionApplicationCommand
	InteractionMessageComponent
)

const (
	_ ApplicationCommandType = iota
	ChatInputCommand
	UserCommand
	MessageCommand
)

const (
	_ ApplicationCommandOptionType = iota
	SubCommandOption
	SubCommandGroup
	StringOption
	IntegerOption
	BooleanOption
	UserOption
	ChannelOption
	RoleOption
	MentionableOption
	NumberOption
)

const (
	PongCallback                             InteractionCallbackType = 1
	ChannelMessageWithSourceCallback         InteractionCallbackType = 4
	DeferredChannelMessageWithSourceCallback InteractionCallbackType = 5
	DeferredUpdateMessageCallback            InteractionCallbackType = 6
	UpdateMessageCallback                    InteractionCallbackType = 7
)

const Ephemeral InteractionCallbackDataFlags = 1 << 6

type Interaction struct {
	Id            Snowflake       `json:"id"`
	ApplicationID Snowflake       `json:"application_id"`
	Type          InteractionType `json:"type"`
	Data          InteractionData `json:"data,omitempty"`
	GuildID       Snowflake       `json:"guild_id,omitempty"`
	ChannelID     Snowflake       `json:"channel_id,omitempty"`
	Member        GuildMember     `json:"member,omitempty"`
	User          User            `json:"user,omitempty"`
	Token         string          `json:"token"`
	Version       uint8           `json:"version"`
	Message       interface{}     `json:"message,omitempty"`
}

type InteractionData struct {
	ID       Snowflake                                 `json:"id"`
	Name     string                                    `json:"name"`
	Type     ApplicationCommandType                    `json:"type"`
	Resolved ResolvedData                              `json:"resolved,omitempty"`
	Options  []ApplicationCommandInteractionDataOption `json:"options,omitempty"`
	CustomID string                                    `json:"custom_id"`
	Values   []SelectOptionValues                      `json:"values,omitempty"`
	TargetID Snowflake                                 `json:"target_id,omitempty"`
}

type ResolvedData struct {
	Users    map[Snowflake]User               `json:"users,omitempty"`
	Members  map[Snowflake]PartialGuildMember `json:"members,omitempty"`
	Roles    map[Snowflake]Role               `json:"roles,omitempty"`
	Channels map[Snowflake]PartialChannel     `json:"channels,omitempty"`
	Messages map[Snowflake]interface{}        `json:"messages,omitempty"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                   `json:"name"`
	Type    ApplicationCommandOptionType             `json:"type"`
	Value   interface{}                              `json:"value,omitempty"`
	Options *ApplicationCommandInteractionDataOption `json:"options,omitempty"`
}

type SelectOptionValues struct {
	Label       string        `json:"label"`
	Value       string        `json:"value"`
	Description string        `json:"description,omitempty"`
	Emoji       *PartialEmoji `json:"emoji,omitempty"`
	Default     bool          `json:"default,omitempty"`
}

type InteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
	Data InteractionCallbackData `json:"data,omitempty"`
}

type InteractionCallbackData struct {
	TTS             bool                         `json:"tts,omitempty"`
	Content         string                       `json:"content,omitempty"`
	Embeds          []MessageEmbed               `json:"embeds,omitempty"`
	AllowedMentions AllowedMentions              `json:"allowed_mentions,omitempty"`
	Flags           InteractionCallbackDataFlags `json:"flags,omitempty"`
	Components      []ActionRow                  `json:"components,omitempty"`
}
