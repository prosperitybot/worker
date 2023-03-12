package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"github.com/prosperitybot/worker/internal/discord/component"
	"go.uber.org/zap"
)

type SettingsCommand struct {
	discord.SlashCommand
	settingNotificationComponent component.SettingsNotificationComponent
	db                           *sqlx.DB
}

func (m SettingsCommand) Command() discordgo.ApplicationCommand {
	var (
		defaultPermissions int64   = 0
		dmAccess           bool    = false
		minMultiplierValue float64 = 0.0
		minDelay           float64 = 1
	)
	return discordgo.ApplicationCommand{
		Name:                     "settings",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manages settings for the server",
		DefaultMemberPermissions: &defaultPermissions,
		DMPermission:             &dmAccess,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "notifications",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Choose where the notifications are displayed for the server",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "channel",
						Type:        discordgo.ApplicationCommandOptionChannel,
						Description: "The channel to send the notifications to",
						Required:    false,
					},
				},
			},
			{
				Name:        "roles",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Choose the type of logic to apply on the role (stack/single)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "type",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "The type of logic to apply to the role",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Single (Only apply one at a time and remove the previous role)",
								Value: "single",
							},
							{
								Name:  "Stack (Stack all previous roles and do not remove old ones)",
								Value: "stack",
							},
						},
					},
				},
			},
			{
				Name:        "multiplier",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Choose the XP multiplier for the server (1x by default)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "multiplier",
						Type:        discordgo.ApplicationCommandOptionNumber,
						Description: "The multiplier to apply to the XP",
						Required:    true,
						MinValue:    &minMultiplierValue,
						MaxValue:    10.0,
					},
				},
			},
			{
				Name:        "delay",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Choose the delay between each message (in seconds)",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "delay",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The delay between each message",
						Required:    true,
						MinValue:    &minDelay,
						MaxValue:    60 * 60 * 24,
					},
				},
			},
		},
	}
}

func (m SettingsCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "notifications":
		m.subcmd_notifications(c, i, subCommand)
	case "roles":
		m.subcmd_roles(c, i, subCommand)
	case "multiplier":
		m.subcmd_multiplier(c, i, subCommand)
	case "delay":
		m.subcmd_delay(c, i, subCommand)
	}
}

func (m SettingsCommand) subcmd_notifications(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	if len(subCommand.Options) > 0 {
		// Has supplied a channel
		channelId := subCommand.Options[0].ChannelValue(nil).ID

		if _, err := m.db.Exec("UPDATE guilds SET notificationType = ?, notificationChannel = ? WHERE id = ?", "channel", channelId, i.GuildID); err != nil {
			logger.Error(c.Request().Context(), "failed to update guild settings", zap.Error(err))
			utils.SendResponse(c, "Failed to update guild settings", true, true)
			return
		}

		utils.SendResponse(c, fmt.Sprintf("Set the notifications channel to <#%s>", channelId), true, false)
	} else {
		// Has not supplied a channel, go with other
		var (
			embed = utils.CreateEmbed(&discordgo.MessageEmbed{
				Description: "Please select the type of notifications you below",
			}, false)
			components = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						m.settingNotificationComponent.BaseComponent(),
					},
				},
			}
		)

		utils.SendComplexResponse(c, discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsEphemeral,
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		})
	}
}

func (m SettingsCommand) subcmd_roles(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		roleAssignmentType = subCommand.Options[0].StringValue()
	)

	if _, err := m.db.Exec("UPDATE guilds SET roleAssignType = ? WHERE id = ?", roleAssignmentType, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "failed to update guild settings", zap.Error(err))
		utils.SendResponse(c, "Failed to update guild settings", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("Set the role assignment type to `%s`", roleAssignmentType), true, false)
}

func (m SettingsCommand) subcmd_multiplier(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		multiplier = subCommand.Options[0].FloatValue()
	)

	if _, err := m.db.Exec("UPDATE guilds SET xpRate = ? WHERE id = ?", multiplier, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "failed to update guild settings", zap.Error(err))
		utils.SendResponse(c, "Failed to update guild settings", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("Set the XP multiplier to `%f`", multiplier), true, false)
}

func (m SettingsCommand) subcmd_delay(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		delay = subCommand.Options[0].IntValue()
	)

	if _, err := m.db.Exec("UPDATE guilds SET xpDelay = ? WHERE id = ?", delay, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "failed to update guild settings", zap.Error(err))
		utils.SendResponse(c, "Failed to update guild settings", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("Set the XP delay to `%d`", delay), true, false)
}

func NewSettingsCommand(db *sqlx.DB, settingsComponent component.SettingsNotificationComponent) SettingsCommand {
	return SettingsCommand{db: db}
}
