package commands

import (
	"fmt"
	"wakeywakey/database"
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var RegisterWake = discordgo.ApplicationCommand{
	Name:        "register",
	Description: "Registers a PC with a MAC address via alias for Wake-on-LAN.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "alias",
			Description: "The alias of the PC to register.",
			Required: true,
		},
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "mac-address",
			Description: "The MAC address of the PC to register.",
			Required: true,
		},
	},
}

func HandleRegister(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	alias := options[0].StringValue()
	macAddress := options[1].StringValue()

	err := database.AddWakeEntry(i.Member.User.ID, alias, macAddress)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: 1 << 6,
				Embeds: []*discordgo.MessageEmbed{
					utils.EmbedError("Registration Failed", "Failed to register PC: "+err.Error()),
				},
			},
		})
		return
	}

	fmt.Println("Handled registration for user " + i.Member.User.Username + " (" + i.Member.User.ID + "): " + alias + " - " + macAddress)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{
				utils.EmbedSuccess("Registration Successful", "Successfully registered PC."),
			},
		},
	})
}
