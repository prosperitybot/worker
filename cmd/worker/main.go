package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/brpaz/echozap"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"

	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/prosperitybot/common/logger"
	"github.com/prosperitybot/common/model"
	"github.com/prosperitybot/common/utils"
	"github.com/prosperitybot/worker/internal/discord"
	"github.com/prosperitybot/worker/internal/discord/command"
	"github.com/prosperitybot/worker/internal/discord/component"
	"github.com/prosperitybot/worker/internal/http/handler"
	"github.com/prosperitybot/worker/internal/http/middleware"

	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

func main() {
	_ = godotenv.Load()
	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}

	tracer.Start()
	defer tracer.Stop()

	err := profiler.Start(
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer profiler.Stop()

	db := setupDatabase()

	echoInstance := echo.New()

	echoInstance.Use(httptrace.Middleware(httptrace.WithServiceName("worker")))

	// Auth group
	authGroup := echoInstance.Group("")

	middlewareHandler := middleware.NewMiddlewareHandler(db)

	authGroup.Use(middlewareHandler.InteractionAuthMiddleware)
	echoInstance.Use(echozap.ZapLogger(logger.GetLogger()))

	components := map[string]discord.Component{
		"settings::notifications":  component.NewSettingsNotificationComponent(db),
		"whitelabel::botselection": component.NewWhitelabelBotSelectionComponent(db),
		"whitelabel::actions":      component.NewWhitelabelActionsComponent(db),
	}

	commands := map[string]discord.SlashCommand{
		"about":       command.NewAboutCommand(db),
		"ignored":     command.NewIgnoredCommand(db),
		"leaderboard": command.NewLeaderboardCommand(db),
		"level":       command.NewLevelCommand(db),
		"levelroles":  command.NewLevelRolesCommand(db),
		"levels":      command.NewLevelsCommand(db),
		"settings": command.NewSettingsCommand(
			db,
			components["settings::notifications"].(component.SettingsNotificationComponent),
		),
		"whitelabel": command.NewWhitelabelCommand(db),
		"xp":         command.NewXpCommand(db),
	}

	commandList := make([]discordgo.ApplicationCommand, len(commands))
	i := 0
	for j := range commands {
		commandList[i] = commands[j].Command()
		i++
	}

	interactionHandler := handler.InteractionHandler{
		Commands:   commands,
		Components: components,
	}

	healthHandler := handler.HealthHandler{Db: db}

	var (
		whitelabelBots []model.WhitelabelBot
	)

	utils.CreateCommands(commandList, os.Getenv("DISCORD_APPLICATION_ID"), os.Getenv("BOT_TOKEN"), os.Getenv("DEVGUILD_ID"))

	if os.Getenv("ENV") == "prod" {
		if err := db.SelectContext(context.Background(), &whitelabelBots, "SELECT * FROM whitelabel_bots"); err != nil {
			logger.Fatal(context.Background(), "error getting whitelabel bots", zap.Error(err))
		}

		for i := range whitelabelBots {
			utils.CreateCommands(commandList, whitelabelBots[i].Id, whitelabelBots[i].Token, os.Getenv("DEVGUILD_ID"))
		}
	}

	// Routes
	echoInstance.GET("/health", healthHandler.GETHealth)
	authGroup.POST("/interactions/:bot_id", interactionHandler.POSTInteractions)

	data, _ := json.MarshalIndent(echoInstance.Routes(), "", "  ")
	logger.Debug(context.Background(), string(data))

	if err := echoInstance.Start(":3000"); err != nil {
		logger.Fatal(context.Background(), "error starting echo server", zap.Error(err))
	}
}

func setupDatabase() *sqlx.DB {
	sqltrace.Register("mysql", &mysql.MySQLDriver{}, sqltrace.WithServiceName("worker"))
	db, err := sqlxtrace.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
