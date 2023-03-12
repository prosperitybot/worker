package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type InteractionHandler struct {
	Commands   map[string]discord.SlashCommand
	Components map[string]discord.Component
}

func (h InteractionHandler) POSTInteractions(c echo.Context) error {
	var body discordgo.Interaction
	if err := c.Bind(&body); err != nil {
		logger.Error(c.Request().Context(), "Error binding interaction", zap.Error(err))
		return err
	}

	var (
		botId = c.Param("bot_id")
	)

	switch body.Type {
	case discordgo.InteractionPing:
		return c.JSON(http.StatusOK, discordgo.InteractionResponse{
			Type: discordgo.InteractionResponsePong,
		})
	case discordgo.InteractionApplicationCommand:
		c = addContextInfo(c, body, botId)

		if cmd, v := h.Commands[body.ApplicationCommandData().Name]; !v {
			return c.NoContent(404)
		} else {
			logger.Info(c.Request().Context(), fmt.Sprintf("Executing command /%s", cmd.Command().Name), zap.String("command", body.ApplicationCommandData().Name))
			cmd.Execute(c, body)
		}
		break
	case discordgo.InteractionMessageComponent:
		c = addContextInfo(c, body, botId)

		if component, v := h.Components[body.MessageComponentData().CustomID]; !v {
			splitString := strings.Split(body.MessageComponentData().CustomID, "_")
			if comp, w := h.Components[splitString[0]]; !w {
				return c.NoContent(404)
			} else {
				logger.Info(c.Request().Context(), fmt.Sprintf("Handling component %s", splitString[0]), zap.String("component", body.MessageComponentData().CustomID))
				comp.Execute(c, body)
			}
		} else {
			logger.Info(c.Request().Context(), fmt.Sprintf("Handling component %s", body.MessageComponentData().CustomID), zap.String("component", body.MessageComponentData().CustomID))
			component.Execute(c, body)
		}
		break
	}

	return c.NoContent(404)
}

func addContextInfo(c echo.Context, body discordgo.Interaction, botId string) echo.Context {
	ctx := c.Request().Context()
	ctx = context.WithValue(ctx, utils.UserIdContextKey, body.Member.User.ID)
	ctx = context.WithValue(ctx, utils.GuildIdContextKey, body.GuildID)
	ctx = context.WithValue(ctx, utils.ChannelIdContextKey, body.ChannelID)
	ctx = context.WithValue(ctx, utils.BotIdContextKey, botId)

	c.SetRequest(c.Request().WithContext(ctx))

	return c
}
