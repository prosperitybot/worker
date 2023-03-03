package middleware

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/utils"
	"go.uber.org/zap"
)

type MiddlewareHandler struct {
	db *sqlx.DB
}

func (h MiddlewareHandler) InteractionAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			botId     = c.Param("bot_id")
			botExists = botId == utils.GetMainBotId()
			isMainBot = botId == utils.GetMainBotId()
			publicKey string
		)

		if !botExists {
			if err := h.db.GetContext(c.Request().Context(), &botExists, "SELECT exists (SELECT 1 FROM whitelabel_bots WHERE botId = ?)", botId); err != nil {
				logger.Error(c.Request().Context(), "Error checking if bot exists", zap.Error(err))
				return c.NoContent(500)
			}
			if err := h.db.GetContext(c.Request().Context(), &publicKey, "SELECT publicKey FROM whitelabel_bots WHERE botId = ?", botId); err != nil {
				logger.Error(c.Request().Context(), "Error getting public key", zap.Error(err))
				return c.NoContent(500)
			}
		}

		if !botExists {
			return c.NoContent(404)
		}

		if isMainBot {
			publicKey = os.Getenv("DISCORD_PUBLIC_KEY")
		}

		pubKeyBytes, err := hex.DecodeString(publicKey)
		if err != nil {
			logger.Error(c.Request().Context(), "Error decoding public key", zap.Error(err))
			c.Error(err)
		}

		pubKey := ed25519.PublicKey(pubKeyBytes)

		if !discordgo.VerifyInteraction(c.Request(), pubKey) {
			c.NoContent(401)
			return nil
		}

		if err := next(c); err != nil {
			c.Error(err)
			return err
		}

		return nil
	}
}

func NewMiddlewareHandler(db *sqlx.DB) MiddlewareHandler {
	return MiddlewareHandler{
		db: db,
	}
}
