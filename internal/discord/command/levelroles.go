package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/model"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type LevelRolesCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m LevelRolesCommand) Command() discordgo.ApplicationCommand {
	minLevel := float64(1)
	var (
		defaultPermissions int64 = 0
		dmAccess           bool  = false
	)
	return discordgo.ApplicationCommand{
		Name:                     "levelroles",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manages level roles",
		DefaultMemberPermissions: &defaultPermissions,
		DMPermission:             &dmAccess,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "add",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Adds a level role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Type:        discordgo.ApplicationCommandOptionRole,
						Description: "The role to give",
						Required:    true,
					},
					{
						Name:        "level",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The level to give the role at",
						Required:    true,
						MinValue:    &minLevel,
					},
				},
			},
			{
				Name:        "remove",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Removes a level role",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "role",
						Type:        discordgo.ApplicationCommandOptionRole,
						Description: "The role to remove",
						Required:    true,
					},
				},
			},
			{
				Name:        "list",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Lists all level roles",
			},
		},
	}
}

func (m LevelRolesCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "add":
		m.subcmd_add(c, i, subCommand)
	case "remove":
		m.subcmd_remove(c, i, subCommand)
	case "list":
		m.subcmd_list(c, i, subCommand)
	}
}

func (m LevelRolesCommand) subcmd_add(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		role          = subCommand.Options[0].RoleValue(nil, i.GuildID).ID
		level         = int(subCommand.Options[1].IntValue())
		alreadyExists = false
	)

	if err := m.db.GetContext(c.Request().Context(), &alreadyExists, "SELECT exists(SELECT 1 FROM level_roles WHERE guildId = ? AND (level = ? OR id = ?))", i.GuildID, level, role); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether levelrole exists", zap.Error(err))
		utils.SendResponse(c, "Error getting level roles", true, true)
		return
	}

	if alreadyExists {
		utils.SendResponse(c, "Level role already exists", true, true)
		return
	}

	levelRole := &model.LevelRole{
		GuildId:   i.GuildID,
		Level:     level,
		Id:        role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if _, err := m.db.NamedExecContext(c.Request().Context(), "INSERT INTO level_roles (guildId, level, id, createdAt, updatedAt) VALUES (:guildId, :level, :id, :createdAt, :updatedAt)", levelRole); err != nil {
		logger.Error(c.Request().Context(), "Error whilst creating the new levelrole", zap.Error(err))
		utils.SendResponse(c, "Error adding level role", true, true)
		return
	}

	var usersNeedingRole []string
	usersNeedingRoleQuery := "SELECT userId FROM guild_users WHERE guildId = ? AND level >= ? AND level < COALESCE((SELECT level FROM level_roles WHERE guildId = ? AND level > ? ORDER BY level ASC LIMIT 1), 9999)"

	if err := m.db.SelectContext(c.Request().Context(), &usersNeedingRole, usersNeedingRoleQuery, i.GuildID, level, i.GuildID, level); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting a list of users to assign the level role to", zap.Error(err))
		utils.SendResponse(c, "Error getting list of users to apply the role to", true, true)
		return
	}

	for _, userId := range usersNeedingRole {
		if err := levelRole.AddToMember(userId, fmt.Sprintf("New level role added (Level %d)", level)); err != nil {
			logger.Error(c.Request().Context(), "Error whilst adding the role", zap.String("roleId", levelRole.Id), zap.String("userToAdd", userId), zap.Error(err))
			utils.SendResponse(c, "Error adding role to users", true, true)
			return
		}
	}

	responseMsg := fmt.Sprintf("<@&%s> will be granted at level **%d**\n\nAssigning role to **%d** users", role, level, len(usersNeedingRole))

	utils.SendResponse(c, responseMsg, false, false)
}

func (m LevelRolesCommand) subcmd_remove(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		role          = subCommand.Options[0].RoleValue(nil, i.GuildID).ID
		alreadyExists = false
	)

	if err := m.db.GetContext(c.Request().Context(), &alreadyExists, "SELECT exists(SELECT 1 FROM level_roles WHERE guildId = ? AND id = ?)", i.GuildID, role); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether levelrole exists", zap.Error(err))
		utils.SendResponse(c, "Error getting level roles", true, true)
		return
	}

	if !alreadyExists {
		utils.SendResponse(c, "Level role does not exist", true, true)
		return
	}

	if _, err := m.db.NamedExecContext(c.Request().Context(), "DELETE FROM level_roles WHERE id = ?", role); err != nil {
		logger.Error(c.Request().Context(), "Error whilst deleting the levelrole", zap.Error(err))
		utils.SendResponse(c, "Error removing level role", true, true)
		return
	}

	responseMsg := fmt.Sprintf("<@&%s> has been removed as a level role", role)

	utils.SendResponse(c, responseMsg, false, false)
}

func (m LevelRolesCommand) subcmd_list(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		levelRoles []model.LevelRole
	)

	if err := m.db.SelectContext(c.Request().Context(), &levelRoles, "SELECT * FROM level_roles WHERE guildId = ?", i.GuildID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting list of level roles", zap.Error(err))
		utils.SendResponse(c, "Error getting level roles", true, true)
		return
	}

	levelRolesStrings := make([]string, len(levelRoles))

	for i := range levelRoles {
		levelRolesStrings[i] = fmt.Sprintf("- <@&%s> at level **%d**", levelRoles[i].Id, levelRoles[i].Level)
	}

	utils.SendResponse(c, fmt.Sprintf("**Level Roles**\n\n%s", strings.Join(levelRolesStrings, "\n")), false, false)
}

func NewLevelRolesCommand(db *sqlx.DB) LevelRolesCommand {
	return LevelRolesCommand{db: db}
}
