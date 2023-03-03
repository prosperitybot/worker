package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/labstack/echo/v4"
)

type SlashCommand interface {
	Command() discordgo.ApplicationCommand
	Execute(c echo.Context, i discordgo.Interaction)
}

type Component interface {
	BaseComponent() discordgo.MessageComponent
	Execute(c echo.Context, i discordgo.Interaction)
}
