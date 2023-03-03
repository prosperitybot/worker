package command

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/model"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"go.uber.org/zap"
)

type WhitelabelCommand struct {
	discord.SlashCommand
	db *sqlx.DB
}

func (m WhitelabelCommand) Command() discordgo.ApplicationCommand {
	var defaultPermissions int64 = discordgo.PermissionAdministrator
	return discordgo.ApplicationCommand{
		Name:                     "whitelabel",
		Type:                     discordgo.ChatApplicationCommand,
		Description:              "Manages a whitelabel bot",
		DefaultMemberPermissions: &defaultPermissions,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "setup",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Runs the initial setup for whitelabel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "token",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "Bot Token",
						Required:    true,
					},
					{
						Name:        "public_key",
						Type:        discordgo.ApplicationCommandOptionString,
						Description: "Public Key",
						Required:    true,
					},
				},
			},
			{
				Name:        "actions",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Description: "Manages whitelabel actions",
			},
		},
	}
}

func (m WhitelabelCommand) Execute(c echo.Context, i discordgo.Interaction) {
	subCommand := i.ApplicationCommandData().Options[0]

	switch subCommand.Name {
	case "setup":
		m.subcmd_setup(c, i, subCommand)
	case "actions":
		m.subcmd_actions(c, i, subCommand)
	}
}

func (m WhitelabelCommand) subcmd_setup(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {
	var (
		botToken          = subCommand.Options[0].StringValue()
		publicKey         = subCommand.Options[1].StringValue()
		botAlreadyExists  = false
		userAlreadyHasBot = false
		action            = "start"
		bot               = model.WhitelabelBot{
			Token:     botToken,
			PublicKey: &publicKey,
			Action:    &action,
		}
	)

	if err := bot.FillInfoByToken(); err != nil {
		logger.Error(c.Request().Context(), "Error whilst collecting bot user information", zap.Error(err))
		utils.SendResponse(c, "Invalid bot token", true, true)
		return
	}

	if err := m.db.GetContext(c.Request().Context(), &botAlreadyExists, "SELECT exists (SELECT 1 FROM whitelabel_bots WHERE botId = ?)", bot.Id); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether bot already exists", zap.Error(err))
		utils.SendResponse(c, "Could not activate whitelabel bot", true, true)
		return
	}

	if botAlreadyExists {
		utils.SendResponse(c, "Bot already exists", true, true)
		return
	}

	if err := m.db.GetContext(c.Request().Context(), &userAlreadyHasBot, "SELECT exists (SELECT 1 FROM whitelabel_bots WHERE userId = ?)", i.Member.User.ID); err != nil {
		logger.Error(c.Request().Context(), "Error whilst checking whether user already has a bot", zap.Error(err))
		utils.SendResponse(c, "Could not activate whitelabel bot", true, true)
		return
	}

	if userAlreadyHasBot {
		// Get old bot information
		var oldBot model.WhitelabelBot
		if err := m.db.GetContext(c.Request().Context(), &oldBot, "SELECT * FROM whitelabel_bots WHERE userId = ?", i.Member.User.ID); err != nil {
			logger.Error(c.Request().Context(), "Error whilst getting old bot information", zap.Error(err))
			utils.SendResponse(c, "Could not activate whitelabel bot", true, true)
			return
		}
		bot = oldBot
		bot.OldId = &oldBot.Id
		bot.Token = botToken
		bot.PublicKey = &publicKey
		action := "recreate"
		bot.Action = &action
		if err := bot.FillInfoByToken(); err != nil {
			logger.Error(c.Request().Context(), "Error whilst logging bot user information", zap.Error(err))
			utils.SendResponse(c, "Invalid bot token", true, true)
		}
	}

	// Insert bot into database
	if _, err := m.db.NamedExecContext(c.Request().Context(), "INSERT INTO whitelabel_bots (userId, botId, oldBotId, token, publicKey, action, botName, botDiscrim, botAvatarHash, createdAt, updatedAt) VALUES (:userId, :botId, :oldBotId, :token, :publicKey, :action, :botName, :botDiscrim, :botAvatarHash, :createdAt, :updatedAt)", bot); err != nil {
		logger.Error(c.Request().Context(), "Error whilst inserting bot into database", zap.Error(err))
		utils.SendResponse(c, "Could not activate whitelabel bot", true, true)
		return
	}

	interactionsEndpointUrl := os.Getenv("WORKER_BASE_URL") + "/interactions/" + bot.Id
	developerPage := fmt.Sprintf("https://discord.com/developers/applications/%s/information", bot.Id)

	utils.SendResponse(c, fmt.Sprintf("Whitelabel bot activated\n\nPlease put the following link in `INTERACTIONS ENDPOINT URL` [here](%s): \n`%s`", developerPage, interactionsEndpointUrl), true, false)
}

func (m WhitelabelCommand) subcmd_actions(c echo.Context, i discordgo.Interaction, subCommand *discordgo.ApplicationCommandInteractionDataOption) {

	var (
		embed = utils.CreateEmbed(&discordgo.MessageEmbed{
			Title:       "Whitelabel Bot Actions",
			Description: "Please select an action below",
		}, false)

		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Style:    discordgo.SuccessButton,
						Label:    "Start",
						CustomID: "whitelabel::start",
					},
					discordgo.Button{
						Style:    discordgo.DangerButton,
						Label:    "Stop",
						CustomID: "whitelabel::stop",
					},
					discordgo.Button{
						Style:    discordgo.PrimaryButton,
						Label:    "Restart",
						CustomID: "whitelabel::restart",
					},
					discordgo.Button{
						Style:    discordgo.DangerButton,
						Label:    "Delete",
						CustomID: "whitelabel::delete",
					},
				},
			},
		}
	)

	utils.SendComplexResponse(c, discordgo.InteractionResponseData{
		Flags:      discordgo.MessageFlagsEphemeral,
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
}

func NewWhitelabelCommand(db *sqlx.DB) WhitelabelCommand {
	return WhitelabelCommand{db: db}
}
