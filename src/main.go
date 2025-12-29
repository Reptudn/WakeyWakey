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
var SELF_USER_ID string

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

	SELF_USER_ID, exists = os.LookupEnv("SELF_USER_ID")
	if !exists || SELF_USER_ID == "" {
		SELF_USER_ID = ""
	}

	var err error

	dbPath := "wakeywakey.db"
	if os.Getenv("PRODUCTION") == "true" {
		dbPath = "./data/wakeywakey.db"
	}
	err = database.Init(dbPath)
	if err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	BOT, err = discordgo.New("Bot " + BOT_TOKEN)
	if err != nil {
		panic("Failed to create Discord session: " + err.Error())
	}

	BOT.AddHandler(func (session *discordgo.Session, interaction *discordgo.InteractionCreate){

		// this is to restrict the bot to only respond to a specific user if SELF_USER_ID is set
		if SELF_USER_ID != "" &&interaction.Member.User.ID != SELF_USER_ID {
			return
		}

		switch interaction.Type {

			case discordgo.InteractionApplicationCommandAutocomplete: {
				switch interaction.ApplicationCommandData().Name {
					case "wake":
						commands.HandleWakeAutocomplete(session, interaction)
					case "unregister":
						commands.HandleUnregisterAutocomplete(session, interaction)
					case "register":
						commands.HandleRegisterAutocomplete(session, interaction)
				}
			}

			case discordgo.InteractionApplicationCommand: {
				switch interaction.ApplicationCommandData().Name {
					case "wake":
						commands.HandleWake(session, interaction)
					case "register":
						commands.HandleRegister(session, interaction)
					case "unregister":
						commands.HandleUnregister(session, interaction)
				}
			}

			default:
				return
		}

	})

	BOT.AddHandlerOnce(func (session *discordgo.Session, ready *discordgo.Ready){
		appId := session.State.User.ID
		_, err := session.ApplicationCommandCreate(appId, GUILD_ID, &commands.Wake)
		if err != nil {
			panic("Failed to register command: " + err.Error())
		}

		_, err = session.ApplicationCommandCreate(appId, GUILD_ID, &commands.RegisterWake)
		if err != nil {
			panic("Failed to register command: " + err.Error())
		}

		_, err = session.ApplicationCommandCreate(appId, GUILD_ID, &commands.UnregisterWake)
		if err != nil {
			panic("Failed to register command: " + err.Error())
		}
	})

	err = BOT.Open()
	if err != nil {
		panic("Failed to open Discord session: " + err.Error())
	}
	defer BOT.Close()
	
	BOT.UpdateGameStatus(0, "with MAC Adresses")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}