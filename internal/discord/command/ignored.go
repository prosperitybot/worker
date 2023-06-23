package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type IgnoredCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m IgnoredCommand) Command() discordgo.ApplicationCommand {
	var (
		defaultPermissions int64 = 0
		dmAccess           bool  = false
	)
	return discordgo.ApplicationCommand{
		Name:                     "ignored",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manage ignored channels and roles",
		DefaultMemberPermissions: &defaultPermissions,
		DMPermission:             &dmAccess,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "channels",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Description: "Manage ignored channels",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "add",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "Add a channel to the ignored list",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "channel",
								Type:        discordgo.ApplicationCommandOptionChannel,
								Description: "The channel to add to the ignored list",
								Required:    true,
							},
						},
					},
					{
						Name:        "remove",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "Remove a channel from the ignored list",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "channel",
								Type:        discordgo.ApplicationCommandOptionChannel,
								Description: "The channel to remove from the ignored list",
								Required:    true,
							},
						},
					},
					{
						Name:        "list",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "List all ignored channels",
					},
				},
			},
			{
				Name:        "roles",
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Description: "Manage ignored roles",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "add",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "Add a role to the ignored list",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "role",
								Type:        discordgo.ApplicationCommandOptionRole,
								Description: "The role to add to the ignored list",
								Required:    true,
							},
						},
					},
					{
						Name:        "remove",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "Remove a role from the ignored list",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "role",
								Type:        discordgo.ApplicationCommandOptionRole,
								Description: "The role to remove from the ignored list",
								Required:    true,
							},
						},
					},
					{
						Name:        "list",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Description: "List all ignored roles",
					},
				},
			},
		},
	}
}

func (m IgnoredCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "channels":
		m.subcmd_channels(c, i, subCommand)
	case "roles":
		m.subcmd_roles(c, i, subCommand)
	}
}

func (m IgnoredCommand) subcmd_channels(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	subSubCommand := subCommand.Options[0]

	switch subSubCommand.Name {
	case "add":
		m.subcmd_channels_add(c, i, subSubCommand)
	case "remove":
		m.subcmd_channels_remove(c, i, subSubCommand)
	case "list":
		m.subcmd_channels_list(c, i, subSubCommand)
	}
}

func (m IgnoredCommand) subcmd_channels_add(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		channelId     = subCommand.Options[0].ChannelValue(nil).ID
		alreadyExists = false
	)

	if err := m.db.Get(&alreadyExists, "SELECT EXISTS(SELECT 1 FROM ignored_channels WHERE id = ? AND guildId = ?)", channelId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether channel is already ignored", zap.Error(err))
		utils.SendResponse(c, "Could not check whether channel is already ignored", true, true)
		return
	}

	if alreadyExists {
		utils.SendResponse(c, fmt.Sprintf("<#%s> is already being ignored", channelId), true, true)
		return
	}

	if _, err := m.db.Exec("INSERT INTO ignored_channels (id, guildId) VALUES (?, ?)", channelId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst adding channel to ignored list", zap.Error(err))
		utils.SendResponse(c, "Could not add channel to ignored list", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("<#%s> will be ignored from gaining xp", channelId), false, false)
}

func (m IgnoredCommand) subcmd_channels_remove(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		channelId     = subCommand.Options[0].ChannelValue(nil).ID
		alreadyExists = false
	)

	if err := m.db.Get(&alreadyExists, "SELECT EXISTS(SELECT 1 FROM ignored_channels WHERE id = ? AND guildId = ?)", channelId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether channel is already ignored", zap.Error(err))
		utils.SendResponse(c, "Could not check whether channel is already ignored", true, true)
		return
	}

	if !alreadyExists {
		utils.SendResponse(c, fmt.Sprintf("<#%s> is not being ignored", channelId), true, true)
		return
	}

	if _, err := m.db.Exec("DELETE FROM ignored_channels WHERE id = ? AND guildId = ?", channelId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst removing channel from ignored list", zap.Error(err))
		utils.SendResponse(c, "Could not remove channel from ignored list", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("<#%s> will no longer be ignored from gaining xp", channelId), false, false)
}

func (m IgnoredCommand) subcmd_channels_list(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		channelIds []string
	)

	if err := m.db.SelectContext(c.Request().Context(), &channelIds, "SELECT id FROM ignored_channels WHERE guildId = ?", i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting list of ignored channels", zap.Error(err))
		utils.SendResponse(c, "Error getting ignored channels", true, true)
		return
	}

	ignoredChannelStrings := make([]string, len(channelIds))

	for i := range channelIds {
		ignoredChannelStrings[i] = fmt.Sprintf("- <#%s>", channelIds[i])
	}

	utils.SendResponse(c, fmt.Sprintf("**Ignored Channels**\n\n%s", strings.Join(ignoredChannelStrings, "\n")), false, false)
}

func (m IgnoredCommand) subcmd_roles(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	subSubCommand := subCommand.Options[0]

	switch subSubCommand.Name {
	case "add":
		m.subcmd_roles_add(c, i, subSubCommand)
	case "remove":
		m.subcmd_roles_remove(c, i, subSubCommand)
	case "list":
		m.subcmd_roles_list(c, i, subSubCommand)
	}
}

func (m IgnoredCommand) subcmd_roles_add(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		roleId        = subCommand.Options[0].RoleValue(nil, "").ID
		alreadyExists = false
	)

	if err := m.db.Get(&alreadyExists, "SELECT EXISTS(SELECT 1 FROM ignored_roles WHERE id = ? AND guildId = ?)", roleId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether role is already ignored", zap.Error(err))
		utils.SendResponse(c, "Could not check whether role is already ignored", true, true)
		return
	}

	if alreadyExists {
		utils.SendResponse(c, fmt.Sprintf("<@&%s> is already being ignored", roleId), true, true)
		return
	}

	if _, err := m.db.Exec("INSERT INTO ignored_roles (id, guildId) VALUES (?, ?)", roleId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst adding role to ignored list", zap.Error(err))
		utils.SendResponse(c, "Could not add role to ignored list", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("<@&%s> will be ignored from gaining xp", roleId), false, false)
}

func (m IgnoredCommand) subcmd_roles_remove(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		roleId        = subCommand.Options[0].RoleValue(nil, "").ID
		alreadyExists = false
	)

	if err := m.db.Get(&alreadyExists, "SELECT EXISTS(SELECT 1 FROM ignored_roles WHERE id = ? AND guildId = ?)", roleId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether role is already ignored", zap.Error(err))
		utils.SendResponse(c, "Could not check whether role is already ignored", true, true)
		return
	}

	if !alreadyExists {
		utils.SendResponse(c, fmt.Sprintf("<@&%s> is not being ignored", roleId), true, true)
		return
	}

	if _, err := m.db.Exec("DELETE FROM ignored_roles WHERE id = ? AND guildId = ?", roleId, i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst removing role from ignored list", zap.Error(err))
		utils.SendResponse(c, "Could not remove role from ignored list", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("<@&%s> will no longer be ignored from gaining xp", roleId), false, false)
}

func (m IgnoredCommand) subcmd_roles_list(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		roleIds []string
	)

	if err := m.db.SelectContext(c.Request().Context(), &roleIds, "SELECT id FROM ignored_roles WHERE guildId = ?", i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting list of ignored roles", zap.Error(err))
		utils.SendResponse(c, "Error getting ignored roles", true, true)
		return
	}

	ignoredRoleStrings := make([]string, len(roleIds))

	for i := range roleIds {
		ignoredRoleStrings[i] = fmt.Sprintf("- <@&%s>", roleIds[i])
	}

	utils.SendResponse(c, fmt.Sprintf("**Ignored Roles**\n\n%s", strings.Join(ignoredRoleStrings, "\n")), false, false)
}

func NewIgnoredCommand(db *sqlx.DB) IgnoredCommand {
	return IgnoredCommand{db: db}
}
