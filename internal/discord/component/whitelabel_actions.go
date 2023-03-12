package component

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

type WhitelabelActionsComponent struct {
	discord.Component
	db *sqlx.DB
}

func (s WhitelabelActionsComponent) BaseComponent() discordgo.MessageComponent {
	return discordgo.SelectMenu{
		CustomID: "whitelabel::actions_FAKE_BOT_ID",
		MenuType: discordgo.StringSelectMenu,
		Options:  []discordgo.SelectMenuOption{},
	}
}

func (s WhitelabelActionsComponent) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		botId  = strings.Split(i.MessageComponentData().CustomID, "_")[1]
		action = i.MessageComponentData().Values[0]
	)

	if err := s.db.Get(&botId, "UPDATE whitelabel_bots SET action = ? WHERE botId = ?", action, botId); err != nil {
		logger.Error(c.Request().Context(), "failed to update whitelabel bot action", zap.Error(err))
	}

	utils.SendResponse(c, fmt.Sprintf("Whitelabel bot has been set to `%s`", action), true, false)
}

func NewWhitelabelActionsComponent(db *sqlx.DB) WhitelabelActionsComponent {
	return WhitelabelActionsComponent{
		db: db,
	}
}
