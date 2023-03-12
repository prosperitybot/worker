package component

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type SettingsNotificationComponent struct {
	discord.Component
	db *sqlx.DB
}

func (s SettingsNotificationComponent) BaseComponent() discordgo.MessageComponent {
	return discordgo.SelectMenu{
		CustomID: "settings::notifications",
		MenuType: discordgo.StringSelectMenu,
		Options: []discordgo.SelectMenuOption{
			{
				Label:       "Reply to Message",
				Description: "Reply to the message that triggered the level up",
				Value:       "settings::notifications::reply",
				Emoji: discordgo.ComponentEmoji{
					Name: "💬",
				},
			},
			{
				Label:       "Specify Channel",
				Description: "Specify a channel to send the notifications to",
				Value:       "settings::notifications::channel",
				Emoji: discordgo.ComponentEmoji{
					Name: "📃",
				},
			},
			{
				Label:       "Direct Messages",
				Description: "Send the notifications to the user's DMs",
				Value:       "settings::notifications::dm",
				Emoji: discordgo.ComponentEmoji{
					Name: "🔏",
				},
			},
			{
				Label:       "Disable Notifications",
				Description: "Disable notifications for the server",
				Value:       "settings::notifications::disable",
				Emoji: discordgo.ComponentEmoji{
					Name: "🚫",
				},
			},
		},
	}
}

func (s SettingsNotificationComponent) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		notificationTypeValue = i.MessageComponentData().Values[0]
		notificationType      string
		responseMsg           string
	)

	switch notificationTypeValue {
	case "settings::notifications::reply":
		notificationType = "reply"
		responseMsg = "Successfully updated level up notifications to be sent via **replies**"
	case "settings::notifications::channel":
		notificationType = "NOT_UPDATED"
		responseMsg = "Please use the slash command `/settings notifications` and specify the channel"
	case "settings::notifications::dm":
		notificationType = "dm"
		responseMsg = "Successfully updated level up notifications to be sent via **Direct Messages**"
	case "settings::notifications::disable":
		notificationType = "disable"
		responseMsg = "Successfully updated level up notifications to be **disabled**"
	default:
		notificationType = "NOT_UPDATED"
		responseMsg = "Please use the slash command `/settings notifications` command and choose an option"
	}

	if notificationType != "NOT_UPDATED" {
		if _, err := s.db.Exec("UPDATE guilds SET notificationType = ?, notificationChannel = NULL WHERE id = ?", notificationType, i.GuildID); err != nil {
			logger.Error(c.Request().Context(), "failed to update guild notification type", zap.Error(err))
			utils.SendResponse(c, "Failed to update guild notification type", true, true)
		}
	}

	utils.SendResponse(c, responseMsg, true, false)
}

func NewSettingsNotificationComponent(db *sqlx.DB) SettingsNotificationComponent {
	return SettingsNotificationComponent{
		db: db,
	}
}
