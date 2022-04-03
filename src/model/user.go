package model

type UserFlags uint32

const (
	FlagNone                      UserFlags = 0
	FlagDiscordEmployee           UserFlags = 1 << 0
	FlagPartneredServerOwner      UserFlags = 1 << 1
	FlagHypeSquadEvents           UserFlags = 1 << 2
	FlagBugHunterLevel1           UserFlags = 1 << 3
	FlagHouseBravery              UserFlags = 1 << 6
	FlagHouseBrilliance           UserFlags = 1 << 7
	FlagHouseBalance              UserFlags = 1 << 8
	FlagEarlySupporter            UserFlags = 1 << 9
	FlagTeamUser                  UserFlags = 1 << 10
	FlagBugHunterLevel2           UserFlags = 1 << 14
	FlagVerifiedBot               UserFlags = 1 << 16
	FlagEarlyVerifiedBotDeveloper UserFlags = 1 << 17
	FlagDiscordCertifiedModerator UserFlags = 1 << 18
)

type PremiumType uint8

const (
	PremiumNone PremiumType = iota
	PremiumNitroClassic
	PremiumNitro
)

type User struct {
	ID            Snowflake   `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	Avatar        string      `json:"avatar"`
	Bot           bool        `json:"bot,omitempty"`
	System        bool        `json:"system,omitempty"`
	MFAEnabled    bool        `json:"mfa_enabled,omitempty"`
	Banner        string      `json:"banner,omitempty"`
	AccentColor   int         `json:"accent_color,omitempty"`
	Locale        string      `json:"locale,omitempty"`
	Verified      bool        `json:"verified,omitempty"`
	Email         string      `json:"email,omitempty"`
	Flags         UserFlags   `json:"flags,omitempty"`
	PremiumType   PremiumType `json:"premium_type,omitempty"`
	PublicFlags   UserFlags   `json:"public_flags,omitempty"`
}
