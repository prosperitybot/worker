package model

import (
	"time"
)

type ChannelType uint8
type EmbedType string

const (
	GuildTextChannel ChannelType = iota
	DMChannel
	GuildVoiceChannel
	GroupDMChannel
	GuildCategoryChannel
	GuildNewsChannel
	GuildStoreChannel
	GuildNewsThreadChannel = iota + 3
	GuildPublicThreadChannel
	GuildPrivateThreadChannel
	GuildStageVoiceChannel
)

const (
	EmbedRichType    EmbedType = "rich"
	EmbedImageType   EmbedType = "image"
	EmbedVideoType   EmbedType = "video"
	EmbedGifvType    EmbedType = "gifv"
	EmbedArticleType EmbedType = "article"
	EmbedLinkType    EmbedType = "link"
)

type PartialChannel struct {
	ID             Snowflake      `json:"id"`
	Name           string         `json:"name,omitempty"`
	Type           ChannelType    `json:"type"`
	Permissions    string         `json:"permissions"`
	ThreadMetadata ThreadMetaData `json:"thread_metadata,omitempty"`
	ParentID       Snowflake      `json:"parent_id,omitempty"`
}

type ThreadMetaData struct {
	Archived            bool   `json:"archived"`
	AutoArchiveDuration uint8  `json:"auto_archive_duration"`
	ArchiveTimestamp    string `json:"archive_timestamp"`
	Locked              bool   `json:"locked"`
	Invitable           bool   `json:"invitable,omitempty"`
}

type AllowedMentions struct {
	Parse       []string    `json:"parse,omitempty"`
	Roles       []Snowflake `json:"roles,omitempty"`
	Users       []Snowflake `json:"users,omitempty"`
	RepliedUser bool        `json:"replied_user,omitempty"`
}

type MessageEmbed struct {
	Title       string         `json:"title,omitempty"`
	Type        EmbedType      `json:"type,omitempty"`
	Description string         `json:"description,omitempty"`
	URL         string         `json:"url,omitempty"`
	Timestamp   *time.Time     `json:"timestamp,omitempty"`
	Color       uint32         `json:"color,omitempty"`
	Footer      EmbedFooter    `json:"footer,omitempty"`
	Image       EmbedImage     `json:"image,omitempty"`
	Thumbnail   EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       EmbedVideo     `json:"video,omitempty"`
	Provider    EmbedProvider  `json:"provider,omitempty"`
	Author      EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int8   `json:"height,omitempty"`
	Width    int8   `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int8   `json:"height,omitempty"`
	Width    int8   `json:"width,omitempty"`
}

type EmbedVideo struct {
	URL      string `json:"url,omitempty"`
	ProxyURL string `json:"proxy_url,omitempty"`
	Height   int8   `json:"height,omitempty"`
	Width    int8   `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name,omitempty"`
	URL          string `json:"url,omitempty"`
	IconURL      string `json:"icon_url,omitempty"`
	ProxyIconURL string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}
