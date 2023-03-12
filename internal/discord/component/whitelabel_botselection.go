package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
)

type WhitelabelBotSelectionComponent struct {
	discord.Component
	db *sqlx.DB
}

func (s WhitelabelBotSelectionComponent) BaseComponent() discordgo.MessageComponent {
	return discordgo.SelectMenu{
		CustomID: "whitelabel::botselection",
		MenuType: discordgo.StringSelectMenu,
		Options:  []discordgo.SelectMenuOption{},
	}
}

func (s WhitelabelBotSelectionComponent) Execute(c echo.Context, i discordgo.Interaction) {
	var (
		botId    = i.MessageComponentData().Values[0]
		menuName = fmt.Sprintf("whitelabel::actions_%s", botId)
	)

	utils.SendComplexResponse(c, discordgo.InteractionResponseData{
		Flags: discordgo.MessageFlagsEphemeral,
		Embeds: []*discordgo.MessageEmbed{utils.CreateEmbed(&discordgo.MessageEmbed{
			Description: fmt.Sprintf("Select an action for the bot with id `%s`", botId),
		}, false)},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID: menuName,
						MenuType: discordgo.StringSelectMenu,
						Options: []discordgo.SelectMenuOption{
							{
								Label: "Start",
								Value: "start",
							},
							{
								Label: "Stop",
								Value: "stop",
							},
							{
								Label: "Restart",
								Value: "restart",
							},
							{
								Label: "Delete",
								Value: "delete",
							},
						},
					},
				},
			},
		},
	})
}

func NewWhitelabelBotSelectionComponent(db *sqlx.DB) WhitelabelBotSelectionComponent {
	return WhitelabelBotSelectionComponent{
		db: db,
	}
}
