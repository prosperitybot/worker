package command

import (
	"database/sql"
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

type LevelCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m LevelCommand) Command() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "level",
		Type:        discordgo.ChatApplicationCommand,
		Description: "See your or someone elses current level",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "user",
				Type:        discordgo.ApplicationCommandOptionUser,
				Description: "The user you want to check the level of",
				Required:    false,
			},
		},
	}
}

func (m LevelCommand) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		userId  = i.Member.User.ID
		guildId = i.GuildID
	)

	if len(i.ApplicationCommandData().Options) > 0 {
		userId = i.ApplicationCommandData().Options[0].UserValue(nil).ID
	}

	var guildUser model.GuildUser

	if err := m.db.Get(&guildUser, "SELECT * FROM guild_users WHERE guildId = ? AND userId = ?", guildId, userId); err != nil {
		if err == sql.ErrNoRows {
			utils.SendResponse(c, fmt.Sprintf("<@%s> has never talked before", userId), true, true)
			return
		}
		logger.Error(c.Request().Context(), "Error whilst getting user level", zap.Error(err))
		utils.SendResponse(c, "Error getting guild user", true, true)
		return
	}

	var (
		xpNeeded    = utils.GetXPRequired(guildUser.Level+1) - guildUser.Xp
		responseMsg = fmt.Sprintf("Your current level is **%d**\nYou need **%d** xp to get to the next level", guildUser.Level, xpNeeded)
	)

	if userId != i.Member.User.ID {
		responseMsg = fmt.Sprintf("<@%s>'s current level is **%d**\nThey need **%d** xp to get to the next level", userId, guildUser.Level, xpNeeded)
	}

	utils.SendResponse(c, responseMsg, false, false)
}

func NewLevelCommand(db *sqlx.DB) LevelCommand {
	return LevelCommand{db: db}
}
