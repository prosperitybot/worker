package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/model"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type LevelsCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m LevelsCommand) Command() discordgo.ApplicationCommand {
	var (
		minLevel                 = float64(1)
		defaultPermissions int64 = 0
		dmAccess           bool  = false
	)
	return discordgo.ApplicationCommand{
		Name:                     "levels",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manages user levels",
		DefaultMemberPermissions: &defaultPermissions,
		DMPermission:             &dmAccess,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "give",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Gives levels to a user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "user",
						Type:        discordgo.ApplicationCommandOptionUser,
						Description: "The user to give levels to",
						Required:    true,
					},
					{
						Name:        "levels",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The amount of levels to give",
						Required:    true,
						MinValue:    &minLevel,
					},
				},
			},
			{
				Name:        "take",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Takes levels from a user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "user",
						Type:        discordgo.ApplicationCommandOptionUser,
						Description: "The user to take levels from",
						Required:    true,
					},
					{
						Name:        "levels",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The amount of levels to take",
						Required:    true,
						MinValue:    &minLevel,
					},
				},
			},
		},
	}
}

func (m LevelsCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "give":
		m.subcmd(c, i, subCommand, true)
	case "take":
		m.subcmd(c, i, subCommand, false)
	}
}

func (m LevelsCommand) subcmd(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption, shouldGive bool) {
	var (
		userId    = subCommand.Options[0].UserValue(nil).ID
		levels    = subCommand.Options[1].IntValue()
		guildUser model.GuildUser
	)

	if err := m.db.Get(&guildUser, "SELECT * FROM guild_users WHERE guildId = ? AND userId = ?", i.GuildID, userId); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting user level", zap.Error(err))
		utils.SendResponse(c, "Error getting user", true, true)
		return
	}

	if shouldGive {
		guildUser.Level += int(levels)
	} else {
		guildUser.Level -= int(levels)
	}
	guildUser.Xp = utils.GetXPRequired(guildUser.Level-1) + 1

	if _, err := m.db.NamedExecContext(c.Request().Context(), "UPDATE guild_users SET level = :level, xp = :xp WHERE guildId = :guildId AND userId = :userId", guildUser); err != nil {
		logger.Error(c.Request().Context(), "Error whilst updating user level", zap.Error(err))
		utils.SendResponse(c, "Error updating user", true, true)
		return
	}

	utils.SendResponse(c, fmt.Sprintf("Given **%d** level(s) to <@%s>", levels, userId), false, false)
}

func NewLevelsCommand(db *sqlx.DB) LevelsCommand {
	return LevelsCommand{db: db}
}
