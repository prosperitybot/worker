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

type XpCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m XpCommand) Command() discordgo.ApplicationCommand {
	var (
		minLevel                 = float64(1)
		defaultPermissions int64 = 0
		dmAccess           bool  = false
	)
	return discordgo.ApplicationCommand{
		Name:                     "xp",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manages user xp",
		DefaultMemberPermissions: &defaultPermissions,
		DMPermission:             &dmAccess,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "give",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Gives xp to a user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "user",
						Type:        discordgo.ApplicationCommandOptionUser,
						Description: "The user to give xp to",
						Required:    true,
					},
					{
						Name:        "xp",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The amount of xp to give",
						Required:    true,
						MinValue:    &minLevel,
					},
				},
			},
			{
				Name:        "take",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Takes xp from a user",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "user",
						Type:        discordgo.ApplicationCommandOptionUser,
						Description: "The user to take xp from",
						Required:    true,
					},
					{
						Name:        "levels",
						Type:        discordgo.ApplicationCommandOptionInteger,
						Description: "The amount of xp to take",
						Required:    true,
						MinValue:    &minLevel,
					},
				},
			},
		},
	}
}

func (m XpCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "give":
		m.subcmd(c, i, subCommand, true)
	case "take":
		m.subcmd(c, i, subCommand, false)
	}
}

func (m XpCommand) subcmd(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption, shouldGive bool) {
	var (
		userId     = subCommand.Options[0].UserValue(nil).ID
		xp         = subCommand.Options[1].IntValue()
		guildUser  model.GuildUser
		levelFound = false
	)

	if err := m.db.Get(&guildUser, "SELECT * FROM guild_users WHERE guildId = ? AND userId = ?", i.GuildID, userId); err != nil {
		logger.Error(c.Request().Context(), "Error whilst getting user xp", zap.Error(err))
		utils.SendResponse(c, "Error getting user", true, true)
		return
	}

	prefix := "Given"

	if shouldGive {
		guildUser.Xp += xp
	} else {
		guildUser.Xp -= xp
		prefix = "Taken"
	}

	for !levelFound {
		if shouldGive {
			if guildUser.Xp >= utils.GetXPRequired(guildUser.Level) {
				guildUser.Level++
			} else {
				levelFound = true
			}
		} else {
			if guildUser.Xp < utils.GetXPRequired(guildUser.Level) {
				guildUser.Level--
			} else {
				levelFound = true
			}
		}
	}

	if _, err := m.db.NamedExecContext(c.Request().Context(), "UPDATE guild_users SET level = :level, xp = :xp WHERE guildId = :guildId AND userId = :userId", guildUser); err != nil {
		logger.Error(c.Request().Context(), "Error whilst updating user xp", zap.Error(err))
		utils.SendResponse(c, "Error updating user", true, true)
		return
	}
	utils.SendResponse(c, fmt.Sprintf("%s **%d** xp to <@%s>", prefix, xp, userId), false, false)
}

func NewXpCommand(db *sqlx.DB) XpCommand {
	return XpCommand{db: db}
}
