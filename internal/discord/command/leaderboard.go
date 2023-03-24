package command

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/model"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type LeaderboardCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m LeaderboardCommand) Command() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        "leaderboard",
		Type:        discordgo.ChatApplicationCommand,
		Description: "Displays the top users and their levels",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "page",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Description: "The page you want to display",
				Required:    false,
			},
		},
	}
}

func (m LeaderboardCommand) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		page             = 1
		pageSize         = 10
		guildId          = i.GuildID
		guildUsers       []model.GuildUser
		userCount        int
		leaderboardLines []string
		offset           = pageSize * (page - 1)
	)

	if len(i.ApplicationCommandData().Options) > 0 {
		page = int(i.ApplicationCommandData().Options[0].IntValue())
	}

	if err := m.db.GetContext(c.Request().Context(), &userCount, "SELECT COUNT(*) FROM guild_users WHERE guildId = ?", guildId); err != nil {
		logger.Error(c.Request().Context(), "Error getting amount of users in guild for leaderboard", zap.Error(err))
		utils.SendResponse(c, "Could not fetch leaderboard", true, true)
		return
	}

	leaderboardQuery := `SELECT gu.*, CONCAT(u.username, "#", u.discriminator) AS username FROM guild_users gu INNER JOIN users u ON gu.userId = u.id WHERE guildId = ? ORDER BY xp DESC LIMIT %d OFFSET %d`
	leaderboardQuery = fmt.Sprintf(leaderboardQuery, pageSize, offset)

	if err := m.db.SelectContext(c.Request().Context(), &guildUsers, leaderboardQuery, guildId); err != nil {
		logger.Error(c.Request().Context(), "Error getting list of users for the leaderboard", zap.Error(err))
		utils.SendResponse(c, "Error getting leaderboard", true, true)
		return
	}

	for i := range guildUsers {
		userIndex := i + offset
		leaderboardLines = append(leaderboardLines, fmt.Sprintf("%d. %s - Level %d", userIndex, guildUsers[i].Username, guildUsers[i].Level))
	}
	responseMsg := fmt.Sprintf("Top 10 Members (Page %d of %d)\n\n %s", page, userCount/pageSize, strings.Join(leaderboardLines, "\n"))

	utils.SendResponse(c, responseMsg, false, false)
}

func NewLeaderboardCommand(db *sqlx.DB) LeaderboardCommand {
	return LeaderboardCommand{db: db}
}
