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

type AboutCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m AboutCommand) Command() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "about",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Information about the bot",
	}
}

func (m AboutCommand) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		aboutStats model.AboutStats
	)

	if err := m.db.Get(&aboutStats, "SELECT (SELECT COUNT(id) FROM guilds WHERE active = true) AS servers, COUNT(DISTINCT guildId, userId) AS users FROM guild_users"); err != nil {
		logger.Error(c.Request().Context(), "failed to get about stats", zap.Error(err))
		utils.SendResponse(c, "Failed to load the about command", true, true)
		return
	}

	responseMsg := "Prosperity is a levelling bot ready to skill up and boost up your Discord server. We pride ourselves on openness, transparency and collaboration	"
	embed := utils.CreateEmbed(&discordgo.MessageEmbed{
		Description: responseMsg,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Bot Statistics",
				Value:  fmt.Sprintf("Servers: %d\nMembers: %d", aboutStats.Servers, aboutStats.Users),
				Inline: true,
			},
		},
	}, false)

	utils.SendComplexResponse(c, discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}})
}

func NewAboutCommand(db *sqlx.DB) AboutCommand {
	return AboutCommand{db: db}
}
