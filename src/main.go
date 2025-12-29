package main

import (
	"os"
	"os/signal"
	"syscall"
	"wakeywakey/commands"
	"wakeywakey/database"

	"github.com/bwmarrin/discordgo"
)

// TODO: Implement Guild Id for the bot to only operate within a specific server.

var BOT *discordgo.Session
var GUILD_ID string
var BOT_TOKEN string

func main() {
	var exists bool

	BOT_TOKEN, exists = os.LookupEnv("DISCORD_BOT_TOKEN")
	if !exists || BOT_TOKEN == "" {
		panic("DISCORD_BOT_TOKEN environment variable not set")
	}

	GUILD_ID, exists = os.LookupEnv("DISCORD_GUILD_ID")
	if !exists || GUILD_ID == "" {
		panic("DISCORD_GUILD_ID environment variable not set")
	}

	var err error
	
	db, err := database.Init("wakeywakey.db")
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	defer db.Close()
	_ = db

	BOT, err = discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		panic("Failed to create Discord session: " + err.Error())
	}

	BOT.AddHandler(func (session *discordgo.Session, interaction *discordgo.InteractionCreate){

		if interaction.Member.User.ID != os.Getenv("SELF_USER_ID") {
			return
		}

		if interaction.Type != discordgo.InteractionApplicationCommand {
			return
		}

		switch interaction.ApplicationCommandData().Name {
			case "wake":
				commands.CommandWakeHandle(session, interaction)
		}
	
	})

	BOT.AddHandlerOnce(func (session *discordgo.Session, ready *discordgo.Ready){
		appId := session.State.User.ID
		_, err := session.ApplicationCommandCreate(appId, GUILD_ID, &commands.CommandWake)
		if err != nil {
			panic("Failed to register command: " + err.Error())
		}
	})

	err = BOT.Open()
	if err != nil {
		panic("Failed to open Discord session: " + err.Error())
	}
	defer BOT.Close()

	_, err = BOT.ApplicationCommandCreate(BOT.State.User.ID, GUILD_ID, &commands.CommandWake)
	if err != nil {
		panic("Failed to register command: " + err.Error())
	}
	
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}