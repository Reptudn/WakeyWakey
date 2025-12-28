package commands

import (
	"wakeywakey/utils"

	"github.com/bwmarrin/discordgo"
)

var CommandWake = discordgo.ApplicationCommand{
	Name:        "wake",
	Description: "Sends a Wake-on-LAN packet to wake up your PC.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionString,
			Name: "mac-address",
			Description: "The MAC address of the PC to wake up.",
			Required: true,
		},
	},
}

func CommandWakeHandle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	macAddress := options[0].StringValue()

	err := utils.SendWakeOnLANPacket(macAddress)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to send Wake-on-LAN packet: " + err.Error(),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Wake-on-LAN packet sent successfully!",
		},
	})
}